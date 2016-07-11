var slides = 0;
var current = 0;

var rightToLeft = true;

$(document).ready(function () {
	var rtl = localStorage.getItem('rightToLeft');
	if (rtl != null) {
		if (rtl == 'false') {
			rightToLeft = false;
		} else {
			rightToLeft = true;
		}
	} else {
		localStorage.setItem('rightToLeft', rightToLeft);
	}

	$('#rtl').prop('checked', rightToLeft);
});

$('#rtl').change(function (e) {
	rightToLeft = $(this).prop('checked');
	localStorage.setItem('rightToLeft', rightToLeft);
});

$(document).keydown(function(e) {
	switch(e.which) {
		case 27:
			resetSlides();
			break;

		case 37:
			left();
			break;

		case 39:
			right();
			break;

		default: return;
	}
	e.preventDefault();
});

function resetSlides() {
	$('#gallery').hide();
	slides = [];
	current = 0;
}

function setSlide() {
	if (slides < 1) {
		return;
	}

	if (current < 0 || current >= slides) {
		return;
	}

	var kids = $('#images').children();
	kids.removeClass('current');

	$(kids[current]).addClass('current');

	$('#gallery').css('display', 'flex');
	$('#page-number').text((current + 1) + '/' + slides);
}

function nextSlide() {
	current++;
	if (current >= slides) {
		current = slides - 1;
	}

	setSlide();
}

function previousSlide() {
	current--;
	if (current < 0) {
		current = 0;
	}

	setSlide();
}

function left() {
	if (rightToLeft) {
		nextSlide();
	} else {
		previousSlide();
	}
}

function right() {
	if (rightToLeft) {
		previousSlide();
	} else {
		nextSlide();
	}
}

$('.chapter-link').click(function (e) {
	e.preventDefault();
	$.get($(this).attr('href'), function (data) {
		var images = $('#images');
		images.empty();
		slides = data.length;
		data.forEach(function (i) {
			images.append('<div style="background-image: url(' + i + ')"></div>');
		});
		current = 1;
		previousSlide();
	});
});

$('#gallery').click(function (e) {
	nextSlide();
});

$('#gallery nav').click(function (e) {
	e.stopPropagation();
});
