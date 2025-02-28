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
	ID          string                   `json:"id" yaml:"-"`
	Name        string                   `json:"name" yaml:"-"`
	Title       string                   `json:"title" yaml:"title"`
	Description string                   `json:"description" yaml:"description"`
	Parameters  map[string]ParameterSpec `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Data        map[string]DataSpec      `json:"data,omitempty" yaml:"data,omitempty"`
}

type ParameterSpec struct {
	Name        string   `json:"name" yaml:"name"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	ToolType    string   `json:"type" yaml:"type"`
	IsArray     bool     `json:"array,omitempty" yaml:"array,omitempty" default:"false"`
	Default     bool     `json:"default,omitempty" yaml:"default,omitempty" default:"false"`
	Values      []string `json:"values,omitempty" yaml:"values,omitempty"`
	Min         float64  `json:"min,omitempty" yaml:"min,omitempty"`
	Max         float64  `json:"max,omitempty" yaml:"max,omitempty"`
}

type DataSpec struct {
	Path        string      `json:"path" yaml:"path"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Example     string      `json:"example,omitempty" yaml:"example,omitempty"`
	Extension   interface{} `json:"-" yaml:"extension,omitempty"`
	Extensions  []string    `json:"extension,omitempty"`
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
