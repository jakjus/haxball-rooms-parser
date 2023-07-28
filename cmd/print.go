package cmd

import (
	"fmt"
        "os"
	"github.com/spf13/cobra"
        "github.com/jedib0t/go-pretty/v6/table"
)

func printCmdE() {
	body, _ := GetData()
	serverList := Parse(body)
        t := table.NewWriter()
        t.SetOutputMirror(os.Stdout)
        t.AppendHeader(table.Row{"Room Name", "Players", "Flag", "Pass", "Link"})
        for _, s := range serverList {
          t.AppendRow(table.Row{s.Name, fmt.Sprintf("%v/%v", s.PlayersNow, s.PlayersMax), fmt.Sprintf("%s", s.Flag), s.Private, fmt.Sprintf("%s", s.Link)})
	}
        t.Render()
	//fmt.Printf("\nTotal Servers: %v\n", len(serverList))
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print Servers information",
	Run: func(cmd *cobra.Command, args []string) {
		printCmdE()
	},
}
