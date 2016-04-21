package main
import (
	"github.com/codegangsta/cli"
	"os"
	"fmt"
	"strconv"
	"github.com/tornyak/quinn/db"
	"text/tabwriter"
)

type Journal struct {
	dbHandler *db.DBHandler
}

// Log interruption
// journal log <person name> <duration minutes> <reason>
func (j *Journal)Log(c *cli.Context) {
	if len(c.Args()) != 3 {
		fmt.Printf("Usage: %v\n", c.Command.ArgsUsage)
		return
	}
	name := c.Args().Get(0)
	reason := c.Args().Get(2)

	duration, err := strconv.ParseInt(c.Args().Get(1), 10, 64)
	if err != nil {
		fmt.Printf("Error reading duration: %v", c.Args().Get(1))
	}
	if duration <= 0 {
		fmt.Printf("Duration: %d must be greater then 0", duration)
	}
	j.dbHandler.Log(name, duration, reason)
}

// List interruption entries, discards extra params if exist
// journal list
func (j *Journal)List(c *cli.Context) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 8, 2, '\t', 0)
	for _, i := range j.dbHandler.List() {
		fmt.Fprintf(w, "%s\t%d\t\"%s\"\n", i.Name, i.Duration, i.Reason)
	}
	w.Flush()
}

// Total prints out total time in minutes spent under interrupts
func (j *Journal)Total(c *cli.Context) {
	fmt.Println(j.dbHandler.Total())
}

func (j *Journal)Hitlist(c *cli.Context) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 8, 2, '\t', 0)
	for _, hit := range j.dbHandler.Hitlist() {
		fmt.Fprintf(w, "%s\t%d\n", hit.Name, hit.Duration)
	}
	w.Flush()
}


func main() {
	dbHandler := db.NewDBHandler()
	journal := &Journal{dbHandler}

	app := cli.NewApp()
	app.Name = "journal"
	app.Usage = "Quinn's journal of interrupts"
	app.UsageText = ""
	app.ArgsUsage = ""
	app.HelpName = ""
	app.HideVersion = true
	app.HideHelp = true

	app.Commands = []cli.Command{
		{
			Name: "log",
			Usage:  "Log interruption",
			ArgsUsage: "journal log <person name> <duration minutes> <reason>",
			Action: journal.Log,
		},
		{
			Name: "list",
			Usage:  "List logged entries",
			Action: journal.List,
		},
		{
			Name: "total",
			Usage:  "Show total interruption time in minutes",
			Action: journal.Total,
		},
		{
			Name: "hitlist",
			Usage:  "Show interruptions per person sorted by duration",
			Action: journal.Hitlist,
		},
	}

	app.Run(os.Args)
}
