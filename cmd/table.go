package cmd

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

func tableCmdE() {
	body, _ := GetData()
	serverList := Parse(body)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Room Name", "Players", "Flag", "Pass", "Link"})
	for _, s := range serverList {
		t.AppendRow(table.Row{s.Name, fmt.Sprintf("%v/%v", s.PlayersNow, s.PlayersMax), fmt.Sprintf("%s", s.Flag), s.Private, fmt.Sprintf("%s", s.Link)})
	}
	t.SetStyle(table.StyleColoredBright)
	t.SortBy([]table.SortBy{{Name: "Room Name", Mode: table.Asc}})
	t.AppendFooter(table.Row{"Total Servers", len(serverList), "", "", ""})

	t.Render()
}

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Print out HaxBall rooms as table",
	Run: func(cmd *cobra.Command, args []string) {
		tableCmdE()
	},
}
