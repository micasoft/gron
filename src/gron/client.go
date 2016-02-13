package gron

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"log"
	"net"
	"os"
	"time"
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
	bcr.Request = "job"
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
		fmt.Printf("\n\n%s : \t%s\n", cyan("Maximum of process"), yellow(s.MaxProcess))
		fmt.Printf("%s : \t%s\n", cyan("Total of process"), yellow(s.Process))
		fmt.Printf("%s : \t%s\n", cyan("Total of running"), yellow(s.Running))
		fmt.Printf("%s: \t\t%s\n\n", cyan("Total of seq."), yellow(s.Sequence))

		if len(s.Waiting.([]*Job)) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Sequence", "Landing", "Command", "Priority"})
			for _, j := range s.Waiting.([]*Job) {
				v := []string{fmt.Sprintf("%09d", j.Sequence), fmt.Sprintf("%.3f", time.Now().Sub(j.Created).Seconds()), j.RawCommand, fmt.Sprintf("%d", j.Prio)}
				table.Append(v)
			}
			table.Render()
		}

		if len(s.Finished.([]*Job)) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Sequence", "Took", "Command", "Priority", "Exit"})
			for _, j := range s.Finished.([]*Job) {
				v := []string{fmt.Sprintf("%09d", j.Sequence), fmt.Sprintf("%.3f", j.Took.Seconds()), j.RawCommand, fmt.Sprintf("%d", j.Prio), fmt.Sprintf("%d", j.ExitStatus)}
				table.Append(v)
			}
			table.Render()
		}
	}
}
