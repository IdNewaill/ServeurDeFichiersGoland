En envoyant un message à un client, vous pouvez utiliser ces balises en préfix pour effectuer différentes actions :

- >[\<ERROR>]Message

    Permet d'indiquer que le message qui suit est une erreur. Le message sera donc affiché en rouge.

- >[\<LISTER=**N**>]Message

    Permet de lister **N** informations les unes en dessous des autres. Cette valeur **N** étant donc un nombre et donc, l'utilisateur attendra **N** messages (avec \n pour la fin d'un message) de la part du serveur.

    Il faut savoir que << Message >> dans la commande est optionnel, mais l'écrire garantira que ce Message soit bien écrit en première position.
- >[\<SAVEFILE=**Size**|**AtPath**>]

    Permet de faire sauvegarder au client dans le répertoire << fileStored+**AtPath** >> le fichier qui sera envoyé juste après au serveur.
    Il ne faudra pas oublier de prévenir le serveur que le fichier a bien été reçu avec le mot clé "OK" que le client envoie après avoir reçu le fichier !
- >[\<QUIT>]
   
    Permet de dire au client que le serveur l'a expulsé et donc le client sait qu'il peut arrêter son programme.
