package managers

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

/*
Cette fonction envoie au client la liste des fichiers disponibles dans dirPath.
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func List(c net.Conn, args []string) error {
	var dirPath string

	// Si un argument est fourni, l'utiliser comme sous-dossier
	if len(args) == 0 || args[0] == "" {
		dirPath = DefaultFilesDir
	} else {
		dirPath = filepath.Join(DefaultFilesDir, args[0])
	}

	// Vérifier si cette commande part de la racine dossier de stockage
	success := IsPathInsideBase(DefaultFilesDir, dirPath)
	if !success {
		slog.Error("Tentative de sortir de la racine", "dir", dirPath)
		ConnexionSendText(c, "[<ERROR>]Can't go outside root")
		return errors.New("[<ERROR>]Tentative de sortir de la racine")
	}
	_, isIn := IsInHiddenFile(dirPath)
	if isIn {
		slog.Error("Hidden Folder")
		ConnexionSendText(c, "[<ERROR>]Unable to read directory")
		return errors.New("[<ERROR>]Unable to read directory because this folder is hidden")
	}

	resolvedDir, resolveErr := resolveDir(dirPath)
	if resolveErr != nil {
		slog.Error("Unable to resolve directory", "dir", dirPath, "err", resolveErr)
		ConnexionSendText(c, "[<ERROR>]Unable to read directory")
		return resolveErr
	}

	// Lecture sécurisée du dossier cible (résolu en chemin absolu)
	entries, err := os.ReadDir(resolvedDir)
	if err != nil {
		slog.Error("Unable to read directory", "dir", resolvedDir, "err", err)
		ConnexionSendText(c, "[<ERROR>]Unable to read directory")
		return err
	}

	var NewEntries = []os.DirEntry{}

	for i := 0; i < len(entries); i++ {
		_, b := IsInHiddenFile(dirPath + entries[i].Name())
		if !b {
			NewEntries = append(NewEntries, entries[i])
		}

	}
	type fileData struct {
		name string
		size int64
	}

	files := make([]fileData, 0, len(NewEntries))
	for _, entry := range NewEntries {
		info, infoErr := entry.Info()
		if entry.IsDir() {
			files = append(files, fileData{name: entry.Name(), size: -1})
			continue
		}

		if infoErr != nil {
			slog.Warn("List: unable to stat file", "file", entry.Name(), "err", infoErr)
			continue
		}
		files = append(files, fileData{name: entry.Name(), size: info.Size()})
	}

	// Envoi du header FileCnt puis d'une ligne par fichier
	ConnexionSendText(c, fmt.Sprintf("[<LISTER="+strconv.Itoa(len(NewEntries))+">]FileCnt %d", len(files)))
	for _, f := range files {
		if f.size == -1 { // Il s'agit d'un dossier
			ConnexionSendText(c, fmt.Sprintf("%s folder", f.name))
		} else {
			ConnexionSendText(c, fmt.Sprintf("%s %d", f.name, f.size))
		}
	}

	return nil
}
