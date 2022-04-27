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
    })
})

