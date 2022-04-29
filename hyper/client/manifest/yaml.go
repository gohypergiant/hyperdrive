package manifest

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Manifest struct {
	StudyName   string `yaml:"study_name"`
	ModelFlavor string `yaml:"model_flavor"`
	Training    struct {
		Data struct {
			Features struct {
				Source string `yaml:"source"`
			} `yaml:"features"`
			Target struct {
				Source string `yaml:"source"`
			} `yaml:"target"`
		} `yaml:"data"`
	} `yaml:"training"`
}


func GetManifest(manifestPath string) Manifest {
	var m Manifest
	yamlFile, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return m
}

func GetName(manifestPath string) string {
	return GetManifest(manifestPath).StudyName
}
