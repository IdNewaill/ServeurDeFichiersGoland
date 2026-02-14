// Author : Grégoire DELUGRE
// Date : 12/06/2025
// FileName : executor.go

package managers

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// resolveDir trouve un chemin absolu vers dirPath, même si le binaire
// est lancé depuis un sous-dossier (cmd/server, go run, etc.).
func resolveDir(dirPath string) (string, error) {
	if filepath.IsAbs(dirPath) {
		return dirPath, nil
	}

	tryPaths := make([]string, 0, 6)

	if cwd, err := os.Getwd(); err == nil {
		tryPaths = append(tryPaths, filepath.Join(cwd, dirPath))
		parent := cwd
		for i := 0; i < 3; i++ {
			parent = filepath.Dir(parent)
			tryPaths = append(tryPaths, filepath.Join(parent, dirPath))
		}
	}

	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		tryPaths = append(tryPaths, filepath.Join(execDir, dirPath))
		parent := execDir
		for i := 0; i < 3; i++ {
			parent = filepath.Dir(parent)
			tryPaths = append(tryPaths, filepath.Join(parent, dirPath))
		}
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		srcDir := filepath.Dir(file)
		srcRoot := filepath.Clean(filepath.Join(srcDir, "..", "..", "..", ".."))
		tryPaths = append(tryPaths, filepath.Join(srcRoot, dirPath))
	}

	seen := make(map[string]struct{})
	for _, p := range tryPaths {
		p = filepath.Clean(p)
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		if info, err := os.Stat(p); err == nil && info.IsDir() {
			return p, nil
		}
	}

	return "", fmt.Errorf("cannot resolve directory %s", dirPath)
}

// Dossier par défaut, résolu au démarrage.
var DefaultFilesDir string

func InitialiseExecutor() {
	if abs, err := resolveDir("internal/app/server/filesStored/"); err == nil {
		// Ajoute un séparateur pour faciliter les concaténations "dir"+"filename"
		DefaultFilesDir = abs + string(os.PathSeparator)
	} else {
		// Fallback: garder le chemin relatif. Les fonctions appelantes géreront l'erreur si inexistante.
		DefaultFilesDir = "internal/app/server/filesStored/"
	}
	slog.Debug(DefaultFilesDir)
}

/*
Cette fonction permet d'executer une commande
@Param : conn net.Conn, commande string
@Return : erreur error
*/
func ExecuteCommand(conn net.Conn, textCommand string, admin bool) error {
	var cut = strings.Split(textCommand, " ")
	var command = strings.ToLower(cut[0])
	var args = cut[1:]
	var err error
	if admin {
		switch command {
		case "end":
			err = End(conn, args)
		case "list":
			err = List(conn, args)
		case "terminate":
			err = Terminate(conn, args)
		case "hide":
			err = Hide(conn, args)
		case "reveal":
			err = Reveal(conn, args)
		case "help": // Commande ajoutée
			err = Help(conn, args, true)
		default:
			// Commande inconnue
			_ = ConnexionSendText(conn, "[<ERROR>]Unknown action")
			err = errors.New("Unknown command")
		}
	} else {
		switch command {
		case "end":
			err = End(conn, args)
		case "list":
			err = List(conn, args)
		case "get":
			err = Get(conn, args)
		case "help": // Commande ajoutée
			err = Help(conn, args, false)
		default:
			// Commande inconnue
			_ = ConnexionSendText(conn, "[<ERROR>]Unknown action")
			err = errors.New("Unknown command")
		}
	}
	return err
}
