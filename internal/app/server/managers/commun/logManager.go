// Author : Grégoire DELUGRE
// Date : 12/04/2025
// FileName : logManager.go

package commun

import (
	"log"
	"os"
	"time"
)

var logFile *os.File

/*
Cette fonction est à appeler avant même d'utiliser le logAction.
@Param : None
*/
func InitialiseLoggingManager(nameFile string) {
	var erreur error
	nameFile = nameFile + ".log"
	logFile, erreur = os.OpenFile(nameFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if erreur != nil {
		panic("[ logManager ] Une erreur est survenue, la tentative de création/ouverture du fichier log a échouée, voici l'erreur > " + erreur.Error())
	}

	log.SetOutput(logFile)
}

/*
Cette fonction permet de remplir le fichier log.
@Param : user string, command string
*/
func LogAction(user string, command string) {
	var dateToUse = time.Now()
	var res = dateToUse.String() + " / " + user + " / " + command
	log.Println(res)
}

func CloseLoggingManager() {
	if logFile != nil {
		logFile.Close()
	}
}
