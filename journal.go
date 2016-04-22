////////////////////////////////////////////////////////////////////////////////
// Program for keeping track of interruptions.
// Interruption is logged from command line, with the name of the person who
// interrupted, the duration in minutes, and the reason for the interruption.
//
// Supported commands are:
// $ journal log <name> <duration> <reason>
//
// Exapmles:
// $ journal log ralph 10 "asked me about test reports"
// $ journal log sara 5 "needed help with a SQL query"
// $ journal log ralph 17 "asked me again about test reports. aargh, why won't he shut up??"
// $ journal log lynn 13 "reported a bug"
//
// $ journal list
// ralph 10 "asked me about test reports"
// sara   5 "needed help with SQL query"
// ralph 17 "asked me again about test reports. aargh, why won't he shut up??"
// lynn  13 "reported a bug"
//
// $ journal total
// 45
//
// $ journal hitlist
// ralph 27
// lynn 13
// sara 5
////////////////////////////////////////////////////////////////////////////////

package main
import (
	"github.com/codegangsta/cli"
	"os"
	"fmt"
	"strconv"
	"github.com/tornyak/quinn/db"
	"text/tabwriter"
)

// Limit input size
const (
	MaxDuration = 24 * 60
	MaxNameLength = 20
	MaxReasonLength = 120
)

// Argument positions for journal log command
const(
	ArgPositionLogName = iota
	ArgPositionLogDuration
	ArgPositionLogReason
)

// Journal implements methods for handling CLI commands
// In the background it communicates with the DB
type Journal struct {
	dbHandler *db.DBHandler
}

// Log single interruption
// CLI: journal log <person name> <duration minutes> <reason>
func (j *Journal)Log(c *cli.Context) {
	if len(c.Args()) != 3 {
		fmt.Printf("Usage: %v\n", c.Command.ArgsUsage)
		return
	}
	name := c.Args().Get(ArgPositionLogName)
	reason := c.Args().Get(ArgPositionLogReason)

	duration, err := strconv.ParseInt(c.Args().Get(ArgPositionLogDuration), 10, 64)
	if err != nil {
		fmt.Printf("Error reading duration: %v", c.Args().Get(1))
		os.Exit(1)
	}
	// assume that interruption cannot be longer than whole day :-)

	if duration <= 0 || duration > MaxDuration {
		fmt.Printf("Duration: %d must be in range 0 - %d", duration, MaxDuration)
		os.Exit(1)
	}
	// Name and description will be truncated if too long
	if len(name) > MaxNameLength {
		name = name[0:MaxNameLength]
	}
	if len(reason) > MaxReasonLength {
		reason = name[0:MaxReasonLength]
	}

	j.dbHandler.Log(name, duration, reason)
}

// List interruption entries, discards extra params if exist
// CLI: journal list
func (j *Journal)List(c *cli.Context) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 4, 8, 0, '\t', 0)
	for _, i := range j.dbHandler.List() {
		fmt.Fprintf(w, "%s\t%d\t\"%s\"\n", i.Name, i.Duration, i.Reason)
	}
	w.Flush()
}

// Total prints out total time in minutes spent under interrupts
// Extra parameters are discarded
// CLI: journal total
func (j *Journal)Total(c *cli.Context) {
	fmt.Println(j.dbHandler.Total())
}

// HitList summs up all interrupts grouped by person's name
// Result is sorted by total time in descending order
// Extra parameters are discarded
// CLI: journal hitlist
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
