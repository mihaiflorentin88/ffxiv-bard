package cli

import (
	"ffxvi-bard/container"
	"ffxvi-bard/infrastructure/database/sql/fixtures"
	"github.com/spf13/cobra"
)

var FixturesCMD = &cobra.Command{
	Use:   "fixtures",
	Short: "f",
	Run: func(cmd *cobra.Command, args []string) {
		generate, _ := cmd.Flags().GetBool("generate")
		execute, _ := cmd.Flags().GetBool("execute")
		if !generate && !execute {
			panic("Please provide a command type. Supported commands `generate`, `execute`")
		}
		if generate {
			fixtures.GenerateFixtures()
		}
		if execute {
			databaseDriver, err := container.GetDatabaseDriver()
			if err != nil {
				panic("Cannot connect to database")
			}
			fixtures := fixtures.NewFixtures(databaseDriver)
			fixtures.Execute()
		}
	},
}

func init() {
	rootCmd.AddCommand(FixturesCMD)
	initFixturesFlags()
	initRequiredFixturesFlags()
}

func initFixturesFlags() {
	FixturesCMD.PersistentFlags().Bool("generate", false, "Generate fixtures")
	FixturesCMD.PersistentFlags().Bool("execute", false, "Execute fixtures")
}

func initRequiredFixturesFlags() {
}
