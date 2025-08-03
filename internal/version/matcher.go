package version

import (
	"strings"

	"github.com/HMZElidrissi/eol-checker/internal/models"
)

// Matcher handles version matching logic
type Matcher struct{}

// NewMatcher creates a new version matcher
func NewMatcher() *Matcher {
	return &Matcher{}
}

// FindBestMatch finds the best matching EOL cycle for a given version
func (m *Matcher) FindBestMatch(version string, cycles []models.EOLCycle) *models.EOLCycle {
	if version == "" || len(cycles) == 0 {
		return nil
	}

	// Try exact match first
	if cycle := m.findExactMatch(version, cycles); cycle != nil {
		return cycle
	}

	// Try prefix match
	if cycle := m.findPrefixMatch(version, cycles); cycle != nil {
		return cycle
	}

	// Try semantic version match (major.minor)
	if cycle := m.findSemanticMatch(version, cycles); cycle != nil {
		return cycle
	}

	// Try major version match
	return m.findMajorVersionMatch(version, cycles)
}

func (m *Matcher) findExactMatch(version string, cycles []models.EOLCycle) *models.EOLCycle {
	for i, cycle := range cycles {
		if string(cycle.Cycle) == version {
			return &cycles[i]
		}
	}
	return nil
}

func (m *Matcher) findPrefixMatch(version string, cycles []models.EOLCycle) *models.EOLCycle {
	for i, cycle := range cycles {
		if strings.HasPrefix(version, string(cycle.Cycle)) {
			return &cycles[i]
		}
	}
	return nil
}

func (m *Matcher) findSemanticMatch(version string, cycles []models.EOLCycle) *models.EOLCycle {
	versionParts := strings.Split(version, ".")
	if len(versionParts) < 2 {
		return nil
	}

	for i, cycle := range cycles {
		cycleParts := strings.Split(string(cycle.Cycle), ".")
		if len(cycleParts) >= 2 {
			if versionParts[0] == cycleParts[0] && versionParts[1] == cycleParts[1] {
				return &cycles[i]
			}
		}
	}
	return nil
}

func (m *Matcher) findMajorVersionMatch(version string, cycles []models.EOLCycle) *models.EOLCycle {
	versionParts := strings.Split(version, ".")
	if len(versionParts) < 1 {
		return nil
	}

	for i, cycle := range cycles {
		if strings.HasPrefix(string(cycle.Cycle), versionParts[0]+".") {
			return &cycles[i]
		}
	}
	return nil
}
