const createsub_popup = document.getElementById("create_sub_popup");

function hide_popup() {
  createsub_popup.style.display = "none"
}

function show_popup() {
  createsub_popup.style.display = "block"
}

if (document.getElementById("create_subtidder_error").innerHTML != "") {
  show_popup()
} else {
  hide_popup()
}


const createpost_popup = document.getElementById("create_post_popup");

function hide_post_popup() {
  createpost_popup.style.display = "none"
}

function show_post_popup() {
  createpost_popup.style.display = "block"
}

if (document.getElementById("create_post_error").innerHTML != "") {
  show_post_popup()
} else {
  hide_post_popup()
}


const signup_popup = document.getElementById("signup_popup");

function hide_signup_popup() {
  signup_popup.style.display = "none"
}

function show_signup_popup() {
  signup_popup.style.display = "block"
}

if (document.getElementById("signup_error").innerHTML != "") {
  show_signup_popup()
} else {
  hide_signup_popup()
}

const signin_popup = document.getElementById("signin_popup");

function hide_signin_popup() {
  signin_popup.style.display = "none"
}

function show_signin_popup() {
  signin_popup.style.display = "block"
}

if (document.getElementById("signin_error").innerHTML != "") {
  show_signin_popup()
} else {
  hide_signin_popup()
}

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


