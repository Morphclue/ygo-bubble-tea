package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int64

const (
	Search Mode = iota
	Select
	View
)
const baseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php?fname="

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type model struct {
	textInput textinput.Model
	cardList  []Card
	mode      Mode
}

type Card struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	FrameType string `json:"frameType"`
	Desc      string `json:"desc"`
}

func getCards(cardName string, m model) {
	url := baseURL + cardName
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var body bytes.Buffer
	_, err = io.Copy(&body, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var data struct {
		Data []Card `json:"data"`
	}
	json.Unmarshal(body.Bytes(), &data)
	m.cardList = data.Data
}

func clearConsole() {
	var clearCommand string
	if runtime.GOOS == "windows" {
		clearCommand = "cls"
	} else {
		clearCommand = "clear"
	}
	cmd := exec.Command("cmd", "/c", clearCommand)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Dark Magician"
	ti.PromptStyle = focusedStyle
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		mode:      Search,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.mode = Select
			getCards(m.textInput.Value(), m)
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.mode {
	case Search:
		return fmt.Sprintf(
			"Enter a card name: \n%s\n",
			m.textInput.View(),
		) + helpStyle("\n enter: choose • q/ctrl+c: quit\n")
	case Select:
		return helpStyle("\n enter: choose • ↑/↓: select • q/ctrl+c: quit\n")
	}

	return "Unknown mode"
}

func main() {
	clearConsole()
	app := tea.NewProgram(initialModel())
	if _, err := app.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
