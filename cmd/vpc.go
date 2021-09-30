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
	"github.com/spf13/viper"
)

// vpcCmd represents the vpc command
var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Provides a VPC resource.",
	Long:  `Provides a VPC resource.`,
	// PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
	// 	// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
	// 	return initializeConfig(cmd)
	// },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Vpc creating...")

		// ./terrathor provider --region=eu-west-1
		cidr, _ := cmd.Flags().GetString("cidr")
		availabilityZone, _ := cmd.Flags().GetString("availabilityZone")

		if cidr != "" && availabilityZone != "" {
			CreateVPCModule("modules/vpc", cidr)
			CreateVPC("test", cidr, availabilityZone)
		}

	},
}

func init() {
	rootCmd.AddCommand(vpcCmd)

	vpcCmd.PersistentFlags().String("cidr", "", "VPC cidr block")
	vpcCmd.PersistentFlags().String("availabilityZone", "", "Availability zone")

	vpcCmd.MarkPersistentFlagRequired("cidr")
	vpcCmd.MarkPersistentFlagRequired("availabilityZone")

	//biind the flags to the configuration
	viper.BindPFlag("cidr", vpcCmd.Flags().Lookup("cidr"))
	viper.BindPFlag("availabilityZone", vpcCmd.Flags().Lookup("availabilityZone"))
	viper.Set("database.user", "newuser")
	viper.Set("owner.name", "John")
	viper.WriteConfig()

}

func CreateVPCModule(folderName string, cidr string) {
	a := viper.GetString("url")
	fmt.Println(a)

	vars := make(map[string]string)

	pies := []File{
		{
			TemplatePath: "./templates/vpc/main.tmpl",
			FilePath:     "modules/vpc/main.tf",
			FolderPath:   folderName,
			configName:   "conf",
		},
		{
			TemplatePath: "./templates/vpc/outputs.tmpl",
			FilePath:     "modules/vpc/outputs.tf",
			FolderPath:   folderName,
			configName:   "conf",
		},
		{
			TemplatePath: "./templates/vpc/variables.tmpl",
			FilePath:     "modules/vpc/variables.tf",
			FolderPath:   folderName,
			configName:   "conf",
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

func CreateVPC(folderName string, cidr string, availabilityZone string) {

	vars := make(map[string]string)
	vars["cidr_block"] = cidr
	vars["availability_zone"] = availabilityZone

	p := &File{
		TemplatePath: "./templates/vpc/vpc.tmpl",
		FilePath:     "test/vpc.tf",
		FolderPath:   folderName,
		configName:   "conf",
	}

	p.InitConfig()

	p.CreateFolder()

	t := p.Parse()

	fs := p.CreateFile()

	err := t.Execute(fs, vars)

	p.ErrorHandler(err)

	defer fs.Close()

}
