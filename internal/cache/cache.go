package cache

import (
	"github.com/hydrocode-de/tool-spec-go"
)

type Cache struct {
	images      map[string]toolspec.SpecFile
	tools       map[string]toolspec.ToolSpec
	Initialised bool
}

func (c *Cache) GetToolSpec(key string) (*toolspec.ToolSpec, bool) {
	spec, ok := c.tools[key]
	return &spec, ok
}

func (c *Cache) SetToolSpec(key string, spec *toolspec.ToolSpec) {
	c.tools[key] = *spec
}

func (c *Cache) ListToolSpecs() []toolspec.ToolSpec {
	specs := make([]toolspec.ToolSpec, 0)
	for _, spec := range c.tools {
		specs = append(specs, spec)
	}
	return specs
}

func (c *Cache) GetImageSpec(key string) (*toolspec.SpecFile, bool) {
	spec, ok := c.images[key]
	return &spec, ok
}

func (c *Cache) SetImageSpec(key string, spec toolspec.SpecFile) {
	c.images[key] = spec
}

func (c *Cache) Reset() {
	c.tools = make(map[string]toolspec.ToolSpec)
	c.images = make(map[string]toolspec.SpecFile)
	c.Initialised = false
}
