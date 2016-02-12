package main

import (
	"flag"
	"gron"
)

var daemon = flag.Bool("d", false, "Starting gron as a daemon")
var status = flag.Bool("status", false, "Status of daemon")
var cmd = flag.String("c", "ls", "Commmand in bash that will be run")
var max = flag.Int("max", 10, "Maximum of process can be executed at same time")
var prio = flag.Int("p", 0, "Set a priority at process level")

func main() {
	//Read config
	flag.Parse()

	switch {
	case *daemon:
		gron.Server(*max)
		break
	case *status:
		gron.GetStatus()
		break
	default:
		gron.Client(cmd, prio)
	}
}
