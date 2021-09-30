/*
Copyright © 2021 Faruk AK <kakuraf@gmail.com>

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
	"fmt"

	"github.com/spf13/cobra"
)

// providerCmd represents the provider command
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Add a provider your module",
	Long:  `A module is a container for multiple resources that are used together.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Provider creating...")

		region, _ := cmd.Flags().GetString("region") // ./terrathor provider --region=eu-west-1

		if region != "" {
			createProvider("test", region)
		}
	},
}

func init() {
	rootCmd.AddCommand(providerCmd)

	providerCmd.PersistentFlags().String("region", "", "Region of provider")

	providerCmd.MarkPersistentFlagRequired("region")
}

func createProvider(folderName string, region string) {

	vars := make(map[string]string)
	vars["region"] = region

	pies := []File{
		{
			TemplatePath: "./templates/provider/provider.tmpl",
			FilePath:     "test/provider.tf",
			FolderPath:   folderName,
		},
		{
			TemplatePath: "./templates/provider/tfvars.tmpl",
			FilePath:     "test/terraform.tfvars",
			FolderPath:   folderName,
		},
		{
			TemplatePath: "./templates/provider/variables.tmpl",
			FilePath:     "test/variables.tf",
			FolderPath:   folderName,
		},
	}

	for _, p := range pies {

		p.CreateFolder()

		t := p.Parse()

		fs := p.CreateFile()

		err := t.Execute(fs, vars)

		p.ErrorHandler(err)

		defer fs.Close()

	}

}
