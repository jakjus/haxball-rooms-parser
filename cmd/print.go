package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func printCmdE() {
	body, _ := GetData()
	serverList := Parse(body)
	for i, s := range serverList {
		fmt.Println(i+1)
		s.Print()
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
