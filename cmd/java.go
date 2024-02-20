/*
Copyright © 2024 FF911
*/
package cmd

import (
	"chiyou.code/mmc/app"
	"github.com/spf13/cobra"
)

// javaCmd represents the java command
var javaCmd = &cobra.Command{
	Use:   "java",
	Short: "mmc java",
	Long:  `mmc java command`,
	Run: func(cmd *cobra.Command, args []string) {
		app.RunJava()
	},
}

func init() {
	rootCmd.AddCommand(javaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// javaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// javaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
