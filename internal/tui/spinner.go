package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel manages spinning progress state
type SpinnerModel struct {
	Spinner  spinner.Model
	Message  string
	Progress string
	Done     bool
	Err      error
}

// NewSpinner creates an initialized SpinnerModel
func NewSpinner(initialMessage string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(ColorPrimary)
	return SpinnerModel{
		Spinner: s,
		Message: initialMessage,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case string:
		m.Message = msg
		return m, nil
	default:
		return m, nil
	}
}

func (m SpinnerModel) View() string {
	if m.Done {
		return fmt.Sprintf(" %s %s\n", SuccessBadge.Render("✓"), m.Message)
	}
	return fmt.Sprintf(" %s %s %s\n", m.Spinner.View(), m.Message, SubtitleStyle.Render(m.Progress))
}
