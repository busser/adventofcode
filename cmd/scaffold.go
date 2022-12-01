package cmd

import (
	"fmt"
	"time"

	"github.com/busser/adventofcode/scaffolding"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scaffoldCmd represents the scaffold command
var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Generate code to kickstart your solution",
	Long: `Generate code to kickstart your solution.

Examples:
  # Build scaffolding for day 1.
  adventofcode scaffold --day=1

  # Overwrite existing code.
  adventofcode scaffold --day=1 --force

  # Provide a session cookie to download your input of the day.
  adventofcode scaffold --day=1 --cookie=abcdef0123...

The CLI assumes you are working on the latest Advent of Code. You can always
override this behavior with the '--year' flag.

To download your input, provide the value of the 'session' cookie for the
adventofcode.com website. You can do this with the '--cookie' flag, the
ADVENTOFCODE_COOKIE environment variable, or by setting the 'cookie' field in
your configuration file.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		gen, err := scaffolding.NewGenerator(
			viper.GetInt("day"),
			viper.GetInt("year"),
			viper.GetString("workdir"),
			viper.GetString("cookie"),
			viper.GetBool("force"),
		)
		if err != nil {
			return fmt.Errorf("making code generator: %w", err)
		}

		if err := gen.Run(); err != nil {
			return fmt.Errorf("building scaffolding: %w", err)
		}

		fmt.Println("🎅🏻 Merry coding!")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)

	scaffoldCmd.Flags().IntP("day", "d", 0, "The day to build scaffolding for")

	// Year defaults to latest Advent of Code.
	year, month, _ := time.Now().Date()
	if month < time.December {
		year--
	}
	scaffoldCmd.Flags().IntP("year", "y", year, "The year of Advent of Code you are working on")

	scaffoldCmd.Flags().StringP("workdir", "w", "", "Your Advent of Code working directory")
	scaffoldCmd.Flags().StringP("cookie", "c", "", "Your session cookie for adventofcode.com")
	scaffoldCmd.Flags().BoolP("force", "f", false, "If true, overwrite existing files")

	viper.BindPFlags(scaffoldCmd.Flags())
}
