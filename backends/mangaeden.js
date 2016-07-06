var fs = require('fs');
var request = require('request');
var util = require('util');

var moment = require('moment');

var fuzzy = require('../lib/fuzzy.js');
var log = require('../lib/log.js');
var reltime = require('../lib/reltime.js');

const DEBUG = process.env.DEBUG == 1;
const MAXIMUM_AGE = 3600000; // milliseconds

const LIST_URL = 'https://www.mangaeden.com/api/list/0/';
const IMAGE_URL = 'https://cdn.mangaeden.com/mangasimg/%s';
const INFO_URL = 'https://www.mangaeden.com/api/manga/%s/';
const CHAPTER_URL = 'https://www.mangaeden.com/api/chapter/%s/';

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
			log.warn('reading from file');
			fs.readFile('./backends/mangaeden.json', function (err, data) {
				handler(err, data);
			});
		} else {
			request.get(LIST_URL, function (err, response, body) {
				handler(err, body);
			});
		}
		latestRetrieval = now;
	} else {
		if (typeof callback == "function") callback(null);
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
					title: series.t,
					id: series.i,
					genres: series.c,
					views: series.h,
				};

				if (series.im != null) {
					r.image = util.format(IMAGE_URL, series.im);
				}

				if (series.s != null) {
					r.status = STATUSES[series.s];
				}

				if (series.ld != null) {
					r.lastChapterDate = reltime(series.ld * 1000);
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

		if (response.statusCode != 200) {
			if (typeof callback == "function") callback(new Error('Manga not found'), 404);
			return;
		}

		var info = JSON.parse(body);
		var r = {
			title: info.title,
			image: util.format(IMAGE_URL, info.image),
			genres: info.categories,
			views: info.hits,
			description: info.description,
			numChapters: info.chapters_len,
			chapters: [],
		}

		if (info.status != null) {
			r.status = STATUSES[info.status];
		}

		if (info.last_chapter_date != null) {
			r.lastChapterDate = reltime(info.last_chapter_date * 1000);
		}

		for (var i = 0; i < info.chapters.length; i++) {
			var chapter = info.chapters[i];
			r.chapters.push({
				number: String(chapter[0]),
				date: moment(chapter[1] * 1000).format('YYYY-MM-DD'),
				title: chapter[2],
				id: chapter[3],
			});
		}

		if (typeof callback == "function") callback(null, r);
	});
}

function chapter (id, callback) {
	request.get(util.format(CHAPTER_URL, id), function (err, response, body) {
		if (err != null) {
			log.error(err);
			if (typeof callback == "function") callback(err, null);
		}

		if (response.statusCode != 200) {
			if (typeof callback == "function") callback(new Error('Chapter not found'), 404);
			return;
		}

		var info = JSON.parse(body);
		var imagesO = {};
		for (var i = 0; i < info.images.length; i++) {
			var im = info.images[i];
			imagesO[im[0]] = util.format(IMAGE_URL, im[1]);
		}

		var images = [];
		Object.keys(imagesO).sort().forEach(function (k) {
			images.push(imagesO[k]);
		});

		if (typeof callback == "function") callback(null, images);
	});
}

module.exports = {
	name: 'Manga_Eden',
	search: search,
	get: get,
	chapter: chapter,
}
