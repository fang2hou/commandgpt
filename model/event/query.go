package event

import tea "github.com/charmbracelet/bubbletea"

type QuerySentMsg string

func SendQuery(query string) tea.Cmd {
	return func() tea.Msg {
		return QuerySentMsg(query)
	}
}
