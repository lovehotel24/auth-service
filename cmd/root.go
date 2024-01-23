package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lovehotel24/auth-service/pkg/foundation/logger"
	"github.com/lovehotel24/auth-service/pkg/handlers"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "auth-service",
	Short: "authentication and authorization service for love hotel24",
	Run:   runCommand,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.auth-service.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCommand(cmd *cobra.Command, args []string) {
	// Construct the application logger.
	log, err := logger.New("AUTH")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := handlers.Run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}
