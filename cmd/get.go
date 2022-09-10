/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/lawrsp/ddcli/monitor"
	"github.com/spf13/cobra"
)

type GetResult struct {
	DeviceName string
	Brightness *monitor.MonitorBrightness
	Contrast   *monitor.MonitorContrast
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
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

		var result []GetResult
		for name, handle := range list {

			brightness, err := monitor.GetMonitorBrightness(handle)
			if err != nil {
				fmt.Println(err)
				return
			}
			contrast, err := monitor.GetMonitorContrast(handle)
			if err != nil {
				fmt.Println(err)
				return
			}

			result = append(result, GetResult{
				DeviceName: name,
				Brightness: brightness,
				Contrast:   contrast,
			})
		}

		bs, err := json.Marshal(result)
		if err != nil {
			fmt.Println("output result failed:", err)
			return
		}
		fmt.Println(string(bs))

	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
