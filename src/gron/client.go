package gron

import (
	"log"
	"net"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"fmt"
	"os"
	"strconv"
)

func connect() net.Conn {
	c, err := net.Dial("unix", sock)
	if err != nil {
		panic(err)
	}
	return c
}

func Client(cmd *string, prio *int) {
	c := connect()
	defer c.Close()
	bcr := NewClientRequest()
	bcr.Request= "job"
	bcr.Object = Job{RawCommand: *cmd, RawPrio: *prio}
	_, err := c.Write(bcr.Encode())
	if err != nil {
		log.Fatal("write error:", err)
	}
}

func GetStatus() {
	c := connect()
	defer c.Close()
	bcr := ClientRequest{Request: "status"}
	_, err := c.Write(bcr.Encode())
	if err != nil {
		log.Fatal("write error:", err)
	} else {
		s := NewStatus()
		s.Decode(c)
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Printf("%s : \t%s\n", cyan("Total of process"), yellow(s.Process))
		fmt.Printf("%s : \t%s\n", cyan("Total of running"), yellow(s.Running))
		fmt.Printf("%s: \t\t%s\n\n", cyan("Total of seq."), yellow(s.Sequence))

		if (len(s.Waiting.([]*Job)) > 0) {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Landing", "Command", "Priority"})
			for _, j := range s.Waiting.([]*Job) {
				v := []string{j.Created.Format("15:04:05.000"), j.RawCommand,  strconv.Itoa(j.Prio)}
	    		table.Append(v)
			}
			table.Render()
		}
	}
}

