package gron

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"sort"
	"syscall"
	"time"
)

const sock = "/tmp/gron.sock"
const sh_cmd = "/bin/bash"

var jobs []*Job
var fjobs []*Job
var sem *Semaphore
var sequence int
var waiting *Semaphore

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
	waiting = NewSemaphore(1)
	go engine()

	kv := <-ksignal
	log.Printf("Signal to finish : %s", kv.String())
	os.Remove(sock)
}

func appendJob(j *Job) bool {
	for _, ja := range jobs {
		if j.RawCommand == ja.RawCommand {
			ja.Prio++
			return false
		}
	}
	sequence++
	//reset sequence
	if sequence > 999999 {
		sequence = 1
	}
	j.Sequence = sequence
	j.Created = time.Now()
	j.Prio = j.RawPrio
	jobs = append(jobs, j)
	return true
}

func trait_request(c net.Conn) {
	defer c.Close()
	cr := NewClientRequest()
	cr.Decode(c)
	switch cr.Request {
	case "job":
		j := cr.Object.(Job)
		if appendJob(&j) {
			waiting.Release()
		}
		break
	case "status":
		s := NewStatus()
		s.MaxProcess = sem.Permits()
		s.Process = sem.Available() + len(jobs)
		s.Running = sem.Available()
		s.Sequence = sequence
		s.Waiting = jobs
		s.Finished = fjobs
		c.Write(s.Encode())
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
		} else {
			waiting.Acquire()
		}
	}
}

func finalJob(j *Job) {
	if len(fjobs) > 20 {
		fjobs = fjobs[:len(fjobs)-1]
	}
	fjobs = append([]*Job{j}, fjobs...)
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
	finalJob(j)
	sem.Release()
	Stats()
}
