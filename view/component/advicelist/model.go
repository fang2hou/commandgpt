package advicelist

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"commandgpt/model/event"
	"commandgpt/model/resource"
	"commandgpt/view/color"
)

var (
	keyMap = KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
	}

	containerStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(color.GruvboxBG0))

	adviceCommandStyle = lipgloss.NewStyle().
				Bold(true).
				Align(lipgloss.Left).
				Foreground(lipgloss.Color(color.GruvboxFG0)).
				Background(lipgloss.Color(color.GruvboxBG0))

	adviceCommandSelectedStyle = adviceCommandStyle.Copy().
					Foreground(lipgloss.Color(color.GruvboxLighterYellow))

	adviceDescriptionStyle = lipgloss.NewStyle().
				Align(lipgloss.Left).
				Foreground(lipgloss.Color(color.GruvboxFG2)).
				Background(lipgloss.Color(color.GruvboxBG0))

	adviceDescriptionSelectedStyle = adviceDescriptionStyle.Copy().
					Foreground(lipgloss.Color(color.GruvboxYellow))

	indicatorStyle = lipgloss.NewStyle()

	indicatorSelectedStyle = indicatorStyle.Copy().
				Foreground(lipgloss.Color(color.GruvboxYellow))

	backgroundStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(color.GruvboxBG0))

	objectStyle = lipgloss.NewStyle().
			Padding(0, 0, 1, 0)
)

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Confirm key.Binding
}

type Model struct {
	style               lipgloss.Style
	advices             *resource.Advices
	selectedAdviceIndex int
}

func New() Model {
	return Model{
		style: lipgloss.NewStyle().
			Padding(1).
			Foreground(lipgloss.Color(color.GruvboxFG0)).
			Background(lipgloss.Color(color.GruvboxBG0)),
	}
}

func (m *Model) SetAdvices(advices *resource.Advices) {
	m.advices = advices
	m.selectedAdviceIndex = 0
}

func (m *Model) Next() {
	if m.advices == nil || len(m.advices.Advices) == 0 {
		return
	}

	m.selectedAdviceIndex = (m.selectedAdviceIndex + 1) % len(m.advices.Advices)
}

func (m *Model) Previous() {
	if m.advices == nil || len(m.advices.Advices) == 0 {
		return
	}

	m.selectedAdviceIndex = (m.selectedAdviceIndex - 1 + len(m.advices.Advices)) % len(
		m.advices.Advices,
	)
}

func (m *Model) SelectedAdvice() (*resource.Advice, error) {
	if m.advices == nil || len(m.advices.Advices) == 0 {
		return nil, errors.New("no advice")
	}

	return &m.advices.Advices[m.selectedAdviceIndex], nil
}

func (m *Model) SetHeight(height int) {
	m.style = m.style.Height(height)
}

func (m *Model) SetWidth(width int) {
	m.style = m.style.Width(width)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keyMap.Confirm):
			if m.advices == nil || len(m.advices.Advices) == 0 ||
				m.selectedAdviceIndex < 0 || m.selectedAdviceIndex >= len(m.advices.Advices) {
				return nil
			}
			return event.ConfirmCommand(m.advices.Advices[m.selectedAdviceIndex].Command)
		case key.Matches(msg, keyMap.Up):
			if m.selectedAdviceIndex > 0 {
				m.selectedAdviceIndex--
			}
		case key.Matches(msg, keyMap.Down):
			if m.selectedAdviceIndex < len(m.advices.Advices)-1 {
				m.selectedAdviceIndex++
			}
		}
	}

	return nil
}

func (m *Model) viewAdvice(advice resource.Advice, selected bool) string {
	var cs, is, ds lipgloss.Style

	fx, _ := m.style.GetFrameSize()
	indicatorChar := ' '

	if selected {
		indicatorChar = '│'
		cs = adviceCommandSelectedStyle
		ds = adviceDescriptionSelectedStyle
		is = indicatorSelectedStyle
	} else {
		cs = adviceCommandStyle
		ds = adviceDescriptionStyle
		is = indicatorStyle
	}

	indicator := is.Render(string(indicatorChar) + " ")

	textWidth := m.style.GetWidth() - fx - 2

	cmd := advice.Command
	if len(cmd) > textWidth {
		cmd = cmd[:textWidth-3] + "..."
	}
	cmdLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		indicator,
		cs.Render(cmd),
		backgroundStyle.Render(strings.Repeat(" ", textWidth-len(cmd))),
	)

	desc := advice.Description
	if len(desc) > textWidth {
		desc = desc[:textWidth-3] + "..."
	}
	descLine := lipgloss.JoinHorizontal(
		lipgloss.Center,
		indicator,
		ds.Render(desc),
		backgroundStyle.Render(strings.Repeat(" ", textWidth-len(desc))),
	)

	s := lipgloss.JoinVertical(lipgloss.Left, cmdLine, descLine)

	return objectStyle.Render(s)
}

func (m *Model) View() string {
	content := ""

	if m.advices != nil && len(m.advices.Advices) > 0 {
		advices := make([]string, 0, len(m.advices.Advices))
		for i, advice := range m.advices.Advices {
			advices = append(advices, m.viewAdvice(advice, i == m.selectedAdviceIndex))
		}

		content = lipgloss.JoinVertical(lipgloss.Top, advices...)
	}

	return containerStyle.Render(m.style.Render(content))
}
