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

// ec2Cmd represents the ec2 command
var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ec2 called")

		instanceName, _ := cmd.Flags().GetString("instanceName") // ./terrathor provider --region=eu-west-1
		instanceType, _ := cmd.Flags().GetString("instanceType")
		//availabilityZone, _ := cmd.Flags().GetString("availabilityZone")
		test, _ := cmd.Flags().GetString("test")

		if instanceName != "" && instanceType != "" {
			CreateEC2Module("modules/ec2", instanceName, test)
			//CreateEC2("test", instanceName, availabilityZone, instanceType)
		}

	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)

	ec2Cmd.PersistentFlags().String("instanceName", "", "EC2 resource name")
	ec2Cmd.PersistentFlags().String("instanceType", "", "Instance type")
	ec2Cmd.PersistentFlags().String("availabilityZone", "", "Availability zone")

	ec2Cmd.MarkPersistentFlagRequired("instanceName")
	ec2Cmd.MarkPersistentFlagRequired("instanceType")
	ec2Cmd.MarkPersistentFlagRequired("availabilityZone")

}

func CreateEC2Module(folderName string, instanceName string, test string) {

	vars := make(map[string]string)
	vars["instance_name"] = instanceName
	vars["test"] = test

	pies := []File{
		{
			TemplatePath: "./templates/ec2/main.tmpl",
			FilePath:     "modules/ec2/main.tf",
			FolderPath:   folderName,
		},
		{
			TemplatePath: "./templates/ec2/outputs.tmpl",
			FilePath:     "modules/ec2/outputs.tf",
			FolderPath:   folderName,
		},
		{
			TemplatePath: "./templates/ec2/data.tmpl",
			FilePath:     "modules/ec2/data.tf",
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

func CreateEC2(folderName string, instanceName string, availabilityZone string, instanceType string) {

	vars := make(map[string]string)
	vars["instance_name"] = instanceName
	vars["availability_zone"] = availabilityZone
	vars["instance_type"] = instanceType

	p := &File{
		TemplatePath: "./templates/ec2/ec2.tmpl",
		FilePath:     "test/ec2.tf",
		FolderPath:   folderName,
	}

	p.CreateFolder()

	t := p.Parse()

	fs := p.CreateFile()

	err := t.Execute(fs, vars)

	p.ErrorHandler(err)

	defer fs.Close()

}
