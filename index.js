var fs = require('fs');

var express = require('express');
var app = express();

var api = require('./api.js');

app.use('/api', api);

fs.readdir('./backends', function (err, files) {
  if (err != null) {
    console.log('Error loading backends');
    process.exit(1);
  }

  var backends = new Backends();
  files.forEach(function (file) {
    var ext = file.split('.').pop();
    if (ext != 'js') {
      return;
    }

    var modulePath = './backends/' + file;
    try {
      backends.add(file, modulePath);
      console.log(modulePath, 'loaded');
    } catch (e) {
      console.log(modulePath, 'not loaded:', e.message);
    }
  });

  if (backends.length() === 0) {
    console.log('No backends found, exiting');
    process.exit(0);
  }

  app.locals.backends = backends;

  app.listen('3000', function () {
    console.log('mrd listening on port 3000');
  });
});

function Backends () {
  this.backends = {};
}

Backends.prototype.add = function (name, modulePath) {
  this.backends[name] = require(modulePath);
}

Backends.prototype.length = function () {
  return Object.keys(this.backends).length;
}

Backends.prototype.search = function (query) {
  var results = [];
  var that = this;
  Object.keys(this.backends).forEach(function (key) {
    var back = that.backends[key];
    results.push(back.search(query));
  });

  console.log(results);
}
