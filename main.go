package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const baseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php?fname="

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type model struct {
	cursor int
}

func initialModel() model {
	// TODO: refactor model
	return model{
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	// TODO: initialize model
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := "Hello, world!"
	s += helpStyle("\nPress q or ctrl+c to quit.\n")
	return s
}

func main() {
	app := tea.NewProgram(initialModel())
	if _, err := app.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
