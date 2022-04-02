const popup = document.getElementById("create_sub_popup");

function hide_popup() {
  popup.style.display = "none"
}

function show_popup() {
  popup.style.display = "block"
}


if (document.getElementById("create_subtidder_error").innerHTML != "") {
  show_popup()
} else {
  hide_popup()
}


function SubscribeTo(id_account, id_subject) {
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


function SubscribeTo(id_account, id_subject) {
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