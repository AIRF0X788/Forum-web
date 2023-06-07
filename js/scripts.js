//changer le thème
const themeSwitcher = document.querySelector('#theme-switcher');
const body = document.querySelector('body');
// récupérer le thème stocké dans le cookie (s'il existe)
const savedTheme = getCookie('theme');
// si un thème est stocké dans le cookie, l'appliquer
if (savedTheme) {
  body.classList.add(savedTheme);
}
// ajouter un écouteur d'événement pour changer le thème
themeSwitcher.addEventListener('click', () => {
  body.classList.toggle('dark-theme');
  // stocker le thème dans un cookie
  setCookie('theme', body.classList.contains('dark-theme') ? 'dark-theme' : 'light-theme', 365);
});
// fonction pour récupérer la valeur d'un cookie
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}
// fonction pour stocker une valeur dans un cookie
function setCookie(name, value, days) {
  const date = new Date();
  date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
  const expires = `expires=${date.toUTCString()}`;
  document.cookie = `${name}=${value}; ${expires}; path=/`;
}


//profile
(function() {
    $('.btn-return').click(function() {
      $(this).toggleClass('active');
      return $('.box').toggleClass('open');
    });
  
  }).call(this);

  