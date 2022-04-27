const posts_tab = document.getElementById("posts_tab")
const infos_tab = document.getElementById("infos_tab")
const admin_tab = document.getElementById("admin_tab")

const posts = document.getElementById("posts");
const infos = document.getElementById("infos");


function show_posts() {
  posts_tab.style.borderBottom = "1vh solid #1A1A1A"
  infos_tab.style.borderBottom = "0"
  admin_tab.style.borderBottom = "0"

  posts.style.display = "block"
  infos.style.display = "none"
}

function show_infos() {
  posts_tab.style.borderBottom = "0"
  infos_tab.style.borderBottom = "1vh solid #1A1A1A"
  admin_tab.style.borderBottom = "0"

  posts.style.display = "none"
  infos.style.display = "block"
}

const getNextSiblings = (e) => {
  let siblings = [];
  while (e = e.nextSibling) {
      siblings.push(e);
  }
  return siblings;
}


document.querySelectorAll('.upvote_bt').forEach(item => {
  item.addEventListener('click', event => {
    if (item.src == "http://localhost/images/global/upvote.png") {  // to change
      item.src = "../images/global/empty_upvote.png"
      item.nextElementSibling.innerHTML = parseInt(item.nextElementSibling.innerHTML) - 1
      return
    }
    item.src = "../images/global/upvote.png"
    item.nextElementSibling.nextElementSibling.src = "../images/global/empty_downvote.png"
    item.nextElementSibling.innerHTML = parseInt(item.nextElementSibling.innerHTML) + 1
  })
})

document.querySelectorAll('.downvote_bt').forEach(item => {
  item.addEventListener('click', event => {
    if (item.src == "http://localhost/images/global/downvote.png") { // to change
      item.src = "../images/global/empty_downvote.png"
      item.previousElementSibling.innerHTML = parseInt(item.previousElementSibling.innerHTML) + 1
      return
    }
    item.src = "../images/global/downvote.png"
    item.previousElementSibling.previousElementSibling.src = "../images/global/empty_upvote.png"
    item.previousElementSibling.innerHTML = parseInt(item.previousElementSibling.innerHTML) - 1
  })
})


function main() {
  show_posts()
}



// SUBSCRIBE TO SUBTIDDER --> to move to subjtidder.js
let bt = document.getElementById("subscribe_bt")
if (bt.innerHTML == "Subscribe") { 
  bt.innerHTML = "Unsubscribe"
  bt.style.backgroundColor = "#666666"
} else {
  bt.innerHTML = "Subscribe"
  bt.style.backgroundColor = "#148AA6"
}


function SubscribeTo(id_account, id_subject) {
  if (bt.innerHTML == "Subscribe") { 
    bt.innerHTML = "Unsubscribe"
    bt.style.backgroundColor = "#666666"
  } else {
    bt.innerHTML = "Subscribe"
    bt.style.backgroundColor = "#148AA6"
  }
  fetch(location.pathname, {
    method: "post",
    headers: {
      'Content-Type': 'application/json'
    },
  
    //make sure to serialize your JSON body
    body: JSON.stringify({
      "id_account_subscribing": id_account,
      "id_subject_to_subscribe": id_subject
    })
  })
}

main()
