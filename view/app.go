package view

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"commandgpt/controller"
	"commandgpt/model/event"
	"commandgpt/view/color"
	"commandgpt/view/component/advicelist"
	"commandgpt/view/component/gptstatus"
	"commandgpt/view/component/logviewer"
	"commandgpt/view/component/queryarea"
	"commandgpt/view/component/toppanel"
)

var (
	appStyle = lipgloss.NewStyle().Background(lipgloss.Color(color.GruvboxBG0))

	mainContainerStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(color.GruvboxBG0)).
				Padding(1, 2)

	isWaiting = false
)

type Model struct {
	focusOnQuery bool
	showLog      bool

	topPanel   toppanel.Model
	logViewer  logviewer.Model
	adviceList advicelist.Model
	queryArea  queryarea.Model
	gptStatus  gptstatus.Model

	controller controller.Controller

	Result string
}

func New() *Model {
	qa := queryarea.New()
	qa.SetFocus(true)

	return &Model{
		focusOnQuery: true,
		showLog:      true,

		topPanel:   toppanel.New(),
		logViewer:  logviewer.New(),
		adviceList: advicelist.New(),
		queryArea:  qa,
		gptStatus:  gptstatus.New(),

		controller: controller.New(),
	}
}

func (m *Model) Init() tea.Cmd {
	err := m.controller.Init()
	if err != nil {
		return event.LogError(err)
	}

	return nil
}

func (m *Model) updateLayout(x, y int) {
	cfx, cfy := mainContainerStyle.GetFrameSize()
	cx, cy := x-cfx, y-cfy

	m.topPanel.SetWidth(cx)

	mx, _ := m.gptStatus.GetSize()
	m.queryArea.SetWidth(cx - mx)

	_, qay := m.queryArea.GetSize()

	m.logViewer.SetWidth(cx)
	m.logViewer.SetHeight(cy - qay - 1)

	m.adviceList.SetWidth(cx)
	m.adviceList.SetHeight(cy - qay - 1)
}

func (m *Model) fetchAdvices(query string) tea.Cmd {
	m.logViewer.LogInfo("Received query: " + query)

	if isWaiting {
		m.logViewer.LogError("Already fetching response from OpenAI")
		return nil
	}

	isWaiting = true
	m.logViewer.Wait("Fetching response from OpenAI...")
	m.showLog = true

	logViewerCmd := m.logViewer.Init()

	fetchOpenAIResponse := func() tea.Msg {
		advices2, err := m.controller.GetSuggestion(m.gptStatus.CurrentModelID(), query, 5)

		m.logViewer.Done()
		isWaiting = false

		if err != nil {
			m.logViewer.LogError("Error: " + err.Error())
		} else {
			m.logViewer.LogSuccess("Response fetched")
			m.adviceList.SetAdvices(advices2)
			m.showLog = false
		}

		return nil
	}

	return tea.Batch(logViewerCmd, fetchOpenAIResponse)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.updateLayout(msg.Width, msg.Height)
		return m, nil
	case event.LogErrorMsg:
		m.logViewer.LogError(string(msg))
		return m, nil
	case event.QuerySentMsg:
		m.focusOnQuery = false
		return m, m.fetchAdvices(string(msg))
	case event.CommandConfirmedMsg:
		m.Result = string(msg)
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.focusOnQuery {
				m.queryArea.SetFocus(false)
				m.focusOnQuery = false
			}
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlO:
			return m, m.gptStatus.NextModel()
		case tea.KeyCtrlQ:
			if m.focusOnQuery {
				m.queryArea.SetFocus(false)
				m.focusOnQuery = false
			} else {
				m.queryArea.SetFocus(true)
				m.focusOnQuery = true
			}

			return m, nil
		}
	}

	if m.focusOnQuery {
		queryAreaCmd := m.queryArea.Update(msg)
		return m, queryAreaCmd
	} else {
		if !m.showLog {
			adviceListCmd := m.adviceList.Update(msg)
			return m, adviceListCmd
		}

		logViewerCmd := m.logViewer.Update(msg)
		return m, logViewerCmd
	}
}

func (m *Model) View() string {
	var (
		top  string
		main string
	)

	if m.focusOnQuery {
		top = m.topPanel.View("Query")
	} else if m.showLog {
		top = m.topPanel.View("Log")
	} else {
		top = m.topPanel.View("Advices")
	}

	if m.showLog {
		main = m.logViewer.View()
	} else {
		main = m.adviceList.View()
	}

	bottom := lipgloss.JoinHorizontal(lipgloss.Center, m.queryArea.View(), m.gptStatus.View())

	return appStyle.Render(
		mainContainerStyle.Render(lipgloss.JoinVertical(lipgloss.Top, top, main, bottom)),
	)
}
