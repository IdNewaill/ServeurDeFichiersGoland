package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gitlab.univ-nantes.fr/iutna.info2.r305/proj-groupe2/ric-del-pro/internal/app/client/managers"
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

// Dossier par défaut.
var DefaultFilesDir string

func InitialiseClient() {
	if abs, err := resolveDir("internal/app/client/filesStored/"); err == nil {
		// Ajoute un séparateur pour faciliter les concaténations "dir"+"filename"
		DefaultFilesDir = abs + string(os.PathSeparator)
	} else {
		// Fallback: garder le chemin relatif. Les fonctions appelantes géreront l'erreur si inexistante.
		DefaultFilesDir = "internal/app/client/filesStored/"
	}
	slog.Debug(DefaultFilesDir)
}

// -------------------------------------- > Main < --------------------------------------

func RunClient(remote string) {
	fmt.Println("\033[34mConnecting to the server : \033[0m" + remote)
	connexion, err := net.Dial("tcp", remote)
	InitialiseClient()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	managers.InitialiseLoggingManager("client")
	managers.LogAction("Client", "Conncted to the server : "+connexion.LocalAddr().String())
	defer connexion.Close()
	defer fmt.Println("\033[34m\nConnexion closed \033[0m")

	stdin := bufio.NewReader(os.Stdin)
	server := bufio.NewReader(connexion)
	fmt.Println("\033[32mConnected to the server !\033[0m\n")

	for {
		fmt.Print("> ")
		userInput, err := stdin.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		managers.LogAction("Client", "Execute > "+userInput)

		// Vérifier si l'utilisateur a demandé à partir
		if err != nil {
			if err == io.EOF {
				// Demande de partir interceptée (Ctrl-D)
				fmt.Println("\033[34mClosing Connexion without permissions from Server\033[0m")
				return
			}
		}

		// Vérifier si le message est vide
		if userInput == "" {
			continue
		}

		// Vérifier si l'utilisateur à lui même par la commande End demandé à fermer la connexion
		if strings.ToLower(userInput) == "end" {
			fmt.Println("\033[34mClosing Connexion without permissions from Server\033[0m")
			return
		}

		// Timeout
		err = connexion.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Println("\033[34mErreur pour la protection timeout\033[0m")
			return
		}

		// Envoi de la commande en texte
		_, err = connexion.Write([]byte(userInput + "\n"))
		if err != nil {
			slog.Error("Server unreachable:", "err", err)
			return
		}

		// Lire une réponse ligne par ligne
		line, err := server.ReadString('\n')
		if err != nil {
			slog.Warn("Déconnexion du serveur")
			return
		}
		line = strings.TrimSpace(line)

		// Types de requête

		// QUIT
		if line == "[<QUIT>]" {
			slog.Debug("Serveur demande la fermeture.")
			return
		}

		// LIST=N
		if strings.HasPrefix(line, "[<LISTER=") {
			handleListResponse(line, server)
			continue
		}

		// ERROR
		if strings.HasPrefix(line, "[<ERROR>]") {
			text, _ := strings.CutPrefix(line, "[<ERROR>]")
			fmt.Println("\033[31m" + text + "\033[0m")
			continue
		}

		// Save File
		if strings.HasPrefix(line, "[<SAVEFILE=") {
			err := handleSaveFileResponse(connexion, line)
			if err != nil {
				fmt.Println("\033[31m" + "Server seems have issues or is not compatible ! " + "\033[0m")
				return
			}
			continue
		}

		// Par défault, écrire ce texte
		fmt.Println(line)
	}
}

// -------------------------------------- > Simples functions < --------------------------------------

func getHeaderData(headerNeeded string, line string) (headerValue string, sucess bool) {
	// Exemple header: "[<LISTER=5>]"
	// Autre Exemple header : "[<LISTER=1>][[ -------- FileCount 1 -------- ]]"
	start := strings.Index(line, "[<"+headerNeeded+"=")
	end := strings.Index(line, ">")
	if start == -1 || end == -1 {
		return "", false
	}
	if end != len(line) {
		fmt.Println(line[end+2:])
	}
	inside := line[start+len("[<"+headerNeeded+"=") : end]
	return inside, true
}

func convSizeToKo(size int) int {
	ko := float64(size) / 1024.0
	return int(math.Ceil(ko))
}

// -------------------------------------- > Handles Command Functions < --------------------------------------

func handleListResponse(line string, server *bufio.Reader) {
	// Récupérer le d'éléments à afficher
	headerValue, success := getHeaderData("LISTER", line)
	if success == false {
		return
	}

	// Vérifier que le nombre d'éléments passé en header est bien un entier positif
	count, err := strconv.Atoi(headerValue)
	if err != nil || count < 0 {
		slog.Warn("Format LIST reçu invalide:", "header", line)
		return
	}

	// Reception puis affichage des éléments (fait count fois)
	for i := 0; i < count; i++ {
		line, err := server.ReadString('\n')
		if err != nil {
			slog.Warn("Déconnexion pendant LIST")
			return
		}
		fmt.Print(line) // déjà avec \n
	}
}

func handleSaveFileResponse(conn net.Conn, line string) error {
	// Header look like ["[<SAVEFILE=SIZE|NAME]"]
	headerValue, success := getHeaderData("SAVEFILE", line)
	if success == false {
		return os.ErrExist
	}
	splitList := strings.Split(headerValue, "|")
	if len(splitList) != 2 {
		fmt.Println("\033[31m" + "Wrong received format answer" + "\033[0m") // Ecrire en Rouge
		return os.ErrExist
	}

	size, erreur := strconv.Atoi(splitList[0])
	if erreur != nil {
		fmt.Println("\033[31m" + "Wrong received format answer" + "\033[0m") // Ecrire en Rouge
		return os.ErrExist
	}

	fmt.Println("Téléchargement de " + strconv.Itoa(convSizeToKo(size)) + "ko")
	buf := make([]byte, size)

	// Timeout
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		fmt.Println("\033[34mErreur pour la protection timeout\033[0m")
		return errors.New("Timeout")
	}

	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}

	success = managers.SaveFile(DefaultFilesDir+splitList[1], buf)

	if success == true {
		fmt.Println("\033[34mSaved at " + DefaultFilesDir + splitList[1] + "\033[0m")
	}

	// Répondre OK (peut importe si la sauvegarde a été réussie ou non)
	_, _ = conn.Write([]byte("OK\n"))
	return nil
}
