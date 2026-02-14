package managers

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

/*
Cette fonction permet de connaître les commandes disponibles
@Param : conn net.Conn, args []string
@Return : erreur error
*/
func Help(conn net.Conn, args []string, isAdmin bool) (erreur error) {
	// Vérifier le nombre d'arguments
	if len(args) > 1 {
		_ = ConnexionSendText(conn, "[<ERROR>]Wrong number of arguments. Usage: help <command:*optionnal>")
		return errors.New("wrong number of arguments")
	}

	// Si le nombre d'arguments était bon alors faire :

	var commandList []string
	if isAdmin == false {
		commandList = []string{"list", "get", "end", "help"}
	} else {
		commandList = []string{"list", "end", "hide", "reveal", "terminate", "help"}
	}

	if len(args) == 0 {
		// Si aucune commande n'a été passée en argument

		err := ConnexionSendText(conn, "[<LISTER="+strconv.Itoa(len(commandList))+">]## Commands: ##")
		if err != nil {
			return err
		}

		for index, command := range commandList {
			err := ConnexionSendText(conn, strconv.Itoa(index+1)+" "+command)
			if err != nil {
				return err
			}
		}
	} else {
		// Rechercher si une commande correspond à l'argument donné
		var searchingCommand = strings.ToLower(args[0])
		for index, command := range commandList {
			if searchingCommand == command {
				// Commande trouvée
				var commandListHelp []string
				if isAdmin == false {
					commandListHelp = []string{
						"List every files in the filesStored folder.", //List
						"Get a certain file.",                         //Get
						"Close the connexion with the server.",        //End
						"This command let you learn details about commands, you can use it like this 'help <command>' but also without any arguments to list alls the documented commands.", //Help
					}
				} else {
					commandListHelp = []string{
						"List every files in the filesStored folder.",                                        //List
						"Close the connexion with the server.",                                               //End
						"Hide a file so nobody can see it anymore, you can use <<Reveal>> to show it again.", // Hide
						"If a file was hidden with <<Hide>> , you can now see it again",                      // Reveal
						"Close the server, not just the connexion with it.",                                  //Terminate
						"This command let you learn details about commands, you can use it like this 'help <command>' but also without any arguments to list alls the documented commands.", //Help
					}
				}

				// Envoyer le détail de la commande au client
				err := ConnexionSendText(conn, "## Command : "+strings.ToUpper(command[0:1])+command[1:]+" > "+commandListHelp[index])
				if err != nil {
					return err
				} else {
					return nil
				}
			}
		}
		// Aucune commande ne correspond
		err := ConnexionSendText(conn, "[<ERROR>]Help don't know anything about this command : "+args[0])
		return err
	}

	return nil
}
