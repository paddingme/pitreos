package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var timestampString string

var restoreCmd = &cobra.Command{
	Use:   "restore {SOURCE} {DESTINATION}",
	Short: "Restores your files to a specified point in time (default: latest available)",
	Args:  cobra.ExactArgs(2),
	Long: `Restores your files to the closest available backup before 
the requested timestamp (default: now). 
It compares existing chunks of data in your files and downloads only the necessary data. 
This is optimized for large and sparse files, like virtual machines disks or nodeos state.`,
	Run: func(cmd *cobra.Command, args []string) {

		pitr := getPITR()
		t, err := parseUnixTimestamp(timestampString)
		if err != nil {
			fmt.Printf("Got error: %s\n", err)
			os.Exit(1)
		}
		err = pitr.RestoreFromBackup(args[0], args[1], t)
		if err != nil {
			fmt.Printf("Got error: %s\n", err)
			os.Exit(1)
		}
	},
	Example: `  pitreos restore gs://mybackups/projectname file:///home/nodeos/data -c --timestamp $(date -d "2 hours ago" +%s)`,
}

func parseUnixTimestamp(unixTimeStamp string) (tm time.Time, err error) {
	i, err := strconv.ParseInt(unixTimeStamp, 10, 64)
	if err != nil {
		return
	}
	tm = time.Unix(i, 0)
	return
}

// adding the "Args" definition (SOURCE / DESTINATION) right below the USAGE definition
var restoreUsageTemplate = `Usage:{{if .Runnable}}
  {{if .HasAvailableFlags}}{{appendIfNotPresent .UseLine "[flags]"}}{{else}}{{.UseLine}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
  {{ .CommandPath}} [command]{{end}}
  * SOURCE: File path (ex: /var/backups) or Google Storage URL (ex: gs://mybackups/projectname)
  * DESTINATION: File path (ex: ../mydata)
  {{if gt .Aliases 0}}
Aliases:
  {{.NameAndAliases}}
{{end}}{{if .HasExample}}
Examples: 
{{ .Example }}{{end}}{{ if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableLocalFlags}}
Flags:
{{.LocalFlags.FlagUsages | trimRightSpace}}{{end}}{{ if .HasAvailableInheritedFlags}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimRightSpace}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsHelpCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{ if .HasAvailableSubCommands }}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func init() {
	restoreCmd.Flags().StringVarP(&timestampString, "timestamp", "t", "", "Timestamp before which we want the latest available backup")
	restoreCmd.SetUsageTemplate(restoreUsageTemplate)
	if timestampString == "" {
		timestampString = strconv.FormatInt(time.Now().Unix(), 10)
	}
}
