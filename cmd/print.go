package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func printCmdE() {
	body, _ := GetData()
	serverList := Parse(body)
	for _, s := range serverList {
		s.Print()
		fmt.Print("\n")
	}
	fmt.Printf("\nTotal Servers: %v\n", len(serverList))
}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print Servers information",
	Run: func(cmd *cobra.Command, args []string) {
		printCmdE()
	},
}
