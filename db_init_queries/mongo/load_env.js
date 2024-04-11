const fs = require('fs');

const envFile = fs.readFileSync('.env', 'utf8');
const envVariables = envFile.split('\n');

const env = {};
envVariables.forEach(variable => {
    const [key, value] = variable.split('=');
    env[key] = value;
});

module.exports = env;
