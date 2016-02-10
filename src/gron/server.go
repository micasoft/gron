package gron

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

const sock = "/tmp/gron.sock"
const sh_cmd = "/bin/bash"

var jobs []*Job
var sem *Semaphore
var sequence int

type Job struct {
	Sequence   int
	Created    time.Time
	RawCommand string
	RawPrio    int
	Prio       int
	ExitStatus int
	Took       time.Duration
}

type JobsSorter []*Job

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

func (j *Job) Parse() (string, []string) {
	var r []string
	for _, s := range strings.Split(j.RawCommand, " ") {
		r = append(r, strings.Trim(s, " "))
	}
	return r[0], r[1:]
}

// Len is part of sort.Interface.
// Len is the number of elements in the collection.
func (js JobsSorter) Len() int {
	return len(js)
}

// Swap is part of sort.Interface.
// Swap swaps the elements with indexes i and j.
func (js JobsSorter) Swap(i, j int) {
	js[i], js[j] = js[j], js[i]
}

// Less is part of sort.Interface.
// Less reports whether the element with
// index i should sort before the element with index j.
func (js JobsSorter) Less(i, j int) bool {
	return js[i].Prio < js[j].Prio
}

func Stats() {
	log.Printf("running: %d waiting: %d", sem.Available(), len(jobs))
}

func Server(max int) {
	l, err := net.Listen("unix", sock)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	gob.Register(Job{})
	log.Printf("Starting as a daemon with %d max of process", max)
	sequence = 0
	ksignal := make(chan os.Signal, 1)
	//Finish the application gracefully
	signal.Notify(ksignal, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for {
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go trait_request(fd)

		}
	}()
	sem = NewSemaphore(max)
	go engine()

	kv := <-ksignal
	log.Printf("Signal to finish : %s", kv.String())
	os.Remove(sock)
}

func trait_request(c net.Conn) {
	var data []byte
	defer c.Close()
	for {
		buf := make([]byte, 1024)
		nr, err := c.Read(buf)
		if err != nil {
			break
		}
		data = append(data, buf[0:nr]...)
	}
	sequence = sequence + 1
	cr := new(ClientRequest)
	cr.Decode(data)
	switch cr.Request {
	case "job":
		cr.Job.Sequence = sequence
		cr.Job.Created = time.Now()
		cr.Job.Prio = cr.Job.RawPrio
		jobs = append(jobs, &cr.Job)
		break
	case "status":
		d := []byte("Not implemented yet")
		c.Write(d)
		break
	}
}

func engine() {
	var j *Job
	for {
		if len(jobs) > 0 {
			sem.Acquire()
			//should sort the slice
			sort.Sort(JobsSorter(jobs))
			Stats()
			//pop a job
			j, jobs = jobs[len(jobs)-1], jobs[:len(jobs)-1]
			go execute(j)
		}
	}
}

func execute(j *Job) {
	var waitStatus syscall.WaitStatus
	cmd := exec.Command(sh_cmd, "-c", j.RawCommand)
	log.Printf("%09d : <%d|%d> %s", j.Sequence, j.Prio, j.RawPrio, j.RawCommand)
	if err := cmd.Run(); err != nil {
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
		}
	} else {
		// Command was successful
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	}
	j.Took = time.Now().Sub(j.Created)
	j.ExitStatus = waitStatus.ExitStatus()
	log.Printf("%09d : exit(%d), %.3fs", j.Sequence, j.ExitStatus, j.Took.Seconds())
	sem.Release()
	Stats()
}
