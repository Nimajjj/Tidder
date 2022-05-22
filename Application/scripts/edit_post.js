function EditPost() {
    let textContent = document.querySelector("#post_text_content");

    if (typeof(textContent) != "undefined" && textContent != null) {
        let textarea = document.createElement("textarea");
        textarea.id = "post_text_content_edit";
        textarea.value = textContent.innerHTML;
        let submitBt = document.createElement("div");
        submitBt.id = "submit_edit_post";
        submitBt.innerHTML = "Submit";

        submitBt.addEventListener("click", function() {
            let text = textarea.value;
            let url = window.location.href

            let count = 0
            for (let i = 0; i < url.length; i++) {
                if (count != 2) {
                    if (url[i] == "/") {
                        count += 1
                    }
                    url = url.substring(i + 1);
                }
            }
            url = url.substring(1);
            
            fetch(location.pathname, {
                method: "post",
                headers: {
                  'Content-Type': 'application/json'
                },
              
                //make sure to serialize your JSON body
                body: JSON.stringify({
                  "post_id": url,
                  "new_post_text_content": text,
                })
              }).then(() => {
                window.location.reload();
              })

        })

        textContent.parentNode.insertBefore(textarea, textContent);
        textContent.parentNode.insertBefore(submitBt, textContent);

        textContent.remove();
    }
}


function DeletePost() {
    var r = confirm("Are you sure you want to delete this post?");
    if (r == true) {
        let url = window.location.href

        let count = 0
        for (let i = 0; i < url.length; i++) {
            if (count != 2) {
                if (url[i] == "/") {
                    count += 1
                }
                url = url.substring(i + 1);
            }
        }
        url = url.substring(1);

        fetch(location.pathname, {
            method: "post",
            headers: {
              'Content-Type': 'application/json'
            },
          
            //make sure to serialize your JSON body
            body: JSON.stringify({
              "delete_post": url,
            })
          }).then(() => {
            window.location.href = "/";
          })
    }
}