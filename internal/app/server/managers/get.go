package managers

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

/*
Cette fonction permet d'envoyer au client le fichier souhaité
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func Get(conn net.Conn, args []string) (erreur error) {
	// Vérifier le nombre d'arguments
	if len(args) != 1 {

		_ = ConnexionSendText(conn, "[<ERROR>]Wrong number of arguments. Usage: get <filename>")
		return errors.New("wrong number of arguments")
	}
	// Si le nombre d'arguments était bon alors faire :

	_, b := IsInHiddenFile(DefaultFilesDir + args[0])

	if b {
		// Vérifier si le fichier existe et est accessible
		// Il n'est pas accessible, prévenir le client
		_ = ConnexionSendText(conn, "[<ERROR>]FileUnknown")
		return errors.New("[<ERROR>]FileUnknown")
	}

	// Récupérer le fichier
	filename := strings.TrimSpace(args[0])
	file, err := ReadFile(DefaultFilesDir + filename)
	if err != nil { // Vérifier si le fichier existe et est accessible
		// Il n'est pas accessible, prévenir le client
		_ = ConnexionSendText(conn, "[<ERROR>]FileUnknown")
		return err
	}

	//Prévenir le joueur qu'on va lui envoyer un fichier
	err = ConnexionSendText(conn, "[<SAVEFILE="+strconv.Itoa(len(file))+"|"+filename+">]Start")
	if err == nil {
		err := ConnexionSendFile(conn, DefaultFilesDir+filename)
		if err != nil {
			return err
		}
		answer, err := ConnexionReadText(conn)
		if err != nil {
			return err
		}
		if strings.TrimSpace(answer) != "OK" {
			return errors.New("Waiting for message OK but got :" + answer)
		}
		return err
	} else {
		return err
	}
}
