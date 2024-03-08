const http = require('node:http');
const fs = require('fs');
const path = require('path');
const querystring = require("querystring");
const fileBox = require('./file');
const config = require('./config');

const SYNC_DOWNLOAD = 1, SYNC_UPLOAD = 2;

// Check arguments to override default config
const { host, port, apiKey } = parseArgs();
const SERVER_HOST = host || "127.0.0.1", 
SERVER_PORT = port || "8080",
SERVER_API_KEY = apiKey || "12345";

const agent = new http.Agent({
    keepAlive: true
});
const serverOptions = {
    hostname: SERVER_HOST,
    agent,
    method: 'POST',
    protocol: 'http:',
    port: SERVER_PORT,
    headers: {
        'Authorization': SERVER_API_KEY,
        'Content-Type': 'application/json',
        'Connection': 'Keep-Alive'
    }
}

function reconnect() {
    console.log('[Sync]:error reconnecting in 5sec...\n')
    setTimeout(run, 5000)
}

function sync(content) {
    console.info(`[Sync]:connecting to:${SERVER_HOST}:${SERVER_PORT}`);
    const options = {
        ...serverOptions,
        path: '/sync'
    };
    let stream = '';
    const req = http.request(options, (res) => {
        res.setEncoding('utf8');
        res.on('error', (err) => {
            console.log(`[Sync]:error code: ${err.code}`)
        });
        res.on('data', (data) => {
            if (res.statusCode !== 200) {
                console.error(`[Sync]:error code:${res.statusCode} ${data}`)
                return;
            }
            // Receive data from Stream
            stream += data;
            if (data.slice(-1) === '\n') {
                // Current stream is completed
                let file;
                try {
                    file = JSON.parse(stream);
                } catch (err) {
                    console.error(err);
                    return;
                }
                switch (file.Action) {
                    case SYNC_DOWNLOAD:
                        download(file);
                        break;
                    case SYNC_UPLOAD:
                        upload(file);
                        break;
                    default:
                        console.error(`[Sync]:error unsupported action for file:`, file);
                        break;
                }
                stream = '';
            }
        })
        
    }).on('close', () => {
        console.log('[Sync]:disconnection...');
        reconnect();
    })
    .on('response', () => {
        console.info(`[Sync]:sent ${content.files.length} files`);
        console.info(`[Sync]:waiting for updates...`);
    })
    .on('error', (err) => {
        console.log(`[Sync]:error code: ${err.code}`)
        // retry if needed
        if (req.reusedSocket && err.code === 'ECONNRESET') {
            reconnect();
        }
    });
    req.write(JSON.stringify(content));
    req.end();
}

function upload(file) {
    console.info(`[Upload]:sending file: ${file.Path}`);
    // Load file content
    const filePath = path.join(config.BOX_FOLDER_PATH, file.Path)
    file.Data = fs.readFileSync(filePath);
    const options = {
        ...serverOptions,
        path: '/upload',
        headers: {...serverOptions.headers, ...{
            'Content-Type': 'application/octet-stream',
            'Content-Length': file.Data.length,
            'File-Path': file.Path
        }},
    };

    const req = http
        .request(options)
        .on('error', function(e) {
            console.error('[Upload]:error: ' + e.message);
        })
        .on('close', () => {
            console.error(`[Upload]:completed file:${file.Path}`);
        })
    req.write(file.Data);
    req.end();
}

function download(f) {
    console.info(`[[Download]]:starting file:${f.Path}`);
    const filePath = path.join(config.BOX_FOLDER_PATH, f.Path);
    const options = {
        ...serverOptions,
        path: `/download/${querystring.escape(f.Path)}`,
        method: 'GET',
    };
    http.get(options, (res) => {
        if (res.statusCode !== 200) {
            console.error(`[Download]:error with file:${f.Path} code:${res.statusCode}`);
            res.resume();
            return;
        }
        // Create folder if needed
        const dir = path.dirname(filePath);
        if (!fs.existsSync(dir)) {
            console.log(`[Download]:creating sub-directory: ${dir}`);
            fs.mkdirSync(dir, { recursive: true });
        }
        // Allows progressive download and let node manage allocated memory
        const file = fs.createWriteStream(filePath);
        res.pipe(file);
        file.on('finish', () =>{
            file.close(() => {
                console.log(`[Download]:completed file:${f.Path}`)
            });
        });
    })  
}

function parseArgs() {
    let host, port, apiKey;
    process.argv.forEach((arg) => {
        if (arg.startsWith("-host")) {
            host = arg.substring("-host".length+1)
        } else if (arg.startsWith("-port")) {
            port = arg.substring("-port".length+1)
        } if (arg.startsWith("-api-key")) {
            apiKey = arg.substring("-api-key".length+1)
        }
    })
    return {host, port, apiKey};
}

const run = () => {
    fileBox.initFolder(config.BOX_FOLDER_PATH);
    fileBox.getContent(config.BOX_FOLDER_PATH).then(content => {
        sync(content)
    })
}

module.exports = {run}