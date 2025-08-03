package tui

import (
	"fmt"
	"strings"

	"github.com/HMZElidrissi/eol-checker/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// RenderView renders the main TUI view
func RenderView(m Model) string {
	var s strings.Builder

	// Title
	s.WriteString(TitleStyle.Render("🔍 Container Image EOL Checker"))
	s.WriteString("\n\n")

	// Input
	s.WriteString("Enter a container image to check its End-of-Life status:\n")
	s.WriteString(InputStyle.Render(m.textInput.View()))
	s.WriteString("\n")

	// Loading state
	if m.loading {
		s.WriteString(fmt.Sprintf("%s Checking EOL status...", m.spinner.View()))
		s.WriteString("\n")
	}

	// Error state
	if m.err != nil {
		s.WriteString(ErrorStyle.Render(fmt.Sprintf("❌ Error: %v", m.err)))
		s.WriteString("\n")
	}

	// Result
	if m.result != nil {
		s.WriteString(RenderResult(*m.result))
		s.WriteString("\n")
	}

	// Help
	if m.result != nil || m.err != nil {
		s.WriteString(HelpStyle.Render("Enter another image to check • Ctrl+C to quit"))
	} else {
		s.WriteString(HelpStyle.Render("Press Enter to check • Ctrl+C to quit"))
	}

	return s.String()
}

// RenderResult renders the EOL check result
func RenderResult(result models.EOLResult) string {
	var s strings.Builder

	// Status badge
	statusStyle, statusIcon := GetStatusStyle(result.Status)
	header := fmt.Sprintf("%s %s", statusIcon, statusStyle.Render(result.Status))
	s.WriteString(header)
	s.WriteString("\n\n")

	// Product info
	s.WriteString(BoldStyle.Render("Product: "))
	s.WriteString(result.Product)
	if result.Version != "" {
		s.WriteString(fmt.Sprintf(" (version: %s)", result.Version))
	}
	s.WriteString("\n\n")

	// Description
	s.WriteString(BoldStyle.Render("Description:"))
	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Width(80).Render(result.Description))
	s.WriteString("\n\n")

	// Additional info
	if result.EOLDate != "" && result.EOLDate != "false" {
		s.WriteString(BoldStyle.Render("EOL Date: "))
		s.WriteString(result.EOLDate)
		if result.DaysRemaining >= 0 {
			s.WriteString(fmt.Sprintf(" (%d days remaining)", result.DaysRemaining))
		}
		s.WriteString("\n")
	}

	if result.SupportEndDate != "" && result.SupportEndDate != "false" {
		s.WriteString(BoldStyle.Render("Support End: "))
		s.WriteString(result.SupportEndDate)
		s.WriteString("\n")
	}

	if result.Latest != "" {
		s.WriteString(BoldStyle.Render("Latest Version: "))
		s.WriteString(result.Latest)
		s.WriteString("\n")
	}

	// Recommendation
	if result.Recommendation != "" {
		s.WriteString("\n")
		s.WriteString(BoldStyle.Render("Recommendation:"))
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Width(80).Render(result.Recommendation))
		s.WriteString("\n")
	}

	// Link
	if result.Link != "" {
		s.WriteString("\n")
		s.WriteString(BoldStyle.Render("More Info: "))
		s.WriteString(LinkStyle.Render(result.Link))
		s.WriteString("\n")
	}

	return ResultStyle.Render(s.String())
}
