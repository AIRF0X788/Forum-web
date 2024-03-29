# Forum Web en Go

Ce programme en Go implémente un forum web simple avec les fonctionnalités d'inscription, de connexion, de création de sujets, de visualisation de sujets, et d'administration d'utilisateurs.

## Dépendances
Ce projet utilise les packages suivants:

database/sql pour l'accès à la base de données SQLite.
github.com/mattn/go-sqlite3 pour le pilote SQLite.
github.com/google/uuid pour la génération d'identifiants uniques.
golang.org/x/crypto/bcrypt pour le hachage des mots de passe.

## Structure de la base de données
Le programme utilise une base de données SQLite avec trois tables:

users pour stocker les informations des utilisateurs.
topics pour stocker les sujets du forum.
messages pour stocker les messages associés aux sujets.

## Installation et Configuration
Installez les dépendances en utilisant go get.

go get github.com/mattn/go-sqlite3
go get github.com/google/uuid
go get golang.org/x/crypto/bcrypt
Assurez-vous que vous avez une base de données SQLite nommée database.db dans le répertoire du programme.

Exécutez le programme avec go run main.go.

Accédez au forum via votre navigateur à l'adresse http://localhost:8080.

## Fonctionnalités Principales

### Inscription et Connexion

Accédez à la page d'inscription via /signup et connectez-vous via /login.
Validation des champs, gestion des erreurs.

### Création de Sujets

Utilisez la page /topics pour créer de nouveaux sujets.
Possibilité de joindre une image au sujet.

### Visualisation de Sujets

Consultez la liste des sujets sur la page d'accueil (/) et accédez aux détails d'un sujet via /topic/{id}.

### Administration d'Utilisateurs

Accédez à la page d'administration via /admin (réservée aux administrateurs).
Bannissez des utilisateurs via l'interface d'administration.
