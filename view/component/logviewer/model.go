package logviewer

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"commandgpt/view/color"
)

var (
	containerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(color.GruvboxBG0))

	spinnerStyle = lipgloss.NewStyle()
)

type Model struct {
	style          lipgloss.Style
	logs           string
	waiting        bool
	waitingMessage string
	spinner        spinner.Model
}

func New() Model {
	spn := spinner.New()
	spn.Style = spinnerStyle
	spn.Spinner = spinner.Points

	return Model{
		style: lipgloss.NewStyle().
			Padding(1).
			Foreground(lipgloss.Color(color.GruvboxFG0)).
			Background(lipgloss.Color(color.GruvboxBG0)),
		logs:    "",
		spinner: spn,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m *Model) addLog(emoji, text string) {
	if len(m.logs) > 0 {
		m.logs += "\n"
	}
	m.logs += emoji + " " + text
}

func (m *Model) LogError(err string) {
	m.addLog("❌", err)
}

func (m *Model) LogInfo(info string) {
	m.addLog("ℹ️ ", info)
}

func (m *Model) LogSuccess(success string) {
	m.addLog("✅", success)
}

func (m *Model) Wait(msg string) {
	m.waiting = true
	m.waitingMessage = msg
}

func (m *Model) Done() {
	m.waiting = false
	m.waitingMessage = ""
}

func (m *Model) SetWidth(width int) {
	m.style = m.style.Width(width)
}

func (m *Model) SetHeight(height int) {
	m.style = m.style.Height(height)
}

func (m *Model) GetSize() (int, int) {
	return m.style.GetWidth(), m.style.GetHeight()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if m.waiting {
		switch msg := msg.(type) {
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return cmd
		}
	}

	return nil
}

func (m *Model) View() string {
	s := m.logs

	if m.waiting {
		if len(s) > 0 {
			s += "\n\n"
		}
		s += m.spinner.View() + " " + m.waitingMessage
	}

	return containerStyle.Render(m.style.Render(s))
}
