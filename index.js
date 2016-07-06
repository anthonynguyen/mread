var fs = require('fs');

var async = require('async');
var express = require('express');
var leven = require('leven');
var app = express();

var api = require('./api.js');
var log = require('./lib/log.js');

app.set('view engine', 'pug');
app.use('/api', api);

app.get('/', function (req, res) {
	res.render('index');
});

app.get('/search', function (req, res) {
	var query = req.query.q;
	log.warn('query in');
	if (!(query != null)) {
		res.render('search', {success: false, message: 'Please enter a query'});
	} else if (query.length < 4) {
		res.render('search', {success: false, message: 'That query is too short', q: query});
	} else {
		req.app.locals.backends.search(query, function (err, data) {
			if (err != null) {
				return res.sendStatus(500);
			}

			res.render('search', {success: true, data: data, q: query});
			log.warn('served');
		});
	}
});

fs.readdir('./backends', function (err, files) {
	if (err != null) {
		log.error('Could not load backends');
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
			var module = require(modulePath);
			backends.add(module.name, module);
			log.success(module.name, 'backend loaded');
		} catch (e) {
			log.warn(modulePath, 'backend not loaded:', e.message);
		}
	});

	if (backends.length() === 0) {
		log.error('No backends found, exiting');
		process.exit(0);
	}

	app.locals.backends = backends;

	app.listen('5678', function () {
		log.success('mrd listening on port 5678');
	});
});

function Backends () {
	this.backends = {};
}

Backends.prototype.add = function (name, module) {
	this.backends[name] = module;
}

Backends.prototype.length = function () {
	return Object.keys(this.backends).length;
}

Backends.prototype.get = function (name) {
	try {
		return this.backends[name];
	} catch (e) {
		return null;
	}
}

Backends.prototype.search = function (query, callback) {
	var results = {};
	var that = this;
	async.each(Object.keys(this.backends), function(key, cb) {
		var back = that.backends[key];
		back.search(query, function (err, data) {
			if (err != null) {
				return cb();
			}

			data.sort(function (a, b) {
				var la = leven(a.title, query);
				var lb = leven(b.title, query);

				if (la < lb) return -1;
				else if (la > lb) return 1;

				return 0;
			});

			results[key] = data;
			cb();
		});
	}, function (err) {
		if (err != null) {
			if (typeof callback == "function") callback(err, null);
			return;
		}

		if (typeof callback == "function") callback(null, results);
	});
}
