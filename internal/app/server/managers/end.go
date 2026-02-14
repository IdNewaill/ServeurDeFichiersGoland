package managers

import (
	"net"
)

/*
Cette fonction permet de fermer la connexion avec un client quand il nous le demande.
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func End(conn net.Conn, args []string) error {
	RemoveClient(conn)
	return nil
}
