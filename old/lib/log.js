var chalk = require('chalk');
 
 module.exports = {
	error: function(m) {
		var args = Array.prototype.slice.call(arguments)
		console.log.apply(null, [chalk.red.bold(' ✘ ')].concat(args));
	},
	warn: function(m) {
		var args = Array.prototype.slice.call(arguments)
		console.log.apply(null, [chalk.yellow.bold(' ! ')].concat(args));
	},
	success: function(m) {
		var args = Array.prototype.slice.call(arguments)
		console.log.apply(null, [chalk.green.bold(' ✓ ')].concat(args));
	},
}
