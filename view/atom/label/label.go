package label

import "github.com/charmbracelet/lipgloss"

var defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFCC66"))

type Model struct {
	Text  string
	Style lipgloss.Style
}

func New(text string) Model {
	return Model{
		Text:  text,
		Style: defaultStyle.Copy(),
	}
}

func (l *Model) View() string {
	if len(l.Text) > 0 {
		return l.Style.Render(l.Text)
	}

	return ""
}

func (l *Model) GetSize() (int, int) {
	x, y := 0, 0

	if len(l.Text) > 0 {
		lx, ly := l.Style.GetFrameSize()
		x += lx
		y += ly

		if l.Style.GetWidth() == 0 {
			x += len(l.Text)
		} else {
			x += l.Style.GetWidth()
		}

		if l.Style.GetHeight() == 0 {
			y++
		} else {
			y += l.Style.GetHeight()
		}
	}

	return x, y
}
