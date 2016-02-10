package main

import (
	"flag"
	"gron"
)

var daemon = flag.Bool("d", false, "Starting as a daemon")
var status = flag.Bool("status", false, "Maximun of process")
var cmd = flag.String("c", "ls", "Command")
var max = flag.Int("max", 10, "Maximun of process")
var prio = flag.Int("p", 0, "Set a prio of process")

func main() {
	//Read config
	flag.Parse()

	switch {
	case *daemon:
		gron.Server(*max)
		break
	case *status:
		gron.Status()
		break
	default:
		gron.Client(cmd, prio)
	}
}
