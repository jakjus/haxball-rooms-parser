package cmd

import (
	"github.com/spf13/cobra"
)

var echoTimes int

var rootCmd = &cobra.Command{Use: "hbparser"}

func init() {
	//cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")
	rootCmd.AddCommand(tableCmd)
	rootCmd.AddCommand(dbCmd)
	//cmdEcho.AddCommand(cmdTimes)
}

func Execute() {
	rootCmd.Execute()
}
