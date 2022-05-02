package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(topCmd)
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Show top sites URL",
	Run: func(cmd *cobra.Command, args []string) {

		providerName := "SEMRush"

		p, ok := providers[providerName]

		if !ok {
			log.Printf("Ranking Provider (%s) not available.\n", providerName)
			return
		}

		urls, err := p.TopSitesGlobal()

		if err != nil {
			log.Print(err)
			return
		}

		fmt.Printf("%v\n", urls)
	},
}
