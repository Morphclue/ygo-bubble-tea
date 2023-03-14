package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Morphclue/ygo-bubble-tea/ui"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int64

const (
	Search Mode = iota
	Select
	View
)

const baseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php?fname="

type Card struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	FrameType string `json:"frameType"`
	Desc      string `json:"desc"`
	CardSets  []struct {
		SetCode       string `json:"set_code"`
		SetRarityCode string `json:"set_rarity_code"`
		SetPrice      string `json:"set_price"`
	} `json:"card_sets"`
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                         { return 1 }
func (d itemDelegate) Spacing() int                        { return 0 }
func (d itemDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*cardListItem)
	if !ok || i == nil {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.card.Name)

	fn := ui.ItemStyle.Render
	if index == m.Index() {
		fn = func(strings ...string) string {
			args := make([]interface{}, len(strings)-1)
			for i, arg := range strings[1:] {
				args[i] = arg
			}
			return ui.SelectedItemStyle.Render("> " + fmt.Sprintf(strings[0], args...))
		}
	}

	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		return
	}
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
	spinner      spinner.Model
	selectedCard Card
	mode         Mode
	isLoading    bool
}

type getCardsMsg struct {
	cards []list.Item
}

func getCardsCmd(cardName string) tea.Cmd {
	return func() tea.Msg {
		url := baseURL + cardName
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)
		var body bytes.Buffer
		_, err = io.Copy(&body, resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		var data struct {
			Data []Card `json:"data"`
		}
		err = json.Unmarshal(body.Bytes(), &data)
		if err != nil {
			return nil
		}
		cardListItems := make([]list.Item, len(data.Data))
		for i, card := range data.Data {
			cardListItems[i] = &cardListItem{card: card}
		}
		return getCardsMsg{cards: cardListItems}
	}
}

func (m model) setInfoTable() table.Model {
	columns := []table.Column{
		{Title: "Code", Width: 10},
		{Title: "Rarity", Width: 10},
		{Title: "Price", Width: 10},
	}

	var rows []table.Row
	for _, cardSet := range m.selectedCard.CardSets {
		rows = append(rows, table.Row{
			cardSet.SetCode,
			cardSet.SetRarityCode,
			cardSet.SetPrice,
		})
	}

	generatedTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return generatedTable
}
func (m model) styleTable() table.Styles {
	return ui.TableStyles()
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Dark Magician"
	ti.PromptStyle = ui.FocusedStyle
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	s := ui.Spinner()

	return model{
		textInput: ti,
		mode:      Search,
		spinner:   s,
		cardList:  list.New([]list.Item{}, itemDelegate{}, 0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			switch m.mode {
			case Search:
				m.mode = Select
				m.isLoading = true
				return m, tea.Batch(m.spinner.Tick, getCardsCmd(m.textInput.Value()))
			case Select:
				m.selectedCard = m.cardList.SelectedItem().(*cardListItem).card
				m.infoTable = m.setInfoTable()
				m.infoTable.SetStyles(m.styleTable())
				m.mode = View
			}
		case "b":
			switch m.mode {
			case Select:
				m.mode = Search
			case View:
				m.mode = Select
			}
		}
	case getCardsMsg:
		m.cardList = list.New(msg.cards, itemDelegate{}, 20, 14)
		m.isLoading = false
	}

	switch m.mode {
	case Search:
		m.textInput, cmd = m.textInput.Update(msg)
	case Select:
		m.cardList, cmd = m.cardList.Update(msg)
	case View:
		m.infoTable, cmd = m.infoTable.Update(msg)
	}

	var sCmd tea.Cmd
	m.spinner, sCmd = m.spinner.Update(msg)
	cmds = append(cmds, sCmd, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.isLoading {
		return fmt.Sprintf("Loading... %s", m.spinner.View())
	}
	switch m.mode {
	case Search:
		return fmt.Sprintf(
			"Enter a card name: \n%s\n",
			m.textInput.View(),
		) + ui.HelpStyle("\n enter: choose • q/ctrl+c: quit\n")
	case Select:
		m.cardList.Title = "Select a card"
		return fmt.Sprintf(
			m.cardList.View(),
		) + ui.HelpStyle("\n enter: choose • b: back\n")
	case View:
		return fmt.Sprintf(
			m.infoTable.View(),
		) + ui.HelpStyle("\n b: back • q/ctrl+c: quit\n")
	}

	return "Unknown mode"
}

func main() {
	app := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := app.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
