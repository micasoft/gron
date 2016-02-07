package main

import (
	"flag"
	"gron"
)

var daemon = flag.Bool("d", false, "Starting as a daemon")
var cmd = flag.String("c", "ls", "Command")
var prio = flag.Int("p", 0, "Set a prio of process")

func main() {
	//Read config
	flag.Parse()

	if *daemon {
		gron.Server()
	} else {
		gron.Client(cmd, prio)
	}
}
