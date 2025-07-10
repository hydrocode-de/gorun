package cache

import (
	"github.com/hydrocode-de/gorun/internal/toolSpec"
)

type Cache struct {
	images      map[string]toolSpec.SpecFile
	tools       map[string]toolSpec.ToolSpec
	Initialised bool
}

func (c *Cache) GetToolSpec(key string) (*toolSpec.ToolSpec, bool) {
	spec, ok := c.tools[key]
	return &spec, ok
}

func (c *Cache) SetToolSpec(key string, spec *toolSpec.ToolSpec) {
	c.tools[key] = *spec
}

func (c *Cache) ListToolSpecs() []toolSpec.ToolSpec {
	specs := make([]toolSpec.ToolSpec, 0)
	for _, spec := range c.tools {
		specs = append(specs, spec)
	}
	return specs
}

func (c *Cache) GetImageSpec(key string) (*toolSpec.SpecFile, bool) {
	spec, ok := c.images[key]
	return &spec, ok
}

func (c *Cache) SetImageSpec(key string, spec toolSpec.SpecFile) {
	c.images[key] = spec
}

func (c *Cache) Reset() {
	c.tools = make(map[string]toolSpec.ToolSpec)
	c.images = make(map[string]toolSpec.SpecFile)
	c.Initialised = false
}
