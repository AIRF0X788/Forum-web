function simulateButtonClick() {
    var currentButton = localStorage.getItem('currentButton');
    
    if (currentButton === null || currentButton === "myButton2") {
      setTimeout(function() {
        document.getElementById("myButton").click();
        localStorage.setItem('currentButton', 'myButton');
      }, 1000);
    } else {
      setTimeout(function() {
        document.getElementById("myButton2").click();
        localStorage.setItem('currentButton', 'myButton2');
      }, 1000);
    }
  }

function getCookie(name) {
    var cookieArr = document.cookie.split("; ");
    for (var i = 0; i < cookieArr.length; i++) {
        var cookiePair = cookieArr[i].split("=");
        if (name === cookiePair[0]) {
            return decodeURIComponent(cookiePair[1]);
        }
    }
    return null;
}
function setCookie(name, value, days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + encodeURIComponent(value) + expires + "; path=/";
}

function toggleTheme() {
    var button = document.getElementById("myButton");
    var currentTheme = button.getAttribute("data-theme");
    var newTheme = currentTheme === "dark" ? "light" : "dark";

    button.classList.add("changing");

    setTimeout(function () {
        button.setAttribute("data-theme", newTheme);
        button.classList.toggle("dark", newTheme === "dark");

        button.classList.remove("changing");

        setCookie("buttonTheme", newTheme, 365);
    }, 900);
}

document.addEventListener("DOMContentLoaded", function () {
    var button = document.getElementById("myButton");
    var storedTheme = getCookie("buttonTheme");

    if (storedTheme) {
        button.setAttribute("data-theme", storedTheme);
        button.classList.toggle("dark", storedTheme === "dark");
    }
});