function banUser(userID) {
    // Envoyer une requête POST au serveur pour bannir l'utilisateur avec l'ID spécifié
    fetch("/ban", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ userID: userID })
    })
    .then(response => {
        if (response.ok) {
            alert("Utilisateur banni avec succès !");
            location.reload();
        } else {
            alert("Une erreur est survenue lors du bannissement de l'utilisateur.");
        }
        
    })
    .catch(error => {
        alert("Une erreur est survenue lors du bannissement de l'utilisateur.");
    });
    
}

function unbanUser(userID) {
    // Envoyer une requête POST au serveur pour bannir l'utilisateur avec l'ID spécifié
    fetch("/unban", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ userID: userID })
    })
    .then(response => {
        if (response.ok) {
            alert("Utilisateur débanni avec succès !");
            location.reload();
        } else {
            alert("Une erreur est survenue lors du débannissement de l'utilisateur.");
        }
    })
    .catch(error => {
        alert("Une erreur est survenue lors du débannissement de l'utilisateur.");
    });
}