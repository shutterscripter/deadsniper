/*
Copyright © 2026 Jayesh Shinde jay2272001@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"deadsniper/config"
	"deadsniper/scrapper"
)

// visitCmd represents the visit command
var visitCmd = &cobra.Command{
	Use:   "check",
	Short: "Check a URL for dead links",
	Long:  `Scrapes the given URL and reports which links are broken or unreachable.`,
	Run: func(cmd *cobra.Command, args []string) {
		deadLinks, blockedByBot, err := scrapper.VisitUrl(config.DefaultConfig.URL)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), err)
			return
		}
		if len(deadLinks) > 0 {
			fmt.Printf("Dead links (%d):\n", len(deadLinks))
			for _, link := range deadLinks {
				fmt.Println(link)
			}
		}
		if len(blockedByBot) > 0 {
			fmt.Printf("\nLinks blocked by server / bot not allowed (403) (%d):\n", len(blockedByBot))
			for _, link := range blockedByBot {
				fmt.Println(link)
			}
		}
		if len(deadLinks) == 0 && len(blockedByBot) == 0 {
			fmt.Println("No dead links and no blocked links found.")
		}
	},
}

func init() {
	rootCmd.AddCommand(visitCmd)
	visitCmd.Flags().StringVarP(&config.DefaultConfig.URL, "url", "u", "", "URL to scrape")
}
