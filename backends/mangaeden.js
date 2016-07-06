var fs = require('fs');
var request = require('request');
var util = require('util');

var fuzzy = require('../lib/fuzzy.js');
var log = require('../lib/log.js');

const DEBUG = process.env.DEBUG == 1;
const MAXIMUM_AGE = 3600000; // milliseconds

const LIST_URL = 'https://www.mangaeden.com/api/list/0/';
const IMAGE_URL = 'https://cdn.mangaeden.com/mangasimg/%s';
const INFO_URL = 'https://www.mangaeden.com/api/manga/%s/';

const STATUSES = ['Suspended', 'Ongoing', 'Completed'];

var latestRetrieval = Date.now() - MAXIMUM_AGE;

var list = null;

function refreshList (callback) {
	var now = Date.now();
	var handler = function (err, data) {
		if (err != null) {
			if (typeof callback == "function") callback(err);
		}
		list = JSON.parse(data);
		if (typeof callback == "function") callback(null);
	}

	if (now - latestRetrieval > MAXIMUM_AGE) {
		log.warn('Mangaeden list too old, getting another');
		if (DEBUG) {
			fs.readFile('./backends/mangaeden.json', function (err, data) {
				handler(err, data);
			});
		} else {
			request.get(LIST_URL, function (err, response, body) {
				handler(err, body);
			});
		}
		latestRetrieval = now;
	}
}

/*
	{
		"a": "one-piece",
		"h": 199232178,
		"ld": 1467287417.0,
		"s": 1,
		"t": "One Piece",
		"i": "4e70ea10c092255ef7004aa2",
		"im": "b0/b0ac7f12d2cb0fc07b9418d5544a3f97cbbc30e967396ae70f98d101.png",
		"c": [
			"Action",
			"Adventure",
			"Comedy",
			"Drama",
			"Fantasy",
			"Sci-fi",
			"Shounen"
		]
	},
*/
function search (query, callback) {
	refreshList(function (err) {
		if (err != null) {
			log.error('Error refreshing list:', err.message);
		}

		if (list === null) {
			if (typeof callback == "function") callback(new Error('No list to search'), null);
			return;
		}

		var results = [];
		list.forEach(function (series) {
			if (fuzzy(query, series.t)) {
				var r = {
					image: util.format(IMAGE_URL, series.im),
					title: series.t,
					id: series.i,
					genres: series.c,
					views: series.h,
				};

				if (series.s != null) {
					r.status = STATUSES[series.s];
				}

				if (series.ld != null) {
					r.lastChapterDate = new Date(series.ld * 1000).toISOString();
				}

				results.push(r);
			}
		});
		if (typeof callback == "function") callback(null, results);
	});
}

function get (id, callback) {
	request.get(util.format(INFO_URL, id), function (err, response, body) {
		if (err != null) {
			log.error(err);
			if (typeof callback == "function") callback(err, null);
		}
		var info = JSON.parse(body);
		if (typeof callback == "function") callback(null, info);
	});
}

module.exports = {
	name: 'Manga Eden',
	search: search,
	get: get,
}
