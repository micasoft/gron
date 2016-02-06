package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const sock = "/tmp/gron.sock"

var daemon = flag.Bool("d", false, "Starting as a daemon")
var cmd = flag.String("c", "ls", "Command")
var prio = flag.Int("p", 0, "Set a prio of process")

type Job struct {
	RawCommand string
	RawPrio    int
}

func (j *Job) Encode() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	encoder.Encode(j.RawCommand)
	encoder.Encode(j.RawPrio)
	return w.Bytes()
}

func (j *Job) Decode(buf []byte) {
	r := bytes.NewBuffer(buf)
	decoder := gob.NewDecoder(r)
	decoder.Decode(&j.RawCommand)
	decoder.Decode(&j.RawPrio)
}

func server() {
	l, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	log.Printf("Starting gron as a daemon")

	ksignal := make(chan os.Signal, 1)

	//Finish the application gracefully
	signal.Notify(ksignal, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go execute(fd)
		}
	}()

	kv := <-ksignal
	log.Printf("Signal to finish : %s", kv.String())
}

func execute(c net.Conn) {
	for {
		buf := make([]byte, 1024)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:nr]
		j := new(Job)
		j.Decode(data)

		log.Printf("let me run this : %s with prio %d", j.RawCommand, j.RawPrio)
		//_, err = c.rite(data)
	}
}

func client(cmd *string, prio *int) {
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

func main() {
	//Read config
	flag.Parse()

	if *daemon {
		server()
	} else {
		client(cmd, prio)
	}
}
