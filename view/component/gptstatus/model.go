package gptstatus

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sashabaranov/go-openai"

	"commandgpt/model/event"
	"commandgpt/model/resource"
	"commandgpt/view/atom/label"
	"commandgpt/view/color"
)

var (
	openAIModels = []resource.GPTModel{
		{
			ID:   openai.GPT3Dot5Turbo,
			Name: "GPT 3.5",
		},
		{
			ID:   openai.GPT4,
			Name: "GPT 4",
		},
	}

	titleStyle = lipgloss.NewStyle().
			Bold(true).Padding(0, 1).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color(color.GruvboxFG0)).
			Background(lipgloss.Color(color.GruvboxBlue))

	modelNameStyle = lipgloss.NewStyle().
			Width(11).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(lipgloss.Color(color.GruvboxBG0)).
			Background(lipgloss.Color(color.GruvboxFG0))
)

type Model struct {
	title        label.Model
	modelName    label.Model
	currentModel int
}

func New() Model {
	title := label.New("OpenAI")
	title.Style = titleStyle

	modelName := label.New(openAIModels[1].Name)
	modelName.Style = modelNameStyle

	return Model{
		title:        title,
		modelName:    modelName,
		currentModel: 1,
	}
}

func (m *Model) NextModel() tea.Cmd {
	m.currentModel = (m.currentModel + 1) % len(openAIModels)
	m.modelName.Text = openAIModels[m.currentModel].Name
	return event.SwitchOpenAIModel(openAIModels[m.currentModel].ID)
}

func (m *Model) CurrentModelID() string {
	return openAIModels[m.currentModel].ID
}

func (m *Model) GetSize() (int, int) {
	titleX, titleY := m.title.GetSize()
	modelNameX, modelNameY := m.modelName.GetSize()

	return titleX + modelNameX, titleY + modelNameY
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Center, m.title.View(), m.modelName.View())
}
