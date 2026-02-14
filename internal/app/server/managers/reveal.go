package managers

import (
	"errors"
	"net"
)

/*
Cette fonction execute l'action de rendre visible un fichier
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func Reveal(conn net.Conn, args []string) (erreur error) {
	// Vérifier le nombre d'arguments
	if len(args) != 1 {

		_ = ConnexionSendText(conn, "[<ERROR>]Wrong number of arguments. Usage: reveal <filename>")
		return errors.New("wrong number of arguments")
	}
	// Si le nombre d'arguments était bon alors faire :
	HiddenFileRemove(conn, DefaultFilesDir+args[0])

	return nil
}
