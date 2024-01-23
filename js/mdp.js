
var input = document.getElementById("password");

input.addEventListener("input", function () {
    var value = input.value;
    var uppercaseCount = 0;
    var specialCharacterCount = (value.match(/[!@#$%&*,.?]/g) || []).length;

    for (var i = 0; i < value.length; i++) {
        if (value[i] === value[i].toUpperCase()) {
            uppercaseCount++;
        }
    }

    if (uppercaseCount < 1 || specialCharacterCount < 1) {
        input.setCustomValidity("Le champ doit contenir au moins: \n - 1 lettre majuscule \n - 1 caractère spécial. (!@#$%&*,.?)");
    } else {
        input.setCustomValidity("");
    }
});
