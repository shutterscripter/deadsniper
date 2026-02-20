/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"fmt"
	"deadsniper/scrapper"
	"deadsniper/config"
	"github.com/spf13/cobra"
	"encoding/json"
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deadsniper",
	Short: "Dead link finder for URLs",
	Long:  `Finds broken or dead links on a given URL by scraping the page and checking each link.`,
	Run: func(cmd *cobra.Command, args []string) {
		deadLinks, blockedByBot, err := scrapper.VisitUrl(config.DefaultConfig.URL)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), err)
			return
		}
		if len(deadLinks) > 0 {
			if config.DefaultConfig.OutputType != 0  {

				switch config.DefaultConfig.OutputType {
				case 1:
					file, err := os.Create("data.txt")
					if err != nil {
						fmt.Println("Error creating file:", err)
						return
					}
					for _, link := range deadLinks {
						
						_, err = file.WriteString(link + "\n")
						if err != nil {
							fmt.Println("Error writing to file:", err)
							return
						}
					}
					defer file.Close()
					fmt.Println("Result written to file: data.txt")
					
				case 2:
					file, err := os.Create("data.json")
					if err != nil {
						fmt.Println("Error creating file:", err)
						return
					}
					for _, link := range deadLinks {
						
						jsonData, err := json.Marshal(link)
						if err != nil {
							fmt.Println("Error marshalling data:", err)
							return
						}
						_, err = file.WriteString(string(jsonData) + "\n")
						if err != nil {
							fmt.Println("Error writing to file:", err)
							return
						}	
					}
					defer file.Close()
					fmt.Println("Result written to file: data.json")
				case 3:
					fmt.Println("CSV output")
				case 4:
					fmt.Println("XML output")
				
				}
			}
		}
		if len(blockedByBot) > 0 {
			// fmt.Printf("\nLinks blocked by server / bot not allowed (403) (%d):\n", len(blockedByBot))
			// for _, link := range blockedByBot {
			// 	fmt.Println(link)
			// }
		}
		if len(deadLinks) == 0 && len(blockedByBot) == 0 {
			fmt.Println("No dead links and no blocked links found.")
		}
		
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.DefaultConfig.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().IntVarP(&config.DefaultConfig.Threads, "threads", "t", 1, "Number of threads to use")
	rootCmd.PersistentFlags().Float64VarP(&config.DefaultConfig.Delay, "delay", "d", 0.5, "Delay between requests")
	rootCmd.PersistentFlags().IntVarP(&config.DefaultConfig.Timeout, "timeout", "T", 10, "Timeout in seconds")
	rootCmd.PersistentFlags().IntVarP(&config.DefaultConfig.OutputType, "output-type", "o", 1, "Output type: 1: text, 2: json, : csv, 4: xml")
	rootCmd.PersistentFlags().BoolVarP(&config.DefaultConfig.Help, "help", "h", false, "Help for the command")
	rootCmd.Flags().StringVarP(&config.DefaultConfig.URL, "url", "u", "", "URL to scrape")

}


