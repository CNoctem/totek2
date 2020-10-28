package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"totek2/config"
	"totek2/git"
	"totek2/util"
	"totek2/yamlreader"
)

var definedOpts []config.Option
var parsedOpts map[string][]string
var properties map[string]string

func main() {
	config.NewCommand(&definedOpts, "package", "Create a package")
	config.NewCommand(&definedOpts, "skeleton", "Create a package skeleton")
	config.NewOption(&definedOpts, "help", "help", 1)
	config.NewOption(&definedOpts, "-config", "Location of the configuration properties file.", 1)
	config.NewOption(&definedOpts, "-workspace", "Location of workspace. Workspace is the directory that would hold git repositories.", 1)

	arguments, err := config.ParseArgs(definedOpts, os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
	parsedOpts = *arguments
	if cmd, ok := parsedOpts["help"]; ok {
		fmt.Println(config.GetDescription(definedOpts, cmd[0]))
		os.Exit(0)
	}

	if _, ok := parsedOpts["package"]; ok {
		createPackage()
	}

}

func createPackage() {
	fmt.Println("creating package")
	readConfig()
	scmUrl := getConfig("scm.cicd.url")
	configUrl := getConfig("config.cicd.url")
	conf, ok := parsedOpts["-config"]

	if !ok {
		log.Fatal("No config yaml specified")
	}

	log.Println("Using configuration from " + conf[0])
	log.Printf("GitLab urls: %s, %s\n", scmUrl, configUrl)

	var workspace string

	value, ok := parsedOpts["-workspace"]
	if ok {
		workspace = value[0]
	} else {
		log.Println("WARNING: no workspace location is given. Current working dir will be used.")
		if !util.AskForConfirmation("Proceed?") {
			fmt.Print("- Háromba vágtad, édes, jó Lajosom?")
			os.Exit(1)
		}
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		workspace = wd
	}
	log.Println("Using workspace " + workspace)

	module, err := yamlreader.GetModule(scmUrl, configUrl, conf[0])
	if err != nil {
		log.Fatal(err)
	}

	user := util.User{Username: "", Password: ""}

	for _, m := range module.ScmProjects {
		log.Println("Cloning " + m)
		err := git.CloneProject(workspace, m, "developer", &user)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, m := range module.ConfigProjects {
		log.Println("Cloning " + m)
		err := git.CloneProject(workspace, m, "developer", &user)
		if err != nil {
			log.Fatal(err)
		}
	}


	os.Exit(0)
}


func getConfig(key string) string {
	val, ok := properties[key]
	if ok {
		return val
	}
	fmt.Printf("Could not find property: %s\n", key)
	return ""
}

func readConfig() {
	configPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configPath = filepath.Join(configPath,"totek.properties")
	if conf, ok := properties["-util"]; ok {
		configPath = conf
	}

	properties, err = util.GetConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
}



