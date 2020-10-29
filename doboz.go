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
	initialize()
	if cmd, ok := parsedOpts["help"]; ok {
		if len(cmd[0]) == 0 {
			for _, o := range definedOpts {
				fmt.Printf("%s\n\t%s\n\n", o.Name, o.Description)
			}
			util.TExit(0)
		}
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
		util.TExit(1, "No config yaml specified")
	}

	_, ok = parsedOpts["-package-version"]
	if !ok {
		util.TExit(1, "Version is not set.")
	}

	log.Println("Using configuration from " + conf[0])
	log.Printf("GitLab urls: %s, %s\n", scmUrl, configUrl)

	workspace := getWorkspace()

	module, err := yamlreader.GetModule(scmUrl, configUrl, conf[0])
	util.TExitOnError(err)

	user := util.User{Username: "", Password: ""}

	createSkeleton(module, workspace, user)

	//fmt.Println(commitIds)

	util.TExit(0)
}

func createSkeleton(module *yamlreader.CicdModule, workspace string, user util.User) {
	commitids, err := cloneModule(module, workspace, user)
	util.TExitOnError(err)

	fmt.Println(commitids)

	version, _ := parsedOpts["-package-version"]
	name, ok := parsedOpts["-package-name"]
	pkgName := (*module).Name
	if ok {
		pkgName = name[0]
	}

	dockerPath := createPackageDirs(workspace, pkgName, version)
	createCommitIdTxts(module, dockerPath, commitids)

}

func createCommitIdTxts(module *yamlreader.CicdModule, dockerPath string, commitids *map[string]string) {
	for _, p := range module.Projects {
		path := filepath.Join(dockerPath, p.Name)
		util.CreateDirOrExit(path)
		f, err := os.Create(filepath.Join(path, "commitid.txt"))
		util.TExitOnError(err)

		text := p.Name + ":       " + (*commitids)[p.Name] + "\n" +
			p.Name + "-infra: " + (*commitids)[p.Name+"-infra"]

		fmt.Println(text)

		_, err = f.WriteString(text)
		util.TExitOnError(err)
	}
}

func createPackageDirs(workspace string, pkgName string, version []string) string {
	pkgPath := filepath.Join(workspace, pkgName+"_v"+version[0])
	dbPath := filepath.Join(pkgPath, "Adatbazis")
	dockerPath := filepath.Join(pkgPath, "Docker")

	util.CreateDirOrExit(pkgPath)
	util.CreateDirOrExit(filepath.Join(pkgPath, "1_Dokumentumok"))
	util.CreateDirOrExit(dbPath)

	util.CreateDirOrExit(dockerPath)
	return dockerPath
}

func cloneModule(module *yamlreader.CicdModule, workspace string, user util.User) (*map[string]string, error) {
	commitIds := make(map[string]string)
	scmPath := filepath.Join(workspace, "scm")
	configPath := filepath.Join(workspace, "config")
	cleanWorkspaceIfNotEmpty(workspace, scmPath, configPath)
	for _, p := range module.Projects {
		log.Println("Cloning " + p.Name)
		_, err := git.CloneProject(filepath.Join(scmPath, p.Name), p.ScmUrl, "developer", &user)
		if err != nil {
			return nil, err
		}
		configCID, err := git.CloneProject(filepath.Join(configPath, p.Name), p.ConfigUrl, "erste_d10", &user)
		if err != nil {
			return nil, err
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
	return &commitIds, nil
}

func cleanWorkspaceIfNotEmpty(workspace string, scmPath string, configPath string) {
	if util.DirExists(filepath.Join(workspace, "scm")) || util.DirExists(filepath.Join(workspace, "config")) {
		shallClean := false
		if _, ok := parsedOpts["-clean-workspace"]; ok {
			log.Println("Cleaning workspace.")
			shallClean = true
		} else {
			fmt.Println("****************************************************************************")
			fmt.Println("Your workspace is not empty." +
				"\nIf you proceed without cleaning it, you might not have a fresh repository," +
				"\ncommitIDs might point to an earlier state." +
				"\n(See: -clean-workspace flag.)")
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
			util.TExit(2)
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
	config.NewOption(&definedOpts, "help", "help <command or option>", 1)

	config.NewCommand(&definedOpts, "package", "Create a package")
	//config.NewCommand(&definedOpts, "skeleton", "Create a package skeleton")
	config.NewCommand(&definedOpts, "login", "Store user data until 'clean' is called. See: clean.")
	config.NewCommand(&definedOpts, "clean", "Cleans user data stored by 'login'. See: login.")

	config.NewOption(&definedOpts, "-config", "Location of the configuration properties file.", 1)
	config.NewOption(&definedOpts, "-workspace", "Location of workspace. Workspace is the directory that would hold git repositories.", 1)
	config.NewOption(&definedOpts, "-package-name", "The name of the package. The name of the directory under workspace that will hold your package.", 1)
	config.NewOption(&definedOpts, "-package-version", "The name of the package. The name of the directory under workspace that will hold your package.", 1)
	config.NewOption(&definedOpts, "-clean-workspace", "Automatically cleans workspace if not empty.", 0)

	arguments, err := config.ParseArgs(definedOpts, os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}
	parsedOpts = *arguments
}
