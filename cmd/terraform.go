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

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Terraform creating...")

		region, _ := cmd.Flags().GetString("region") // ./terrathor provider --region=eu-west-1
		bucket, _ := cmd.Flags().GetString("bucket")
		key, _ := cmd.Flags().GetString("key")

		if region != "" && bucket != "" && key != "" {
			createTerraform("test", region, bucket, key)
		}
	},
}

func init() {
	rootCmd.AddCommand(terraformCmd)

	terraformCmd.PersistentFlags().String("region", "", "Region of S3")
	terraformCmd.PersistentFlags().String("bucket", "", "Bucket name")
	terraformCmd.PersistentFlags().String("key", "", "State file name")

	terraformCmd.MarkPersistentFlagRequired("region")
	terraformCmd.MarkPersistentFlagRequired("bucket")
	terraformCmd.MarkPersistentFlagRequired("key")

}

func createTerraform(folderName string, region string, bucket string, key string) {

	vars := make(map[string]string)
	vars["region"] = region
	vars["bucket_name"] = bucket
	vars["key"] = key

	p := &File{
		TemplatePath: "./templates/terraform/terraform.tmpl",
		FilePath:     "test/terraform.tf",
		FolderPath:   folderName,
	}

	p.CreateFolder()

	t := p.Parse()

	fs := p.CreateFile()

	err := t.Execute(fs, vars)

	p.ErrorHandler(err)

	defer fs.Close()

}
