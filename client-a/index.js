const app = require('./src/app.js')

function checkHelpArg() {
    if (process.argv.includes('-help')) {
        console.log('\nClient Demo JavaScript\n');
        console.log('Usage: node index.js [options]');
        console.log('');
        console.log('Options:');
        console.log('  -host      Server Host | Default to 127.0.0.1');
        console.log('  -port      Server Port | Default to 8080');
        console.log('  -api-key   Server API Key | Default to 12345');
        console.log('');

        process.exit(0); // Exit the script
    }
}

checkHelpArg();
app.run();