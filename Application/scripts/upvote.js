document.querySelectorAll('.upvote_bt').forEach(item => {
    item.addEventListener('click', event => {
        let change = 1
        let src = "../images/global/upvote.png"
        let score = item.nextElementSibling
        let otherBt = item.nextElementSibling.nextElementSibling

        if (item.src == "http://localhost/images/global/upvote.png") {
            change = -1
            src = "../images/global/empty_upvote.png"
        } else if (otherBt.src == "http://localhost/images/global/downvote.png") {
            change = 2
            src = "../images/global/upvote.png"
        }
        item.src = src
        score.innerHTML = parseInt(score.innerHTML) + change
        otherBt.src = "../images/global/empty_downvote.png"


        let id = ""
        for (let i = 0; i < item.alt.length; i++) {
            if (item.alt[i] == " ") {
                break
            }
            id += item.alt[i]
        }

        fetch(location.pathname, {  // vote query
            method: "post",
            headers: {
              'Content-Type': 'application/json'
            },
          
            body: JSON.stringify({
              "id_post": parseInt(id),
              "score": change
            })
          })
    })
})

document.querySelectorAll('.downvote_bt').forEach(item => {
    item.addEventListener('click', event => {
        let change = -1
        let src = "../images/global/downvote.png"
        let score = item.previousElementSibling
        let otherBt = item.previousElementSibling.previousElementSibling

        if (item.src == "http://localhost/images/global/downvote.png") {
            change = 1
            src = "../images/global/empty_downvote.png"
        } else if (otherBt.src == "http://localhost/images/global/upvote.png") {
            change = -2
            src = "../images/global/downvote.png"
        }
        item.src = src
        score.innerHTML = parseInt(score.innerHTML) + change
        otherBt.src = "../images/global/empty_upvote.png"

        let id = ""
        for (let i = 0; i < item.alt.length; i++) {
            if (item.alt[i] == " ") {
                break
            }
            id += item.alt[i]
        }

        fetch(location.pathname, {  // vote query
            method: "post",
            headers: {
              'Content-Type': 'application/json'
            },
          
            body: JSON.stringify({
              "id_post": parseInt(id),
              "score": change
            })
          })
    })
})

