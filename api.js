var express = require('express');
var router = express.Router();

router.all('/', function (req, res) {
  res.send('api');
});

router.get('/search/:query', function (req, res) {
  var query = req.params.query;
  if (query.length < 5) {
    res.send('Search query is too short');
  } else {
    res.send(query);
  }
});

module.exports = router;
