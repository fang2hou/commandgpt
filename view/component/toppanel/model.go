package toppanel

import (
	"github.com/charmbracelet/lipgloss"

	"commandgpt/view/atom/label"
	"commandgpt/view/color"
)

const (
	queryHelpText   = "enter: confirm • ctrl+q: toggle input • ctrl+o: switch openai model • ctrl+c: quit"
	logHelpText     = "ctrl+q: toggle input • ctrl+o: switch openai model • ctrl+c: quit"
	advicesHelpText = "j/k, ↑/↓: select • enter: confirm • ctrl+q: toggle input • ctrl+c: quit"
)

var (
	labelStyle = lipgloss.NewStyle().
			Bold(true).Padding(0, 1).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color(color.GruvboxFG0)).
			Background(lipgloss.Color(color.GruvboxPurple))

	helpTextStyle = lipgloss.NewStyle().
			Align(lipgloss.Right).
			Foreground(lipgloss.Color(color.GruvboxFG4)).
			Background(lipgloss.Color(color.GruvboxBG0))

	modes = map[string]string{
		"Query":   queryHelpText,
		"Log":     logHelpText,
		"Advices": advicesHelpText,
	}
)

type Model struct {
	label         label.Model
	helpText      string
	helpTextStyle lipgloss.Style
	width         int
}

func New() Model {
	l := label.New("Query")
	l.Style = labelStyle
	return Model{
		label:         l,
		helpTextStyle: helpTextStyle,
	}
}

func (m *Model) SetWidth(w int) {
	m.width = w
	fx, _ := m.label.Style.GetFrameSize()
	m.helpTextStyle = m.helpTextStyle.Width(w - fx - len(m.label.Text))
}

func (m *Model) View(mode string) string {
	if v, ok := modes[mode]; ok {
		m.label.Text = mode
		m.helpText = v
		m.SetWidth(m.width)
	}

	ls := m.label.View()
	hts := m.helpTextStyle.Render(m.helpText)
	return lipgloss.JoinHorizontal(lipgloss.Top, ls, hts)
}
