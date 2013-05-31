function submitLinks() {
	if (document.getElementById("spotify_links").value == '') {
		alert('There is nothing to submit! Paste some urls!');
	} else {
		$.post('/submit', $('#songs_input').serialize(), function (data) {
			$('#content').html(data);
		});
	}
}