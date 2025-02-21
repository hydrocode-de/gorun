package toolSpec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type SpecFile struct {
	Tools map[string]ToolSpec `yaml:"tools"`
}

func (s *SpecFile) GetTool(toolName string) (ToolSpec, error) {
	toolSpec, ok := s.Tools[toolName]
	if !ok {
		return ToolSpec{}, fmt.Errorf("tool %s was not found in the given specification file", toolName)
	}

	return toolSpec, nil
}

type ToolSpec struct {
	Title       string                   `yaml:"title"`
	Description string                   `yaml:"description,omitempty"`
	Parameters  map[string]ParameterSpec `yaml:"parameters,omitempty"`
	Data        map[string]DataSpec      `yaml:"data,omitempty"`
}

type ParameterSpec struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description,omitempty"`
	ToolType    string   `yaml:"type"`
	IsArray     bool     `yaml:"array,omitempty" default:"false"`
	Default     bool     `yaml:"default,omitempty" default:"false"`
	Values      []string `yaml:"values,omitempty"`
	Min         float64  `yaml:"min,omitempty"`
	Max         float64  `yaml:"max,omitempty"`
}

type DataSpec struct {
	Path        string      `yaml:"path"`
	Description string      `yaml:"description,omitempty"`
	Example     string      `yaml:"example,omitempty"`
	Extension   interface{} `yaml:"extension,omitempty"`
	Extensions  []string
}

func (d *DataSpec) UnmarshalYAML(value *yaml.Node) error {
	// unmarshal the dataSpecAlias
	type dataSpecAlias DataSpec
	var alias dataSpecAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}

	*d = DataSpec(alias)

	// switch the extension data type
	switch ext := d.Extension.(type) {
	case string:
		d.Extensions = []string{ext}
	case []interface{}:
		d.Extensions = make([]string, len(ext))
		for i, e := range ext {
			d.Extensions[i] = e.(string)
		}
	}

	d.Extension = nil
	return nil
}

func LoadToolSpec(rawData []byte) (SpecFile, error) {
	var toolSpec SpecFile
	err := yaml.Unmarshal(rawData, &toolSpec)
	if err != nil {
		return SpecFile{}, err
	}

	return toolSpec, nil
}
