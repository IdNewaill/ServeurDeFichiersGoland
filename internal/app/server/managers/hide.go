package managers

import (
	"errors"

	"net"
)

/*
Cette fonction permet de cacher un fichier aux autres (momentanément)
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func Hide(conn net.Conn, args []string) (erreur error) {
	// Vérifier le nombre d'arguments
	if len(args) != 1 {

		_ = ConnexionSendText(conn, "[<ERROR>]Wrong number of arguments. Usage: hide <filename>")
		return errors.New("wrong number of arguments")
	}
	// Si le nombre d'arguments était bon alors faire :

	HiddenFileAdd(conn, DefaultFilesDir+args[0])

	return nil
}
