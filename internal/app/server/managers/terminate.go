package managers

import (
	"net"
	"sync"
)

var terminateLock sync.Mutex
var ServerIsClosing bool = false

/*
Cette fonction permet de cacher un fichier aux autres (momentanément)
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func Terminate(conn net.Conn, args []string) (erreur error) {
	terminateLock.Lock()
	defer terminateLock.Unlock()
	if !ServerIsClosing {
		ServerIsClosing = true
		RemoveAllClients()
	}

	return nil
}

func IsServerClosed() bool {
	terminateLock.Lock()
	defer terminateLock.Unlock()
	return ServerIsClosing
}
