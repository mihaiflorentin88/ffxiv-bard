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
		count, _ := cmd.Flags().GetInt("count")
		if !generate && !execute {
			panic("Please provide a command type. Supported commands `generate`, `execute`")
		}
		if generate {
			fixtures.GenerateFixtures(count)
		}
		if execute {
			databaseDriver := container.Load.DatabaseDriver()
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
	FixturesCMD.PersistentFlags().Int("count", 100, "The number of fixtures to generate")
}

func initRequiredFixturesFlags() {
}
