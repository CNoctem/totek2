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
var commitIds = make(map[string]string)

func main() {
	initialize()
	if cmd, ok := parsedOpts["help"]; ok {
		fmt.Println(config.GetDescription(definedOpts, cmd[0]))
		util.TExit(0)
	}

	if _, ok := parsedOpts["package"]; ok {
		createPackage()
	}

	if _, ok := parsedOpts["login"]; ok {
		if !util.AskForConfirmation("Your password will be saved in plain text! Proceed?") {
			util.TExit(2)
		}
		user := util.Authenticate()
		util.SaveUserData(user)
		log.Println("User data has been saved.")
		util.TExit(0)
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

	workspace := getWorkspace()

	module, err := yamlreader.GetModule(scmUrl, configUrl, conf[0])
	if err != nil {
		log.Fatal(err)
	}

	user := util.User{Username: "", Password: ""}

	cloneModule(module, workspace, user)

	fmt.Println(commitIds)

	util.TExit(0)
}

func cloneModule(module *yamlreader.CicdModule, workspace string, user util.User) {
	scmPath := filepath.Join(workspace, "scm")
	configPath := filepath.Join(workspace, "config")
	cleanWorkspaceIfNotEmpty(workspace, scmPath, configPath)
	for _, p := range module.Projects {
		log.Println("Cloning " + p.Name)
		_, err := git.CloneProject(filepath.Join(scmPath, p.Name), p.ScmUrl, "developer", &user)
		if err != nil {
			log.Fatal(err)
		}
		configCID, err := git.CloneProject(filepath.Join(configPath, p.Name), p.ConfigUrl, "erste_d10", &user)
		if err != nil {
			log.Fatal(err)
		}
		p.ConfigCommitId = *configCID
		commitIds[p.Name] = *configCID
		infraCID, err := git.CloneProject(filepath.Join(configPath, p.Name+"-infra"), p.InfraUrl, "erste_d10", &user)
		if err != nil {
			log.Fatal(err)
		}
		p.InfraCommitId = *infraCID
		commitIds[p.Name+"-infra"] = *infraCID
	}
}

func cleanWorkspaceIfNotEmpty(workspace string, scmPath string, configPath string) {
	if util.DirExists(filepath.Join(workspace, "scm")) || util.DirExists(filepath.Join(workspace, "config")) {
		shallClean := false
		if _, ok := parsedOpts["-clean"]; ok {
			log.Println("Cleaning workspace.")
			shallClean = true
		} else {
			fmt.Println("****************************************************************************")
			fmt.Println("Your workspace is not empty." +
				"\nIf you proceed without cleaning it, you might not have a fresh repository," +
				"\ncommitIDs might point to an earlier state." +
				"\n(See: -clean flag.)")
			fmt.Println("****************************************************************************")
			shallClean = util.AskForConfirmation("Do you want me to clean your workspace?")
		}
		if shallClean {
			if err := os.RemoveAll(scmPath); err != nil {
				log.Fatal(err)
			}
			if err := os.RemoveAll(configPath); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getWorkspace() string {
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
	return workspace
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

func initialize() {
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
	config.NewCommand(&definedOpts, "package", "Create a package")
	config.NewCommand(&definedOpts, "skeleton", "Create a package skeleton")
	config.NewCommand(&definedOpts, "login", "Store user data until 'clean' is called. See: clean.")
	config.NewCommand(&definedOpts, "clean", "Cleans user data stored by 'login'. See: login.")
	config.NewOption(&definedOpts, "help", "help", 1)
	config.NewOption(&definedOpts, "-config", "Location of the configuration properties file.", 1)
	config.NewOption(&definedOpts, "-workspace", "Location of workspace. Workspace is the directory that would hold git repositories.", 1)
	config.NewOption(&definedOpts, "-package-name", "The name of the package. The name of the directory under workspace that will hold your package.", 1)
	config.NewOption(&definedOpts, "-package-version", "The name of the package. The name of the directory under workspace that will hold your package.", 1)
	config.NewOption(&definedOpts, "-clean", "Automatically cleans workspace if not empty.", 0)

	arguments, err := config.ParseArgs(definedOpts, os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
	parsedOpts = *arguments
}
