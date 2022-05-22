const posts_tab = document.getElementById("posts_tab")
const infos_tab = document.getElementById("infos_tab")
const admin_tab = document.getElementById("admin_tab")
const tabs = [posts_tab, infos_tab, admin_tab]

const posts = document.getElementById("posts");
const infos = document.getElementById("infos");
const admin = document.getElementById("admin");
const pages = [posts, infos, admin]


function SwitchTab(tab) {
  tabs.forEach(tab => tab.style.borderBottom = "0")
  pages.forEach(page => page.style.display = "none")

  if (tab == "posts") {
    posts_tab.style.borderBottom = "1vh solid #1A1A1A"
    posts.style.display = "block"
  } else if (tab == "infos") {
    infos_tab.style.borderBottom = "1vh solid #1A1A1A"
    infos.style.display = "block"
  } else if (tab == "admin") {
    admin_tab.style.borderBottom = "1vh solid #1A1A1A"
    admin.style.display = "block"
  }
}



// SUBSCRIBE TO SUBTIDDER --> to move to subjtidder.js
let bt = document.getElementById("subscribe_bt")
if (typeof(bt) != 'undefined' && bt != null) {
  if (bt.innerHTML == "Subscribe") { 
    bt.innerHTML = "Unsubscribe"
    bt.style.backgroundColor = "#666666"
  } else {
    bt.innerHTML = "Subscribe"
    bt.style.backgroundColor = "#148AA6"
  }
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

SwitchTab("posts")

let info = document.querySelector("#info")

let edit = document.querySelector("#edit_info")
let textarea = document.querySelector("#textarea_info")
let submit = document.querySelector("#submit_info")
let cancel = document.querySelector("#cancel_info")

if (edit != null) {
  edit.addEventListener('click', event => {
    textarea.style.display = "block"
    submit.style.display = "block"
    cancel.style.display = "block"
    edit.style.display = "none"
    info.style.display = "none"
  
    textarea.value = info.innerHTML
  })
  
}

if (submit != null) {
  submit.addEventListener('click', event => {
    info.innerHTML = textarea.value

    textarea.style.display = "none"
    submit.style.display = "none"
    cancel.style.display = "none"
    edit.style.display = "block"
    info.style.display = "block"

    fetch(location.pathname, {
      method: "post",
      headers: {
        'Content-Type': 'application/json'
      },
    
      //make sure to serialize your JSON body
      body: JSON.stringify({
        "info": info.innerHTML,
      })
    })
  })
}

if (cancel != null) {
  cancel.addEventListener('click', event => {
    textarea.style.display = "none"
    submit.style.display = "none"
    cancel.style.display = "none"
    edit.style.display = "block"
    info.style.display = "block"
  })
}

function readURL(input, what) {
  if (input.files && input.files[0]) {
    var reader = new FileReader();
    reader.onload = function (e) {
      $('#' + what)
        .attr('src', e.target.result)
        .attr('max-height', 660)
    };
    reader.readAsDataURL(input.files[0]);
  }
}




const admin_gen_tab = document.getElementById("admin_tab_gen")
const admin_rol_tab = document.getElementById("admin_tab_rol")
const admin_use_tab = document.getElementById("admin_tab_users")
const admin_tabs = [admin_gen_tab, admin_rol_tab, admin_use_tab]

const gen = document.getElementById("admin_gen");
const rol = document.getElementById("admin_roles");
const use = document.getElementById("admin_users");
const admin_pages = [gen, rol, use]


function SwitchAdminTab(tab) {
  admin_tabs.forEach(tab => tab.classList.remove("admin_tab_selected"))
  admin_pages.forEach(page => page.style.display = "none")

  if (tab == "general") {
    admin_gen_tab.classList.add("admin_tab_selected")
    gen.style.display = "block"
  } else if (tab == "roles") {
    admin_rol_tab.classList.add("admin_tab_selected")
    rol.style.display = "block"
  } else if (tab == "users") {
    admin_use_tab.classList.add("admin_tab_selected")
    use.style.display = "block"
  }
}


function UpdateBannedUser() {
  let changes = []
  document.querySelectorAll(".ban_checkbox").forEach((item) => {
    oldState = item.getAttribute("ban")
    currentState = item.checked
    if ((oldState === "true") != currentState) {
      changes.push(item.getAttribute("name"))
    }
  })

  let res = ""

  changes.forEach((e) => {
    res += e + ";"
  })

  if (res != "") {
    return res;
  }
} 


