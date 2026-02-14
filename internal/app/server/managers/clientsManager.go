// Author : Grégoire DELUGRE
// Date : 12/06/2025
// FileName : clientsManager.go

package managers

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	"gitlab.univ-nantes.fr/iutna.info2.r305/proj-groupe2/ric-del-pro/internal/app/server/managers/commun"
)

// [[ Variables ]]

var clients = make([]net.Conn, 0) // liste des clients connectés
var clientsLock sync.Mutex
var clientsCount int = 0

// [[ Functions ]]

/*
Cette fonction permet d'ajouter un client au serveur
@Param : conn net.Conn
*/
func AddClient(conn net.Conn, admin bool) {
	clientsLock.Lock()
	if admin == false {
		defer clientsLock.Unlock()
	}

	// Ajouter le client à la liste des clients présents
	clientsCount++
	clients = append(clients, conn)
	fmt.Println("\n\033[35mClient Connected> "+conn.RemoteAddr().String()+"\nNumber of clients >", clientsCount, "\n\033[0m")

	// Lancer une go routine pour suivre les demandes du client
	if admin == true {
		// Puisqu'un seul admin peut se connecter en même temps, on execute sans go routine
		clientsLock.Unlock()
		HandleClient(conn, admin)
	} else {
		go HandleClient(conn, admin)
	}
}

/*
Cette fonction permet d'expulser/retirer un client
@Param : conn net.Conn
*/
func RemoveClient(conn net.Conn) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i+1:]...)
			clientsCount--
			fmt.Println("\n\033[35mClient Deconnected> "+conn.RemoteAddr().String()+"\nNumber of clients >", clientsCount, "\n\033[0m")
			commun.LogAction(conn.RemoteAddr().String(), "Removed this client")
			_ = ConnexionSendText(conn, "[<QUIT>]")
			return
		}
	}
}

/*
Cette fonction permet de retirer tous les clients
@Param : None
PS : RICHARD Baptiste
*/

func RemoveAllClients() {

	for _, c := range clients {
		RemoveClient(c)
	}

}

/*
Cette permet de savoir si un client se trouve toujours dans la liste des clients présents
@Param : conn net.Conn
*/
func clientNeedToLogOut(conn net.Conn) bool {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	for _, client := range clients {
		if client == conn {
			return false
		}
	}
	return true
}

/*
Cette fonction permet de gérer un client et lui permettre d'executer des commandes
@Param : conn net.Conn
*/
func HandleClient(conn net.Conn, admin bool) {
	for {
		// Vérifier que le client n'a pas été expulsé
		if clientNeedToLogOut(conn) {
			return
		}

		// voir si une commande a besoin d'être executé
		command, err := ConnexionReadText(conn)
		if err != nil {
			// Le joueur ne semble plus présent ou a rencontré une erreur
			RemoveClient(conn)
		} else {
			// Une commande a bien été reçue donc il faut executer cette commande
			commun.LogAction(conn.RemoteAddr().String(), "Has entered this command : "+command)
			slog.Debug("Cette commande a été reçue : " + command)

			err := ExecuteCommand(conn, command, admin)
			if err == nil {
				slog.Debug("\033[32mSuccès\033[0m lors de l'execution de cette commande : " + command)
			} else {
				slog.Debug("\033[33mErreur lors de l'execution de cette commande : " + command + "\n" + err.Error() + "\033[0m")
			}
		}
	}
}
