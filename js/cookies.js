//Bouton qui fait disparaitre la barre cookie
document.getElementById('cookie-accept').addEventListener('click', function() {
    document.getElementById('cookie-consent').classList.add('hidden');
  });
  
  document.getElementById('cookie-accept').addEventListener('click', function() {
    document.getElementById('cookie-consent').classList.add('hidden');
    document.cookie = "cookieAccepted=true; expires=Fri, 31 Dec 9999 23:59:59 GMT";
  });

  // Vérifier si l'utilisateur a déjà accepté les cookies
if (document.cookie.includes("cookieAccepted=true")) {
  document.getElementById('cookie-consent').classList.add('hidden');
} else {
  document.getElementById('cookie-consent').classList.remove('hidden');
}



function setCookie(name, value, days) {
  var expires = "";
  if (days) {
      var date = new Date();
      date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
      expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}


function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(';');
  for (var i = 0; i < ca.length; i++) {
      var c = ca[i];
      while (c.charAt(0) == ' ') {
          c = c.substring(1, c.length);
      }
      if (c.indexOf(nameEQ) == 0) {
          return c.substring(nameEQ.length, c.length);
      }
  }
  return null;
}


