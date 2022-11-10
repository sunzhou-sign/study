/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var output string

// wgetCmd represents the wget command
var wgetCmd = &cobra.Command{
	Use:   "wget",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Example: "xpower wget iqsing.github.io/download.tar -o /tmp",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wget called")
		fmt.Println(args)
		Download(args[0], output)
	},
}

func init() {
	rootCmd.AddCommand(wgetCmd)
	wgetCmd.Flags().StringVarP(&output, "output", "o", "", "output file")
	wgetCmd.MarkFlagRequired("output")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wgetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wgetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Download(url string, path string) {
	out, err := os.Create(path)
	check(err)
	defer out.Close()

	res, err := http.Get(url)
	check(err)
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	check(err)
	fmt.Printf("save as " + path)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
