/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (

	"github.com/spf13/cobra"
	"deadsniper/scrapper"
)

// visitCmd represents the visit command
var visitCmd = &cobra.Command{
	Use:   "check",
	Short: "Check a URL for dead links",
	Long:  `Scrapes the given URL and reports which links are broken or unreachable.`,
	Run: func(cmd *cobra.Command, args []string) {

		scrapper.VisitUrl(scrapper.DefaultConfig.URL)
	},
}

func init() {
	rootCmd.AddCommand(visitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// visitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	visitCmd.Flags().StringVarP(&scrapper.DefaultConfig.URL, "url", "u", "", "URL to scrape")
}
