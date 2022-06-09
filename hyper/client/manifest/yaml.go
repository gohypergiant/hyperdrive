package manifest

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Manifest struct {
	StudyName   string `yaml:"study_name"`
	ModelFlavor string `yaml:"model_flavor"`
	ProjectName string `yaml:"project_name"`
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
		if os.IsNotExist(err) {
			file, err := os.Create(manifestPath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		} else {
			log.Printf("yamlFile.Get err   #%v ", err)
		}

	}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if m.ProjectName == "" {
		m.ProjectName = fmt.Sprintf("project-%s", generatePathNameString())

		str := string(yamlFile)

		updatedYAML := fmt.Sprintf("project_name: %s\n%s", m.ProjectName, str)
		err := ioutil.WriteFile(manifestPath, []byte(updatedYAML), 0)

		if err != nil {
			log.Printf("yamlFile.Write err   #%v ", err)
		}

	}
	if m.StudyName == "" {
		m.StudyName = fmt.Sprintf("study-%s", generatePathNameString())

		str := string(yamlFile)

		updatedYAML := fmt.Sprintf("study_name: %s\n%s", m.StudyName, str)
		err := ioutil.WriteFile(manifestPath, []byte(updatedYAML), 0)

		if err != nil {
			log.Printf("yamlFile.Write err   #%v ", err)
		}

	}

	return m
}
func generatePathNameString() string {

	cwdPath, _ := os.Getwd()
	cwdName := strings.Replace(cwdPath, "/", "-", -1)
	cwdName = strings.Replace(cwdName, "\\", "-", -1)
	cwdName = strings.Replace(cwdName, ":", "-", -1)
	return cwdName
}

func GetName(manifestPath string) string {
	return GetManifest(manifestPath).StudyName
}

func GetProjectName(manifestPath string) string {
	return GetManifest(manifestPath).ProjectName
}
