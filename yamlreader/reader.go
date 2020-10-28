package yamlreader

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type cicdmodule struct {
	Module   struct {
		Group      string `yaml:"group"`
		Projects []string `yaml:"projects"`
	} `yaml:"module"`
}

type CicdModule struct {
	ConfigProjects []string
	ScmProjects []string
}

func GetModule(scmUrl, configUrl, yamlfile string) (*CicdModule, error) {
	module, err := readyaml(yamlfile)
	if err != nil {
		return nil, err
	}

	group := (*module).Module.Group
	configProjects := make([]string, 0)
	scmProjects := make([]string, 0)
	for _, m := range (*module).Module.Projects {
		scmProjects = append(scmProjects, scmUrl + "/" + group + "/" + m + ".git")
		configProjects = append(configProjects, configUrl + "/" + group + "/" + m + ".git")
		configProjects = append(configProjects, configUrl + "/" + group + "/" + m + "-infra.git")
	}
	return &CicdModule{ConfigProjects: configProjects, ScmProjects: scmProjects}, nil
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


