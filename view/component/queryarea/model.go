package queryarea

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"commandgpt/model/event"
	"commandgpt/view/atom/label"
	"commandgpt/view/color"
)

var (
	titleStyle = lipgloss.NewStyle().
			Width(7).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color(color.GruvboxFG0)).
			Background(lipgloss.Color(color.GruvboxAqua))

	blurredTextAreaStyle = textarea.Style{
		Base:        lipgloss.NewStyle().Background(lipgloss.Color(color.GruvboxBG0)).Padding(0, 1),
		CursorLine:  lipgloss.NewStyle().Background(lipgloss.Color(color.GruvboxBG0)),
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color(color.GruvboxBG3)),
		Text:        lipgloss.NewStyle().Foreground(lipgloss.Color(color.GruvboxFG0)),
	}

	focusedTextAreaStyle = textarea.Style{
		Base:        lipgloss.NewStyle().Background(lipgloss.Color(color.GruvboxBG2)).Padding(0, 1),
		CursorLine:  lipgloss.NewStyle().Background(lipgloss.Color(color.GruvboxBG2)),
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color(color.GruvboxBG3)),
		Text:        lipgloss.NewStyle().Foreground(lipgloss.Color(color.GruvboxFG0)),
	}
)

type Model struct {
	title label.Model
	query textarea.Model
}

func New() Model {
	title := label.New("QUERY")
	title.Style = titleStyle

	query := textarea.New()
	query.Placeholder = "Build dockerfile..."
	query.Prompt = ""
	query.CharLimit = 280
	query.SetWidth(80)
	query.SetHeight(1)
	query.ShowLineNumbers = false
	query.BlurredStyle = blurredTextAreaStyle
	query.FocusedStyle = focusedTextAreaStyle

	return Model{
		title: title,
		query: query,
	}
}

func (m *Model) SetFocus(isFocused bool) {
	if isFocused {
		m.query.Focus()
	} else {
		m.query.Blur()
	}
}

func (m *Model) IsFocused() bool {
	return m.query.Focused()
}

func (m *Model) SetWidth(w int) {
	labelWidth, _ := m.title.GetSize()
	m.query.SetWidth(w - labelWidth)
}

func (m *Model) SetHeight(h int) {
	m.title.Style = m.title.Style.Height(h)
	m.query.SetHeight(h)
}

func (m *Model) GetSize() (int, int) {
	tx, ty := m.title.GetSize()
	return tx + m.query.Width(), ty
}

func (m *Model) Query() string {
	s := m.query.Value()
	if s == "" {
		s = m.query.Placeholder
	}
	return s
}

func (m *Model) Reset() {
	m.query.Reset()
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var queryCmd tea.Cmd

	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		if m.IsFocused() {
			m.SetFocus(false)
			return event.SendQuery(m.Query())
		}
		return nil
	}

	m.query, queryCmd = m.query.Update(msg)
	return queryCmd
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, m.title.View(), m.query.View())
}
