module.exports = function (time) {
	var diff = Date.now() - time;
	var absolute = new Date(time);
	var str = '(' + absolute.getUTCDate() + '/' + (absolute.getUTCMonth() + 1) + '/' + absolute.getUTCFullYear() + ')';
	diff /= 1000;

	var days = diff / 60 / 60 / 24;

	if (days < 1) {
		return 'Today ' + str;
	} else if (days < 2) {
		return 'Yesterday ' + str;
	} else {
		return parseInt(days) + ' days ago ' + str;
	}
}
