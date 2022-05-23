function readURL(input) {
    if (input.files && input.files[0]) {
        var reader = new FileReader();
        reader.onload = function (e) {
            $('#default_png').attr('src', e.target.result).attr('max-height', 660)
        };
        reader.readAsDataURL(input.files[0]);

        if (document.querySelectorAll('.sumbit_bt').length > 0) { return; }

        let submitBt = document.createElement('div');
        submitBt.className = 'sumbit_bt';
        submitBt.innerHTML = 'Submit';

        submitBt.addEventListener('click', function () {
            document.forms['profile_picture'].submit().then(() => {
                window.location.reload();
            })
        })
        

        input.parentNode.insertBefore(submitBt, input.nextSibling);
    }
  }