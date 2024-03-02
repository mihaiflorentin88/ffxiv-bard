window.onload = function() {
    const deleteButton = document.getElementById('deleteButton');

    deleteButton.addEventListener('click', function() {
        const songId = this.getAttribute('data-song-id');
        const userConfirmed = confirm('Are you sure you want to delete this song?');
        if (userConfirmed) {
            fetch(`/song/${songId}`, {
                method: 'DELETE',
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.text();
                })
                .then(() => {
                    window.location.href = '/song/list';
                })
                .catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                });
        }
    });
};
