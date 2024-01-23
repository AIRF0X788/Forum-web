$(document).ready(function () {
    // Gérer la soumission du formulaire
    $('#messageForm').submit(function (event) {
        event.preventDefault();

        // Récupérer le contenu du message saisi par l'utilisateur
        const message = $('#messageInput').val();

        // Envoyer le message au serveur via une requête AJAX
        $.ajax({
            type: 'POST',
            url: '/enregistrer-message',
            data: JSON.stringify({ contenu: message }),
            contentType: 'application/json',
            success: function (response) {
                // Actualiser la page
                location.reload();
            },
            error: function (error) {
                console.error('Erreur lors de l\'enregistrement du message :', error);
            }
        });

        $('#messageInput').val('');
    });
});