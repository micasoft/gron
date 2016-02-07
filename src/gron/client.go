package gron

import (
	"log"
	"net"
)

func Client(cmd *string, prio *int) {
	c, err := net.Dial("unix", sock)
	bcmd := Job{RawCommand: *cmd, RawPrio: *prio}
	if err != nil {
		panic(err)
	}
	defer c.Close()
	log.Printf("run %s by %d", bcmd.RawCommand, bcmd.RawPrio)

	_, errz := c.Write(bcmd.Encode())
	if errz != nil {
		log.Fatal("write error:", err)
	}
}
