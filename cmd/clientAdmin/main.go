package main

import (
	"flag"
	"log/slog"

	"gitlab.univ-nantes.fr/iutna.info2.r305/proj-groupe2/ric-del-pro/internal/app/client"
)

func parseArgs() (remote string) {
	dFlag := flag.Bool("d", false, "enable debug log level")
	aFlag := flag.String("a", "127.0.0.1", "server address (default: 127.0.0.1)")
	pFlag := flag.String("p", "8080", "server port (default: 8080)")
	flag.Parse()

	if *dFlag {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	remote = *aFlag + ":" + *pFlag
	return
}

func main() {
	remote := parseArgs()
	client.RunClient(remote)
}
