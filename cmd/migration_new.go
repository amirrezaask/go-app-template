/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"app/monitoring"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Need a name for migration")
			os.Exit(1)
		}
		now := time.Now()
		fileName := "%.5d_%s.%s.sql"
		ufd, err := os.Create("./db/migrations/" + fmt.Sprintf(fileName, now.Unix(), args[0], "up"))
		if err != nil {
			monitoring.Logger().Fatalln(err)
		}
		err = ufd.Close()
		if err != nil {
			monitoring.Logger().Fatalln(err)
		}
		dfd, err := os.Create("./db/migrations/" + fmt.Sprintf(fileName, now.Unix(), args[0], "down"))
		if err != nil {
			monitoring.Logger().Fatalln(err)
		}
		err = dfd.Close()
		if err != nil {
			monitoring.Logger().Fatalln(err)
		}
	},
}

func init() {
	migrationCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
