document.querySelectorAll('.upvote_bt').forEach(item => {
    item.addEventListener('click', event => {
        let element = document.getElementById("pseudonym")
        if (typeof(element) == 'undefined' || element == null) {
            return
        }

        let change = 1
        let src = "../images/global/upvote.png"
        let state = item.getAttribute('state')
        let score = item.nextElementSibling
        let otherBt = item.nextElementSibling.nextElementSibling

        if (state == "active") {
            change = -1
            item.setAttribute('state', 'empty')
            src = "../images/global/empty_upvote.png"
        } else if (otherBt.getAttribute('state') == "active") {
            change = 2
            item.setAttribute('state', 'active')
            otherBt.setAttribute('state', 'empty')
            otherBt.src = "../images/global/empty_downvote.png"
        } else {
            item.setAttribute('state', 'active')
        }

        item.src = src
        score.innerHTML = parseInt(score.innerHTML) + change

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
        let element = document.getElementById("pseudonym")
        if (typeof(element) == 'undefined' || element == null) {
            return
        }
        let change = -1
        let src = "../images/global/downvote.png"
        let state = item.getAttribute('state')
        let score = item.previousElementSibling
        let otherBt = item.previousElementSibling.previousElementSibling

        if (state == "active") {
            change = 1
            item.setAttribute('state', 'empty')
            src = "../images/global/empty_downvote.png"
        } else if (otherBt.getAttribute('state') == "active") {
            change = -2
            item.setAttribute('state', 'active')
            otherBt.setAttribute('state', 'empty')
            otherBt.src = "../images/global/empty_upvote.png"
        } else {
            item.setAttribute('state', 'active')
        }

        item.src = src
        score.innerHTML = parseInt(score.innerHTML) + change

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

