package cache

import "github.com/hydrocode-de/gorun/internal/toolSpec"

type Cache struct {
	images map[string]toolSpec.SpecFile
	data   map[string]toolSpec.ToolSpec
}

func (c *Cache) GetToolSpec(key string) (*toolSpec.ToolSpec, bool) {
	spec, ok := c.data[key]
	return &spec, ok
}

func (c *Cache) SetToolSpec(key string, spec *toolSpec.ToolSpec) {
	c.data[key] = *spec
}

func (c *Cache) GetImageSpec(key string) (*toolSpec.SpecFile, bool) {
	spec, ok := c.images[key]
	return &spec, ok
}

func (c *Cache) SetImageSpec(key string, spec toolSpec.SpecFile) {
	c.images[key] = spec
}

func (c *Cache) Reset() {
	c.data = make(map[string]toolSpec.ToolSpec)
	c.images = make(map[string]toolSpec.SpecFile)
}
