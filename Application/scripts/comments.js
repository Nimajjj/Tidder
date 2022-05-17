function AnswerToComment(id) {
    let forms = document.querySelectorAll(".answer_form")
    forms.forEach(f => {
        f.remove()
    })


    let parent = document.querySelector("#comment_" + id)

    let form = document.createElement("form")
    form.setAttribute("method", "post")
    form.setAttribute("name", "answer")
    form.setAttribute("id", "answer")
    form.classList.add("answer_form")

    let textarea = document.createElement("textarea")
    textarea.setAttribute("name", "answer_content")
    textarea.rows = 5
    textarea.style.width = "100%"

    let div = document.createElement("div")
    div.classList.add("button")
    div.onclick = function(){answer(id)}


    let p = document.createElement("p")
    p.innerHTML = "Submit"

    let idHolder = document.createElement("input")
    idHolder.setAttribute("type", "hidden")
    idHolder.setAttribute("name", "comment_id")
    idHolder.setAttribute("value", id)

    form.appendChild(textarea)
    div.appendChild(p)
    form.appendChild(div)
    form.appendChild(idHolder)

    parent.appendChild(form)
}


function answer(id) {
    console.log('you actually try something')

    let parent = document.querySelector("#comment_" + id)
    let content

    let done = false
    parent.childNodes.forEach(item => {
        if (item.tagName == "FORM") {
            item.childNodes.forEach(iitem => {
                if (iitem.tagName == "TEXTAREA") {
                    content = iitem.value
                }
            })
        }
    })


    fetch(location.pathname, { 
        method: "post",
        headers: {
          'Content-Type': 'application/json'
        },
      
        body: JSON.stringify({
          "id_response_to": parseInt(id),
          "content": content
        })
      }).then(() => {
        window.location.reload();
    })
}



document.querySelectorAll('.comments_upvote_bt').forEach(item => {
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
              "id_comment": parseInt(id),
              "score_comment": change
            })
          })
    })
})

document.querySelectorAll('.comments_downvote_bt').forEach(item => {
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
              "id_comment": parseInt(id),
              "score_comment": change
            })
          })
    })
})

