package gron

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"
)

type JobsSorter []*Job

type Job struct {
	Sequence   int
	Created    time.Time
	RawCommand string
	RawPrio    int
	Prio       int
	ExitStatus int
	Took       time.Duration
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
	if js[i].Prio == js[j].Prio {
		return js[i].Created.Nanosecond() > js[j].Created.Nanosecond()
	}

	return js[i].Prio < js[j].Prio
}
