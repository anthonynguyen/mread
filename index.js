var fs = require('fs');
var process = require('process');

var express = require('express');
var app = express();

var backends = [];
fs.readdir('./backends', function (err, files) {
  if (err != null) {
    console.log('Error loading backends');
    process.exit(1);
  }

  files.forEach(function (file) {
    var backendPath = './backends/' + file;
    try {
      var backend = require(backendPath);
      backends.push(backend);
      console.log(backendPath, 'loaded');
    } catch (e) {
      console.log(backendPath, 'not loaded');
    }
  });

  if (backends.length === 0) {
    console.log('No backends found, exiting');
    process.exit(0);
  }

  app.listen('3000', function () {
    console.log('mrd listening on port 3000');
  });
});
