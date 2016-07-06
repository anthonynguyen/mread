var express = require('express');
var router = express.Router();

router.all('/', function (req, res) {
	res.send('api');
});

router.get('/:backend/get/:id', function (req, res) {
	var requestedBackend = req.params.backend;
	var id = req.params.id;

	var backend = req.app.locals.backends.get(requestedBackend);
	if (backend != null) {
		backend.get(id, function (err, data) {
			if (err != null) {
				return res.sendStatus(500);
			}

			res.send(data);
		});
	} else {
		res.sendStatus(404);
	}
});

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

module.exports = router;
