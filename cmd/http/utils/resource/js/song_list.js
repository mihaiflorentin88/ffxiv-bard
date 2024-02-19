document.addEventListener('DOMContentLoaded', function () {
    function setQueryParamValue(elementId, paramName) {
        var element = document.getElementById(elementId);
        if (element) {
            const params = new URLSearchParams(window.location.search);
            const value = params.get(paramName);
            if (value !== null) {
                element.value = value;
            }
        }
    }

    setQueryParamValue('#pageParam', 'page');
});


document.addEventListener('DOMContentLoaded', function () {
    function updatePlaceholderStyling(selectElement) {
        if (selectElement.value === "") {
            selectElement.classList.add('placeholder');
        } else {
            selectElement.classList.remove('placeholder');
        }
    }

    var selects = document.querySelectorAll('.form-control');
    selects.forEach(function (select) {
        updatePlaceholderStyling(select);

        select.addEventListener('change', function () {
            updatePlaceholderStyling(select);
        });

        select.addEventListener('focus', function () {
            updatePlaceholderStyling(select);
        });
        select.addEventListener('blur', function () {
            updatePlaceholderStyling(select);
        });
    });
});

document.addEventListener('DOMContentLoaded', function () {
    function getQueryParam(param) {
        var urlParams = new URLSearchParams(window.location.search);
        return urlParams.get(param);
    }

    var ensembleSize = getQueryParam('ensembleSize');
    var genreId = getQueryParam('genre');

    if (ensembleSize) {
        var ensembleSizeElement = document.getElementById('ensembleSize');
        ensembleSizeElement.value = ensembleSize;
    }

    if (genreId) {
        var genreElement = document.getElementById('genre');
        genreElement.value = genreId;
    }
});
