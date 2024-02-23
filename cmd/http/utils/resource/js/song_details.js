document.addEventListener('DOMContentLoaded', function() {
    function editComment(commentId, content, updateUrl) {
        var commentDiv = document.querySelector('.comment-content[data-comment-id="' + commentId + '"]');
        var textarea = document.createElement('textarea');
        textarea.setAttribute('maxlength', '500');
        textarea.value = content;
        textarea.rows = 5;
        textarea.classList.add('form-control');
        var saveButton = document.createElement('button');
        saveButton.textContent = 'Save';
        saveButton.classList.add('btn', 'btn-primary', 'save-comment');
        saveButton.setAttribute('type', 'button');
        saveButton.onclick = function() {
            submitUpdatedComment(commentId, textarea.value, updateUrl);
        };

        commentDiv.innerHTML = '';
        commentDiv.appendChild(textarea);
        commentDiv.appendChild(saveButton);
    }

    function submitUpdatedComment(commentId, updatedContent, updateUrl) {
        fetch(updateUrl, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json', // Set the Content-Type header
            },
            body: JSON.stringify({ commentId: parseInt(commentId), content: updatedContent }),
            redirect: 'error' // Prevent following redirects
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Could not save document');
                } else {
                    var commentDiv = document.querySelector('.comment-content[data-comment-id="' + commentId + '"]');
                    commentDiv.innerHTML = `<p>${updatedContent}</p>`;
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Error saving comment.');
            });
    }

    var editIcons = document.querySelectorAll('.edit-comment');
    editIcons.forEach(function(icon) {
        icon.addEventListener('click', function(event) {
            event.preventDefault();
            var commentId = this.getAttribute('data-comment-id');
            var content = this.parentNode.nextElementSibling.querySelector('p').textContent.trim();
            var updateUrl = this.getAttribute('data-url');
            editComment(commentId, content, updateUrl);
        });
    });
});
