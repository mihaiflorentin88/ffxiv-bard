package cli

import (
	"ffxvi-bard/container"
	"github.com/spf13/cobra"
)

// migrateCMD fetches information about migrations
var migrateCMD = &cobra.Command{
	Use:   "migrate",
	Short: "m",
	Run: func(cmd *cobra.Command, args []string) {
		up, _ := cmd.Flags().GetBool("up")
		down, _ := cmd.Flags().GetBool("down")
		if !up && !down {
			panic("Please provide a command type. Supported commands `up`, `down`")
		}
		driver := container.GetMigrationDriver()
		var command string
		if up {
			command = "up"
		}
		if down {
			command = "down"
		}
		driver.Execute(command)
	},
}

func init() {
	rootCmd.AddCommand(migrateCMD)
	initMigrationlags()
	initRequiredMigrationFlags()
}

func initMigrationlags() {
	migrateCMD.PersistentFlags().Bool("up", false, "Up migration")
	migrateCMD.PersistentFlags().Bool("down", false, "Down migration")
}

func initRequiredMigrationFlags() {
}
