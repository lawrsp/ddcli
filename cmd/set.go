/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/lawrsp/ddcli/monitor"
	"github.com/spf13/cobra"
)

var setBrightness uint
var setContrast uint

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := monitor.GetAllPhysicalMonitors()
		if err != nil {
			fmt.Println("error occured:")
			fmt.Println(err)
			return
		}

		for _, handle := range list {
			if setBrightness <= 100 && setBrightness > 0 {
				err := monitor.SetMonitorBrightness(handle, setBrightness)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if setContrast <= 100 && setContrast > 0 {
				err := monitor.SetMonitorContrast(handle, setContrast)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		fmt.Println("OK")
	},
}

func init() {
	setCmd.Flags().UintVarP(&setBrightness, "brightness", "b", 0, "set the brightness")
	setCmd.Flags().UintVarP(&setContrast, "contrast", "c", 0, "set the contrast")

	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
