package main

import (
	"flag"
	"gron"
)

var daemon = flag.Bool("daemon", false, "Starting gron as a daemon")
var status = flag.Bool("status", false, "Status of daemon")
var cmd = flag.String("command", "", "Commmand in bash that will be run")
var max = flag.Int("max", 10, "Maximum of process can be executed at same time")
var prio = flag.Int("p", 0, "Set a priority at process level")
var last = flag.Bool("last", false, "Last command executed")

func init() {
	flag.BoolVar(daemon, "d", false, "")
	flag.BoolVar(status, "s", false, "")
	flag.StringVar(cmd, "c", "", "")
	flag.BoolVar(last, "l", false, "")
}

func main() {
	//Read config
	flag.Parse()

	switch {
	case *daemon:
		gron.Server(*max)
		break
	case *status:
		gron.GetStatus(last)
		break
	default:
		gron.Client(cmd, prio)
	}
}
