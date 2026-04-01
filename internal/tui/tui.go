package tui

import (
	"fmt"
	"strings"

	"github.com/febriliankr/lazymake/internal/parser"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2"))
	nameStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("4"))
	descStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	fileStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
)

type Model struct {
	targets   []parser.Target
	filtered  []parser.Target
	textInput textinput.Model
	cursor    int
	Selected  *parser.Target
	quitting  bool
	height    int
	multiFile bool
}

func New(targets []parser.Target) Model {
	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.Focus()
	ti.CharLimit = 100

	files := make(map[string]bool)
	for _, t := range targets {
		files[t.File] = true
	}

	return Model{
		targets:   targets,
		filtered:  targets,
		textInput: ti,
		multiFile: len(files) > 1,
		height:    20,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if len(m.filtered) > 0 {
				t := m.filtered[m.cursor]
				m.Selected = &t
				return m, tea.Quit
			}

		case "up", "ctrl+k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "ctrl+j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}

		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			m.filtered = filterTargets(m.targets, m.textInput.Value())
			m.cursor = 0
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting && m.Selected == nil {
		return ""
	}
	if m.Selected != nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("  lazymake") + "\n\n")
	b.WriteString("  " + m.textInput.View() + "\n\n")

	// Calculate visible items
	maxVisible := m.height - 7
	if maxVisible < 3 {
		maxVisible = 3
	}

	start := 0
	if m.cursor >= maxVisible {
		start = m.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	if len(m.filtered) == 0 {
		b.WriteString(descStyle.Render("  No targets found") + "\n")
	}

	for i := start; i < end; i++ {
		t := m.filtered[i]
		cursor := "  "
		name := nameStyle.Render(t.Name)
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
			name = selectedStyle.Render(t.Name)
		}

		line := fmt.Sprintf("%s%-30s", cursor, name)
		if t.Description != "" {
			line += " " + descStyle.Render(t.Description)
		}
		if m.multiFile {
			line += " " + fileStyle.Render(t.File)
		}
		b.WriteString(line + "\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  ↑↓ navigate • enter select • esc quit"))
	return b.String()
}

func filterTargets(targets []parser.Target, query string) []parser.Target {
	if query == "" {
		return targets
	}

	names := make([]string, len(targets))
	for i, t := range targets {
		names[i] = t.Name + " " + t.Description
	}

	matches := fuzzy.Find(query, names)
	result := make([]parser.Target, len(matches))
	for i, m := range matches {
		result[i] = targets[m.Index]
	}
	return result
}

func Run(targets []parser.Target) (*parser.Target, error) {
	m := New(targets)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(Model)
	return result.Selected, nil
}
