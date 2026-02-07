package main

import (
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	cwd string
	entries []os.DirEntry
	cursor int
	height int
}

func loadDir(path string) []os.DirEntry {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	return entries
}

func initialModel() model {
	cwd, _ := os.Getwd()
	return model{
		cwd: cwd,
		entries: loadDir(cwd),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var s string

	header := lipgloss.NewStyle().Bold(true).Render(m.cwd)
	s += header + "\n\n"

	for i, e := range m.entries {
		name := e.Name()
		if e.IsDir() {
			name += "/"
		}

		if i == m.cursor {
			name = lipgloss.NewStyle().Reverse(true).Render(name)
		}

		s += name + "\n"
	}
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
	
	case tea.KeyMsg:
		switch msg.String() {
		
		case "q", "ctrl+c":
			return m, tea.Quit
		
		case "j", "down":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter", "1", "right":
			if len(m.entries) == 0 {
				break
			}
			e := m.entries[m.cursor]
			if e.IsDir() {
				newPath := filepath.Join(m.cwd, e.Name())
				m.cwd = newPath
				m.entries = loadDir(newPath)
				m.cursor = 0
			}

		case "h", "backspace", "left":
			parent := filepath.Dir(m.cwd)
			m.cwd = parent
			m.entries = loadDir(parent)
			m.cursor = 0
		}
	}
	
	return m, nil

}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		os.Exit(1)
	}
}


