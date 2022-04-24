package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/bartmika/bleve-server/pkg/rpc_client" // Import so we can detect bugs on compile-time.
)

// Application environment variables
var (
	applicationAddress string
)

func init() {
	// Load up our `environment variables` from our operating system.
	appAddress := os.Getenv("BLEVE_SERVER_ADDRESS")

	// Attach environment variables to system.
	rootCmd.PersistentFlags().StringVar(&applicationAddress, "appAddress", appAddress, "The applications address.")

	viper.BindPFlag("appAddress", rootCmd.PersistentFlags().Lookup("appAddress"))

	viper.SetDefault("appAddress", appAddress)
}

var rootCmd = &cobra.Command{
	Use:   "bleve-server",
	Short: "RPC server over a single running bleve instance",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

// Execute is the main entry into the application from the command line terminal.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
