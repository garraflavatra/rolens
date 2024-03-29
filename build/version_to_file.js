#!/usr/bin/node

// This script extracts the version number from wails.json in the project root
// and writes it to version.txt

const fs = require('fs');

fs.writeFileSync(
  __dirname + '/version.txt',
  JSON.parse(
    fs.readFileSync(__dirname + '/../wails.json')
  ).info.productVersion
);
