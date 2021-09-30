package cmd

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

type Filer interface {
	CreateFolder()
	CheckFolder(name string)
	Parse() *template.Template
	CreateFile() *os.File
	ErrorHandler(err error)
	InitTerraform()
	ValidateTerraform()
	CreateSSHKeygen()
}

type File struct {
	TemplatePath string
	FilePath     string
	FolderPath   string
	configName   string
}

func (p *File) CreateFolder() {

	err := os.MkdirAll(p.FolderPath, 0755)

	if err != nil {
		p.ErrorHandler(err)
	}

	p.CheckFolder(p.FolderPath)
}

func (p *File) CheckFolder(name string) {
	_, err := os.Stat(name)

	if os.IsNotExist(err) {
		p.ErrorHandler(err)
	}

}

func (p *File) Parse() *template.Template {
	t, err := template.ParseFiles(p.TemplatePath)

	if err != nil {
		p.ErrorHandler(err)
	}

	return t
}

func (p *File) CreateFile() *os.File {

	is, err := p.CheckFile(p.FilePath)

	p.ErrorHandler(err)

	if !is {
		f, err := os.Create(p.FilePath)

		if err != nil {
			p.ErrorHandler(err)
		}

		return f
	} else {

		f, err := os.OpenFile(p.FilePath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			p.ErrorHandler(err)
		}

		return f
	}

}

func (p *File) CheckFile(fileName string) (bool, error) {
	_, err := os.Stat(fileName)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (p *File) ErrorHandler(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (p *File) InitTerraform() {
	out, err := exec.Command("powershell", "terraform init").Output()

	if err != nil {
		fmt.Println(err.Error())
		//log.Fatal(err)
	}

	fmt.Println(string(out))
}

func (p *File) ValidateTerraform() {
	out, err := exec.Command("powershell", "terraform validate").Output()

	if err != nil {
		fmt.Println(err.Error())
		//log.Fatal(err)
	}

	fmt.Println(string(out))
}

func (p *File) CreateSSHKeygen() {
	out, err := exec.Command("powershell", "ssg-keygen").Output()

	if err != nil {
		fmt.Println(err.Error())
		//log.Fatal(err)
	}

	fmt.Println(string(out))
}

func (p *File) InitConfig() {

	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(p.configName)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(".")

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println(err.Error())
			return
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	//v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	//bindFlags(cmd, v)

	//return nil

}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
// func bindFlags(cmd *cobra.Command, v *viper.Viper) {
// 	cmd.Flags().VisitAll(func(f *pflag.Flag) {
// 		// Environment variables can't have dashes in them, so bind them to their equivalent
// 		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
// 		if strings.Contains(f.Name, "-") {
// 			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
// 			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
// 		}

// 		// Apply the viper config value to the flag when the flag is not set and viper has a value
// 		if !f.Changed && v.IsSet(f.Name) {
// 			val := v.Get(f.Name)
// 			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
// 		}
// 	})
// }
