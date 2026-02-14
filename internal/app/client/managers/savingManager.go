// Author : Grégoire DELUGRE
// Date : 12/04/2025
// FileName : savingManager.go

package managers

import (
	"os"
	"sync"
)

// -------------------------------------- > Fonctions anti Races conditions < --------------------------------------

type FileLockManager struct {
	filePath string
	lock     *sync.Mutex
}

var lockFolderLocks sync.Mutex
var folderLocksStrucsList []*FileLockManager = make([]*FileLockManager, 0)

/*
Cette fonction permet d'obtenir un lock lié à un fichier
NE PAS UTILISER, est liée seulement au fichier savingManager
@Param : filePath string
@Return : *sync.Mutex
*/
func getFolderLock(filePath string) *sync.Mutex {
	lockFolderLocks.Lock()
	defer lockFolderLocks.Unlock()

	// Rechercher le Lock dans la liste folderLocksStrucsList<FileLockManager>
	for _, fileLock := range folderLocksStrucsList {
		if fileLock.filePath == filePath {
			return fileLock.lock
		}
	}

	// Créer un lock puisqu'il n'est pas déjà présent dans la liste

	var lock = &FileLockManager{
		filePath: filePath,
		lock:     &sync.Mutex{},
	}

	folderLocksStrucsList = append(folderLocksStrucsList, lock)

	return lock.lock
}

// -------------------------------------- > Fonctions de sauvegardes < --------------------------------------

/*
Cette fonction permet de sauvegarder un fichier au path renseigné.
@Param : filePath string
*/
func SaveFile(filePath string, data []byte) bool {
	// Empêcher l'écriture et/ou la lecture multiple d'un même fichier
	var actualLock *sync.Mutex = getFolderLock(filePath)
	actualLock.Lock()
	defer actualLock.Unlock()

	// Sauvegarder le fichier
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		print("[ savingManager ] Une erreur est survenue lors de la création d'un fichier '" + filePath + "' voici l'erreur > " + err.Error())
		LogAction("Client", "Tryed to save File at '"+filePath+"' but an error occured : "+err.Error())
	} else {
		LogAction("Client", "File at '"+filePath+"' has been saved")
	}
	return err == nil
}

/*
Cette fonction permet de lire un fichier au path renseigné.
@Param : filePath string
@Return : file []byte , erreur error
*/
func ReadFile(filePath string) ([]byte, error) {
	// Empêcher l'écriture et/ou la lecture multiple d'un même fichier
	var actualLock *sync.Mutex = getFolderLock(filePath)
	actualLock.Lock()
	defer actualLock.Unlock()

	// Lire le fichier
	res, err := os.ReadFile(filePath)
	if err != nil {
		LogAction("Client", "File at '"+filePath+"' has been read")
	} else {
		LogAction("Client", "Try to read File at '"+filePath+"' but an error has occured : "+string(res))
	}
	return res, err
}
