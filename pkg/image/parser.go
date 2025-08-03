package image

import (
	"fmt"
	"path"
	"strings"
)

// ImageInfo represents parsed container image information
type ImageInfo struct {
	Registry string
	Name     string
	Tag      string
	Product  string
	Version  string
}

// Parser handles container image name parsing
type Parser struct{}

// NewParser creates a new image parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a container image name into its components
func (p *Parser) Parse(imageName string) (*ImageInfo, error) {
	if imageName == "" {
		return nil, fmt.Errorf("image name cannot be empty")
	}

	// Handle registry prefix (optional)
	parts := strings.Split(imageName, "/")
	var registry, nameWithTag string

	if len(parts) > 2 || (len(parts) == 2 && strings.Contains(parts[0], ".")) {
		// Has registry
		registry = parts[0]
		nameWithTag = strings.Join(parts[1:], "/")
	} else {
		// No registry, assume Docker Hub
		nameWithTag = imageName
	}

	// Split name and tag
	nameAndTag := strings.Split(nameWithTag, ":")
	name := nameAndTag[0]
	tag := "latest"

	if len(nameAndTag) > 1 {
		tag = nameAndTag[1]
	}

	// Extract product name (base name without path)
	product := path.Base(name)

	// Extract version from tag (remove suffixes like -alpine, -slim)
	version := strings.SplitN(tag, "-", 2)[0]
	if version == "latest" {
		version = ""
	}

	return &ImageInfo{
		Registry: registry,
		Name:     name,
		Tag:      tag,
		Product:  product,
		Version:  version,
	}, nil
}
