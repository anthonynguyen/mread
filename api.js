var express = require('express');
var router = express.Router();

router.all('/', function (req, res) {
	res.send('api');
});

/*
Superset of format of results:
not all results will have each of these items
[
	{
		id: unique id,
		title: Manga Title,
		image: http://example.com/url/to/image,
		status: completed/in progress,
		genres: [action, adventure],
		lastChapterDate: DATE,
		views: 10000,
	}
	...
	...
]
*/
router.get('/search/:query', function (req, res) {
	var query = req.params.query;
	if (query.length < 5) {
		res.status(400).send('Search query is too short');
	} else {
		req.app.locals.backends.search(query, function (err, data) {
			if (err != null) {
				return res.sendStatus(500);
			}

			res.send(data);
		});
	}
});

/*
Superset of format of result:
not all results will have each of these items
{
	title: Manga Title,
	image: http://example.com/url/to/image,
	status: completed/in progress,
	genres: [action, adventure],
	lastChapterDate: DATE,
	views: 10000,
	description: blah blah,
	numChapters: 10,
	released: 2016,
	chapters: [
		{
			number: "123",
			date: DATE,
			title: "Chapter title",
			id: asdjalsjl3983a
		}
	]
}
*/
router.get('/manga/:backend/:id', function (req, res) {
	var requestedBackend = req.params.backend;
	var id = req.params.id;

	var backend = req.app.locals.backends.get(requestedBackend);
	if (backend != null) {
		backend.get(id, function (err, data) {
			if (err != null) {
				var code = Number(data);
				return res.sendStatus(isNaN(code) ? 500 : code);
			}

			res.send(data);
		});
	} else {
		res.sendStatus(404);
	}
});

/*
Format of result:
[
	image url,
	image url,
	...
	...
]
*/
router.get('/chapter/:backend/:id', function (req, res) {
	var requestedBackend = req.params.backend;
	var id = req.params.id;

	var backend = req.app.locals.backends.get(requestedBackend);
	if (backend != null) {
		backend.chapter(id, function (err, data) {
			if (err != null) {
				var code = Number(data);
				return res.sendStatus(isNaN(code) ? 500 : code);
			}

			res.send(data);
		});
	} else {
		res.sendStatus(404);
	}
});

module.exports = router;
