const { readFileSync } = require("fs");

const PWD_DIR = "/lps_configuration_data/";
const PWD_FILE = "management_password.txt";

function readPassword() {
    return readFileSync(PWD_DIR+PWD_FILE, "utf-8")
}

module.exports = { readPassword };
