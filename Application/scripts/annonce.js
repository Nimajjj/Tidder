

var slide = new Array("../images/annonce/4.PNG", "../images/annonce/2.PNG", "../images/annonce/3.PNG");

var i = 0;

let chiffre = slide.length;

// fonction qui crée un sommaire en dessous des images

function sommaire(chiffre) {
    for (let i = 0; i < chiffre; i++) {
        let dot = document.createElement("li");
        let newContent = document.createTextNode("");
        dot.classList.add("t"+i);
        dot.dataset.count = i;
        dot.appendChild(newContent);
        document.getElementById("sommaire").appendChild(dot);
    }
}

window.addEventListener("DOMContentLoaded", function () {
    sommaire(chiffre);
    getSlide();
})

function getSlide () {
    document.querySelectorAll("#sommaire li").forEach(function (item) {
        item.addEventListener("click", function () {
                document.getElementById("slide").src = slide[item.dataset.count];
            }
        )
    })
}



// Fonction qui initialise le carrousel d'image

function Carrousel(sens) {
    i = i + sens;
    if (i < 0)
        i = slide.length - 1;
    if (i > slide.length - 1)
        i = 0;
    document.getElementById("slide").src = slide[i];
}

