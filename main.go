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
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
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
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type Card struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	FrameType string `json:"frameType"`
	Desc      string `json:"desc"`
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*cardListItem)
	if !ok || i == nil {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.card.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type cardListItem struct {
	card Card
}

func (c cardListItem) FilterValue() string {
	return c.card.Name
}

type model struct {
	textInput    textinput.Model
	cardList     list.Model
	infoTable    table.Model
	selectedCard Card
	mode         Mode
}

func getCards(cardName string) []list.Item {
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
	cardListItems := make([]list.Item, len(data.Data))
	for i, card := range data.Data {
		cardListItems[i] = &cardListItem{card: card}
	}
	return cardListItems
}

func (m model) setInfoTable() table.Model {
	columns := []table.Column{
		{Title: "Type", Width: 10},
		{Title: "Value", Width: 30},
	}

	rows := []table.Row{
		{"ID", strconv.Itoa(m.selectedCard.Id)},
		{"Name", m.selectedCard.Name},
		{"Type", m.selectedCard.Type},
		{"Frame Type", m.selectedCard.FrameType},
		{"Description", m.selectedCard.Desc},
	}

	generatedTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	return generatedTable
}

func (m model) styleTable() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return s
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
			switch m.mode {
			case Search:
				m.mode = Select
				items := getCards(m.textInput.Value())
				m.cardList = list.New(items, itemDelegate{}, 20, 14)
			case Select:
				m.selectedCard = m.cardList.SelectedItem().(*cardListItem).card
				m.infoTable = m.setInfoTable()
				m.infoTable.SetStyles(m.styleTable())
				m.mode = View
			}
		}
	}

	switch m.mode {
	case Search:
		m.textInput, cmd = m.textInput.Update(msg)
	case Select:
		m.cardList, cmd = m.cardList.Update(msg)
	case View:
		m.infoTable, cmd = m.infoTable.Update(msg)
	}

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
		m.cardList.Title = "Select a card"
		return fmt.Sprintf(
			m.cardList.View(),
		) + helpStyle("\n enter: choose • ↑/↓: select • q/ctrl+c: quit\n")
	case View:
		return fmt.Sprintf(
			m.infoTable.View(),
		)
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
