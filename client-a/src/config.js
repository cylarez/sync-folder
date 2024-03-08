const path = require('path');

const IGNORE_FILES = ".DS_Store"
const BOX_FOLDER_NAME = "box";
const BOX_FOLDER_PATH = path.join(path.dirname(require.main.filename), BOX_FOLDER_NAME);

module.exports = {IGNORE_FILES, BOX_FOLDER_PATH}