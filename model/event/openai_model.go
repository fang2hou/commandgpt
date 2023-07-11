package event

import tea "github.com/charmbracelet/bubbletea"

type OpenAIModelChangedMsg string

func SwitchOpenAIModel(modelID string) tea.Cmd {
	return func() tea.Msg {
		return OpenAIModelChangedMsg(modelID)
	}
}
