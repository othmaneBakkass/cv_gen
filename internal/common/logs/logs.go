package logs

import "github.com/charmbracelet/lipgloss"

var successLogStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1"))

func SuccessLog(msg string) string {
	return successLogStyle.Render(msg)
}

var errorLogStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8"))

func ErrorLog(msg string) string {
	return errorLogStyle.Render(msg)
}

var infoLogStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa"))

func InfoLog(msg string) string {
	return infoLogStyle.Render(msg)
}
