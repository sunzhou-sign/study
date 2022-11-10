package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your application",
	Long:  "A longer description",
}

var mockMsgCmd = &cobra.Command{
	Use:   "mockMsg",
	Short: "批量发送测试文本消息",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mockMsg called")
	},
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "导出数据",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("export called")
	},
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(mockMsgCmd)
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("out", "k", "./backup", "导出路径")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
