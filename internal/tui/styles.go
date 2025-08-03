package tui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(0, 1).
			MarginBottom(1)

	ResultStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			MarginTop(1)

	CriticalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF0000")).
			Padding(0, 1).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFA500")).
			Padding(0, 1).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0066CC")).
			Padding(0, 1).
			Bold(true)

	OKStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#008000")).
		Padding(0, 1).
		Bold(true)

	UnknownStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#666666")).
			Padding(0, 1).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	BoldStyle = lipgloss.NewStyle().Bold(true)

	LinkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#0066CC"))
)

// GetStatusStyle returns the appropriate style for a status
func GetStatusStyle(status string) (lipgloss.Style, string) {
	switch status {
	case "CRITICAL":
		return CriticalStyle, "🚨"
	case "WARNING":
		return WarningStyle, "⚠️"
	case "INFO":
		return InfoStyle, "ℹ️"
	case "OK":
		return OKStyle, "✅"
	default:
		return UnknownStyle, "❓"
	}
}
