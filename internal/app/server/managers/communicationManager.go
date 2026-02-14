// Author : Grégoire DELUGRE
// Date : 12/04/2025
// FileName : communicationManager.go

package managers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

/*
Cette fonction permet de recevoir une ligne de texte envoyé
@Param : conn net.Conn
@Return : string texte , erreur error
*/
func ConnexionReadText(conn net.Conn) (string, error) {
	// Timeout : au bout de 10 minutes d'inactivité
	err := conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	if err != nil {
		fmt.Println("\033[34mErreur pour la protection timeout\033[0m")
		return "", errors.New("Timeout")
	}

	// Pas besoin de timeout car attente potentielle d'une commande depuis l'appel de cette fonction
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

/*
Cette fonction permet de recevoir un fichier, d'obtenir son nom ainsi que son contenu
@Param : conn net.Conn, size int
@Return : string texte , erreur error
*/
func ConnexionReadFile(conn net.Conn, size int) ([]byte, string, error) {
	// Timeout : on laisse du temps au cas où on reçois un fichier énorme et/ou que la connexion n'est pas bonne
	err := conn.SetReadDeadline(time.Now().Add(1 * time.Hour))
	if err != nil {
		return []byte{}, "", errors.New("Timeout")
	}

	headerBuf := make([]byte, size)
	_, err = io.ReadFull(conn, headerBuf)
	if err != nil {
		return nil, "", err
	}

	header := string(headerBuf)
	header = strings.TrimSpace(header)

	if !strings.HasPrefix(header, "<FILE=") {
		return nil, "", fmt.Errorf("header invalide: %s", header)
	}

	// Exemple header : <FILE=120;readme.txt>
	inside := strings.TrimSuffix(strings.TrimPrefix(header, "<FILE="), ">")
	parts := strings.SplitN(inside, ";", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("header mal formé")
	}

	size, err = strconv.Atoi(parts[0])
	if err != nil {
		return nil, "", err
	}
	filename := parts[1]

	// Lire exactement 'size' octets pour le fichier
	buf := make([]byte, size)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, "", err
	}

	return buf, filename, nil
}

/*
Cette fonction permet d'envoyer une ligne de texte à une certaine connexion
@Param : conn net.Conn, text string
@Return : erreur error
*/
func ConnexionSendText(conn net.Conn, text string) error {
	// Timeout
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return errors.New("Timeout")
	}

	data := []byte(text + "\n")
	_, err = conn.Write(data)
	return err
}

/*
Cette fonction permet d'envoyer une ligne de texte à une certaine connexion
@Param : conn net.Conn
@Return : string texte , erreur error
*/
func ConnexionSendFile(conn net.Conn, filename string) error {
	// Timeout : on laisse du temps au cas où on reçois un fichier énorme et/ou que la connexion n'est pas bonne
	err := conn.SetReadDeadline(time.Now().Add(1 * time.Hour))
	if err != nil {
		return errors.New("Timeout")
	}

	// Récupérer le fichier a envoyer
	data, err := ReadFile(filename) // ( Fonction provenant de savingManager )
	if err != nil {
		return err
	}

	// Envoyer le fichier et renvoyer l'erreur
	_, err = conn.Write(data)
	return err
}
