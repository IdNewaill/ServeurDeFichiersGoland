package managers

import (
	"net"
	"os"
	"path/filepath"
	"sync"
)

var HiddenFileLock sync.Mutex
var ListHiddenFile []string = []string{}

func IsInHiddenFile(file string) (int, bool) {
	file = filepath.Clean(file)
	HiddenFileLock.Lock()
	defer HiddenFileLock.Unlock()

	for i, s := range ListHiddenFile {
		theRessource, err := os.Stat(s)
		if err != nil {
			continue
		}
		if s == file || (theRessource.IsDir() && IsPathInsideBase(s, file)) {
			return i, true
		}
	}
	return -1, false
}

func HiddenFileAdd(conn net.Conn, file string) {
	file = filepath.Clean(file)
	_, b := IsInHiddenFile(file)

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		ConnexionSendText(conn, "La ressource n'existe pas")
		return
	}

	if !b {
		HiddenFileLock.Lock()
		defer HiddenFileLock.Unlock()
		ListHiddenFile = append(ListHiddenFile, file)
		ConnexionSendText(conn, "La ressource a été caché")
	} else {
		ConnexionSendText(conn, "La ressource est déjà caché")
	}

}

func HiddenFileRemove(conn net.Conn, file string) {
	file = filepath.Clean(file)
	i, b := IsInHiddenFile(file)

	if b {
		HiddenFileLock.Lock()
		defer HiddenFileLock.Unlock()
		ListHiddenFile = append(ListHiddenFile[:i], ListHiddenFile[i+1:]...)
		ConnexionSendText(conn, "La ressource est de retour !")
	} else {
		ConnexionSendText(conn, "La ressource est déjà visible")
	}

}
