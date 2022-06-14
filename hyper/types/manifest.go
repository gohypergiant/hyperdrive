package types

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
