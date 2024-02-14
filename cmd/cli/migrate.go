package cli

import (
	"ffxvi-bard/config"
	database "ffxvi-bard/infrastructure/database/sql/migration"
	"github.com/spf13/cobra"
)

// migrateCMD fetches information about migrations
var migrateCMD = &cobra.Command{
	Use:   "migrate",
	Short: "m",
	Run: func(cmd *cobra.Command, args []string) {
		up, _ := cmd.Flags().GetBool("up")
		down, _ := cmd.Flags().GetBool("down")
		//version, _ := cmd.Flags().GetUint("version")
		config, err := config.NewConfig()
		if !up && !down {
			panic("Please provide a command type. Supported commands `up`, `down`")
		}
		if err != nil {
			panic(err)
		}
		driver := database.NewMigrationDriver(config.Database.Database, config.Database.Path)
		var command string
		if up {
			command = "up"
		}
		if down {
			command = "down"
		}
		//if version != 0 {
		//	driver.ExecuteOne(command, version)
		//	return
		//}
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
	//migrateCMD.PersistentFlags().Uint("version", 0, "Migration version")
}

func initRequiredMigrationFlags() {
}
