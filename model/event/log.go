package event

import tea "github.com/charmbracelet/bubbletea"

type LogErrorMsg string

func LogError(err error) tea.Cmd {
	return func() tea.Msg {
		return LogErrorMsg(err.Error())
	}
}
