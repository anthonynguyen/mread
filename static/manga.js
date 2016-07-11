var slides = 0;
var current = 0;
var rightToLeft = true;

$(document).ready(function () {
	var rtl = Lockr.get('rtl', true);

	if (rtl) rightToLeft = true;
	else rightToLeft = false;

	$('#rtl').prop('checked', rightToLeft);

	// Add unread badges to everything
	$('.chapter-link').each(function () {
		var backendName = $(this).attr('data-backend');
		var chapterID = $(this).attr('data-chapter-id');
		console.log(backendName + ', ' + chapterID);
		if (!Lockr.sismember(backendName + 'ReadChapters', chapterID)) {
			$(this).find('.is-unread').addClass('new');
		} else {
			$(this).find('.is-unread').removeClass('new');
		}
	});
});

$('#rtl').change(function (e) {
	rightToLeft = $(this).prop('checked');
	Lockr.set('rtl', rightToLeft);
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
	slides = 0;
	current = 0;
}

function setSlide() {
	if (slides < 1) return;

	if (current < 0 || current >= slides) return;

	var kids = $('#images').children();
	kids.removeClass('current');

	$(kids[current]).addClass('current');

	$('#gallery').css('display', 'flex');
	$('#page-number').text((current + 1) + '/' + slides);
}

function nextSlide() {
	current++;
	if (current >= slides) current = slides - 1;
	setSlide();
}

function previousSlide() {
	current--;
	if (current < 0) current = 0;
	setSlide();
}

function left() {
	if (rightToLeft) nextSlide();
	else previousSlide();
}

function right() {
	if (rightToLeft) previousSlide();
	else nextSlide();
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

	Lockr.sadd($(this).attr('data-backend') + 'ReadChapters', $(this).attr('data-chapter-id'));
	$(this).find('.is-unread').removeClass('new');
});

$('#gallery').click(function (e) {
	nextSlide();
});

$('#gallery nav').click(function (e) {
	e.stopPropagation();
});
