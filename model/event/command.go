package event

import tea "github.com/charmbracelet/bubbletea"

type CommandConfirmedMsg string

func ConfirmCommand(command string) tea.Cmd {
	return func() tea.Msg {
		return CommandConfirmedMsg(command)
	}
}
