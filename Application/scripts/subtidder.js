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

function main() {
  show_posts()
}

main()
