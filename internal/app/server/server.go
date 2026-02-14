package server

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"gitlab.univ-nantes.fr/iutna.info2.r305/proj-groupe2/ric-del-pro/internal/app/server/managers"
	"gitlab.univ-nantes.fr/iutna.info2.r305/proj-groupe2/ric-del-pro/internal/app/server/managers/commun"
)

// -------------------------------------- > Main < --------------------------------------

func RunServer(port *string, portAdmin *string) {
	// Initialisation des modules
	fmt.Println("\033[34m[Starting Server on port\033[0m " + *port + "\033[34m ]\033[0m")
	commun.InitialiseLoggingManager("server")
	managers.InitialiseExecutor()

	// Démarrer le serveur
	l, e := net.Listen("tcp", ":"+*port)
	if e != nil {
		slog.Error(e.Error())
		fmt.Println("Le serveur n'a pas réussi à démarrer", e.Error())
		return
	}

	lAdmin, eAdmin := net.Listen("tcp", ":"+*portAdmin)
	if eAdmin != nil {
		slog.Error(eAdmin.Error())
		return
	}

	go listenToAdminPort(lAdmin) // Lancer une go routine qui s'occupe du client Administrateur

	fmt.Println("\033[32mServer has started ! Listening for clients ..\033[0m\n")

	// A l'arrêt du serveur
	defer func() {
		commun.CloseLoggingManager()
		l.Close()
		lAdmin.Close()
		slog.Debug("Stopped listening on port " + *port)
	}()

	//Traitement lorsqu'un client se connecte

	go func() {
		slog.Debug("Now listening on port " + *port)
		for {
			c, e := l.Accept()
			if e != nil {
				slog.Error(e.Error())
				continue
			}
			managers.AddClient(c, false)
		}

	}()

	for managers.IsServerClosed() == false {
		time.Sleep(1 * time.Second)
	}

}

// -------------------------------------- > Fonctions < --------------------------------------

func listenToAdminPort(l net.Listener) {
	for {
		c, e := l.Accept()
		if e != nil {
			slog.Error(e.Error())
			continue
		}

		managers.AddClient(c, true)
	}
}
