package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/HMZElidrissi/eol-checker/internal/api"
	"github.com/HMZElidrissi/eol-checker/internal/models"
	"github.com/HMZElidrissi/eol-checker/internal/version"
	"github.com/HMZElidrissi/eol-checker/pkg/image"
)

// Messages
type eolCheckMsg struct {
	result models.EOLResult
	err    error
}

// Model represents the TUI application state
type Model struct {
	textInput      textinput.Model
	spinner        spinner.Model
	loading        bool
	result         *models.EOLResult
	err            error
	width          int
	height         int
	apiClient      *api.Client
	versionMatcher *version.Matcher
	imageParser    *image.Parser
}

// NewModel creates a new TUI model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter image name (e.g., nginx:1.20, ubuntu:20.04, node:16)"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		textInput:      ti,
		spinner:        s,
		loading:        false,
		apiClient:      api.NewClient(),
		versionMatcher: version.NewMatcher(),
		imageParser:    image.NewParser(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.loading && m.textInput.Value() != "" {
				m.loading = true
				m.result = nil
				m.err = nil
				return m, tea.Batch(
					m.spinner.Tick,
					m.checkEOL(m.textInput.Value()),
				)
			} else if !m.loading && m.textInput.Value() == "" && (m.result != nil || m.err != nil) {
				// Clear previous results when Enter is pressed on empty input
				m.result = nil
				m.err = nil
				return m, nil
			}
		}

	case eolCheckMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.result = &msg.result
			// Clear the input field after successful check to allow for another input
			m.textInput.SetValue("")
		}
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.loading {
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

// checkEOL performs the EOL check asynchronously
func (m Model) checkEOL(imageName string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.performEOLCheck(imageName)
		return eolCheckMsg{result: result, err: err}
	}
}

// performEOLCheck performs the actual EOL check
func (m Model) performEOLCheck(imageName string) (models.EOLResult, error) {
	// Parse image name
	imageInfo, err := m.imageParser.Parse(imageName)
	if err != nil {
		return models.EOLResult{}, fmt.Errorf("failed to parse image: %w", err)
	}

	// Fetch EOL data
	cycles, err := m.apiClient.GetProductCycles(imageInfo.Product)
	if err != nil {
		return models.EOLResult{}, fmt.Errorf("failed to fetch EOL data: %w", err)
	}

	if cycles == nil {
		return models.EOLResult{
			Product:     imageInfo.Product,
			Version:     imageInfo.Version,
			Status:      models.StatusUnknown,
			Description: fmt.Sprintf("Product '%s' not found in EOL database", imageInfo.Product),
		}, nil
	}

	// Find the overall latest version (first cycle is typically the most recent)
	var overallLatest string
	if len(cycles) > 0 {
		// Find the cycle with the most recent release date
		var latestCycle *models.EOLCycle
		for i := range cycles {
			if latestCycle == nil {
				latestCycle = &cycles[i]
				continue
			}

			// Compare release dates to find the most recent
			if cycles[i].ReleaseDate > latestCycle.ReleaseDate {
				latestCycle = &cycles[i]
			}
		}

		if latestCycle != nil {
			overallLatest = latestCycle.Latest
		}
	}

	// Find matching cycle
	cycleInfo := m.versionMatcher.FindBestMatch(imageInfo.Version, cycles)
	if cycleInfo == nil {
		return models.EOLResult{
			Product:     imageInfo.Product,
			Version:     imageInfo.Version,
			Status:      models.StatusUnknown,
			Description: fmt.Sprintf("Version '%s' not found for product '%s'", imageInfo.Version, imageInfo.Product),
			Latest:      overallLatest,
		}, nil
	}

	return m.buildEOLResult(imageName, imageInfo, cycleInfo, overallLatest)
}

// buildEOLResult builds the final EOL result with status analysis
func (m Model) buildEOLResult(imageName string, imageInfo *image.ImageInfo, cycleInfo *models.EOLCycle, overallLatest string) (models.EOLResult, error) {
	result := models.EOLResult{
		Product: imageInfo.Product,
		Version: string(cycleInfo.Cycle),
		Latest:  overallLatest,
	}

	if cycleInfo.Link != nil {
		result.Link = *cycleInfo.Link
	}

	// Parse EOL date
	eolDateStr, hasEOL := cycleInfo.EOL.(string)
	var eolDate time.Time
	var eolDateErr error
	if hasEOL {
		result.EOLDate = eolDateStr
		eolDate, eolDateErr = time.Parse("2006-01-02", eolDateStr)
	}

	// Parse support end date
	supportEndStr, hasSupport := cycleInfo.Support.(string)
	var supportEndDate time.Time
	var supportEndDateErr error
	if hasSupport {
		result.SupportEndDate = supportEndStr
		supportEndDate, supportEndDateErr = time.Parse("2006-01-02", supportEndStr)
	}

	// Parse discontinued
	discontinued := false
	if d, ok := cycleInfo.Discontinued.(string); ok {
		if t, err := time.Parse("2006-01-02", d); err == nil && t.Before(time.Now()) {
			discontinued = true
		}
	} else if d, ok := cycleInfo.Discontinued.(bool); ok {
		discontinued = d
	}

	// Calculate days remaining
	result.DaysRemaining = -1
	if hasEOL && eolDateErr == nil {
		result.DaysRemaining = int(time.Until(eolDate).Hours() / 24)
	}

	daysToSupportEnd := -1
	if hasSupport && supportEndDateErr == nil {
		daysToSupportEnd = int(time.Until(supportEndDate).Hours() / 24)
	}

	// Determine status and messages
	now := time.Now()

	if discontinued {
		result.Status = models.StatusCritical
		result.Description = fmt.Sprintf("The image %s is based on a discontinued version of %s.", imageName, imageInfo.Product)
		result.Recommendation = fmt.Sprintf("Upgrade immediately to the latest version (%s) as this version is no longer maintained.", overallLatest)
	} else if hasSupport && supportEndDateErr == nil && supportEndDate.Before(now) {
		result.Status = models.StatusCritical
		result.Description = fmt.Sprintf("The image %s is based on %s which is no longer supported (support ended on %s).", imageName, imageInfo.Product, supportEndStr)
		result.Recommendation = fmt.Sprintf("Upgrade to a supported version. Latest version is %s.", overallLatest)
	} else if hasEOL && eolDateErr == nil && eolDate.Before(now) {
		result.Status = models.StatusCritical
		result.Description = fmt.Sprintf("The image %s is based on %s which reached End-of-Life on %s.", imageName, imageInfo.Product, eolDateStr)
		result.Recommendation = fmt.Sprintf("Upgrade to a newer version. Latest version is %s.", overallLatest)
	} else if hasSupport && supportEndDateErr == nil && daysToSupportEnd <= 30 {
		result.Status = models.StatusWarning
		result.Description = fmt.Sprintf("The image %s is based on %s which will lose support in %d days (on %s).", imageName, imageInfo.Product, daysToSupportEnd, supportEndStr)
		result.Recommendation = fmt.Sprintf("Plan to upgrade soon. Latest version is %s.", overallLatest)
	} else if hasEOL && eolDateErr == nil && result.DaysRemaining <= 30 {
		result.Status = models.StatusWarning
		result.Description = fmt.Sprintf("The image %s is based on %s which will reach End-of-Life in %d days (on %s).", imageName, imageInfo.Product, result.DaysRemaining, eolDateStr)
		result.Recommendation = fmt.Sprintf("Plan to upgrade soon. Latest version is %s.", overallLatest)
	} else if (hasSupport && supportEndDateErr == nil && daysToSupportEnd <= 90) || (hasEOL && eolDateErr == nil && result.DaysRemaining <= 90) {
		result.Status = models.StatusInfo
		if hasSupport && supportEndDateErr == nil && daysToSupportEnd <= 90 {
			result.Description = fmt.Sprintf("The image %s is based on %s which will lose support in %d days (on %s).", imageName, imageInfo.Product, daysToSupportEnd, supportEndStr)
		} else {
			result.Description = fmt.Sprintf("The image %s is based on %s which will reach End-of-Life in %d days (on %s).", imageName, imageInfo.Product, result.DaysRemaining, eolDateStr)
		}
		result.Recommendation = fmt.Sprintf("Consider planning an upgrade. Latest version is %s.", overallLatest)
	} else {
		result.Status = models.StatusOK
		result.Description = fmt.Sprintf("The image %s is based on a currently supported version of %s.", imageName, imageInfo.Product)
		if cycleInfo.Latest != string(cycleInfo.Cycle) {
			result.Recommendation = fmt.Sprintf("This version is supported, but consider upgrading to the latest version (%s) for the newest features and security updates.", overallLatest)
		}
	}

	return result, nil
}

// View renders the TUI
func (m Model) View() string {
	return RenderView(m)
}
