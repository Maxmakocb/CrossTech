package cmd

import (
	"fmt"
	"os"
   
	"github.com/spf13/cobra"
)

var (
	srvPort = 0
)

// An assumption is made that only the developing team has the access to these scripts.
var rootCmd = &cobra.Command{
	Use:  "cross_tech",
    Short: "cross_tech - cli to interact with the application",

	Run: func(cmd *cobra.Command, args []string) {
		serve()
    },
}

var PopulateDB = &cobra.Command{
	Use: "populate",
	Run: func(cmd *cobra.Command, args []string) {
		populate()
	},
}

func init() {
	rootCmd.AddCommand(PopulateDB)
	rootCmd.Flags().IntVar(&srvPort, "server-port", 8080, "port for the server to listen on")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed during execution: '%s'", err)
        os.Exit(1)
    }
}