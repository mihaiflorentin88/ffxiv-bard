package cli

import (
	"ffxvi-bard/cmd/http"
	"github.com/spf13/cobra"
)

// httpCMD fetches information about migrations
var httpCmd = &cobra.Command{
	Use:   "server",
	Short: "s",
	Run: func(cmd *cobra.Command, args []string) {
		start, _ := cmd.Flags().GetBool("start")
		port, _ := cmd.Flags().GetInt("port")
		poolSize, _ := cmd.Flags().GetInt("pool")
		if start {
			http.Server(port, poolSize)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	initHttpFlags()
	initRequiredHttplags()
}

func initHttpFlags() {
	httpCmd.PersistentFlags().Bool("start", false, "Starts the web server")
	httpCmd.PersistentFlags().Int("port", 80, "Sets the server port")
	httpCmd.PersistentFlags().Int("pool", 10, "Sets the server pool size")
}

func initRequiredHttplags() {
}
