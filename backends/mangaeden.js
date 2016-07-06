var fs = require('fs');
var moment = require('moment');
var request = require('request');

var fuzzy = require('../lib/fuzzy.js');

const DEBUG = process.env.DEBUG == 1;
const LIST_URL = 'https://www.mangaeden.com/api/list/0/';
const MAXIMUM_AGE = 3600000; // milliseconds

var latestRetrieval = moment().subtract(10, 'hours');

var list = null;

function getList () {
	// If debug mode is on use a local copy of the list so that we don't download a new one each time
	if (DEBUG) {
		var data = fs.readFileSync('./backends/mangaeden.json');
		list = JSON.parse(data);
	} else {
		request.get(LIST_URL, function (err, data) {
			list = JSON.parse(data);
		});
	}
}

function refreshList () {
	if (moment().diff(latestRetrieval) > MAXIMUM_AGE) {
		getList();
		latestRetrieval = moment();
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
function search (query) {
	refreshList();

	var results = [];
	list.forEach(function (series) {
		if (fuzzy(query, series.t)) {
			results.push(series);
		}
	});
	return results;
}

module.exports = {
	search: search,
}
