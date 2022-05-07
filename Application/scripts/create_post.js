function readURL(input) {
    if (input.files && input.files[0]) {
      var reader = new FileReader();
      reader.onload = function (e) {
        $('#media_preview')
          .attr('src', e.target.result)
          .attr('max-height', 660)
      };
      reader.readAsDataURL(input.files[0]);
    }
  }