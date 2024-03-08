const fs = require('fs');
const config = require('./config');

const initFolder = function(dir) {
    if (! fs.existsSync(dir)) {
        console.log(`Creating Box folder`);
        fs.mkdirSync(dir);
    }
}

// Prepare existing files list
const getContent = async function(box) {
    const content = [];
    const files = getAllFiles(box);

    for (i in files) {
        const filepath = files[i];
        const path = filepath.substring(config.BOX_FOLDER_PATH.length+1);
        const stats = fs.statSync(filepath);
        const lastUpdate = Math.floor(stats.mtime.getTime() / 1000);
        content.push({path, lastUpdate});
    }
    return {files: content};
}

function getAllFiles(dir) {
    let results = [];
    const list = fs.readdirSync(dir);
    list.forEach((name) => {
        file = dir + '/' + name;
        const stat = fs.statSync(file);
        if (stat && stat.isDirectory()) { 
            // Get dir files
            results = results.concat(getAllFiles(file));
        } else if (! config.IGNORE_FILES.includes(name)) {
            results.push(file);
        }
    });
    return results;
}

module.exports = { initFolder, getContent }