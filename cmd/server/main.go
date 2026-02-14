package main

import (
	"flag"
	"log/slog"

	"gitlab.univ-nantes.fr/iutna.info2.r305/proj-groupe2/ric-del-pro/internal/app/server"
)

func parseArgs() (port *string, portAdmin *string) {

	logLevel := flag.Bool("d", true, "enable debug log level")
	port = flag.String("p", "3333", "server port (default: 3333)")
	portAdmin = flag.String("a", "8080", "server port (default: 8080)")

	flag.Parse()

	if *logLevel {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Set logging level to debug")
	}

	return
}

func main() {
	port, portAdmin := parseArgs()
	server.RunServer(port, portAdmin)
}
