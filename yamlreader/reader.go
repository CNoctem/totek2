package yamlreader

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type cicdmodule struct {
	Module struct {
		Name     string   `yaml:"name"`
		Group    string   `yaml:"group"`
		Projects []string `yaml:"projects"`
	} `yaml:"module"`
}

type CicdModule struct {
	Name     string
	Projects []CicdProject
}

type CicdProject struct {
	Name           string
	ScmUrl         string
	ConfigUrl      string
	ConfigCommitId string
	InfraUrl       string
	InfraCommitId  string
}

func GetModule(scmUrl, configUrl, yamlfile string) (*CicdModule, error) {
	module, err := readyaml(yamlfile)
	if err != nil {
		return nil, err
	}

	group := (*module).Module.Group
	projects := make([]CicdProject, 0)
	for _, m := range (*module).Module.Projects {
		pScmUrl := scmUrl + "/" + group + "/" + m + ".git"
		pConfigfUrl := configUrl + "/" + group + "/" + m + ".git"
		pInfraUrl := configUrl + "/" + group + "/" + m + "-infra.git"
		projects = append(projects, CicdProject{Name: m, ScmUrl: pScmUrl, ConfigUrl: pConfigfUrl, InfraUrl: pInfraUrl})
	}

	return &CicdModule{Name: (*module).Module.Name, Projects: projects}, nil
}

func readyaml(yamlfile string) (*cicdmodule, error) {
	bytes, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		return nil, err
	}
	var config cicdmodule
	err2 := yaml.Unmarshal(bytes, &config)
	if err2 != nil {
		return nil, err
	}
	return &config, nil
}
