package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(countryCmd)
}

var countryCmd = &cobra.Command{
	Use:   "country",
	Short: "Show top sites URL by country",
	Run: func(cmd *cobra.Command, args []string) {

		providerName := "SEMRush"

		p, ok := providers[providerName]

		if !ok {
			log.Printf("Ranking Provider (%s) not available.\n", providerName)
			return
		}

		urls, err := p.TopSitesCountry()

		if err != nil {
			log.Print(err)
			return
		}

		fmt.Printf("%v\n", urls)
	},
}
