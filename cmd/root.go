package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tex-screenshot",
	Short: "tex-screenshot is a cli tool for converting screenshots to LaTex",
	Long:  "tex-screenshot is a cli tool for converting screenshots to LaTex",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

type Response struct {
	RequestID string `json:"request_id"`
	Res       Result `json:"res"`
	Status    bool   `json:"status"`
}

type Result struct {
	Conf  float64 `json:"conf"`
	Latex string  `json:"latex"`
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing tex-screenshot '%s'\n", err)
		os.Exit(1)
	}
}
