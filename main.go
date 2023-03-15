package main

import (
	"fmt"
	"os"

	"github.com/Morphclue/ygo-bubble-tea/api"
	"github.com/Morphclue/ygo-bubble-tea/entity"
	"github.com/Morphclue/ygo-bubble-tea/ui"
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

type model struct {
	textInput    textinput.Model
	cardList     list.Model
	infoTable    table.Model
	spinner      spinner.Model
	selectedCard entity.Card
	mode         Mode
	isLoading    bool
}

type getCardsMsg struct {
	cards []list.Item
}

func getCardsCmd(cardName string) tea.Cmd {
	return func() tea.Msg {
		cards, err := api.GetCards(cardName)
		if err != nil {
			return nil
		}
		cardListItems := make([]list.Item, len(cards))
		for i, card := range cards {
			cardListItems[i] = &ui.CardListItem{Card: card}
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
		cardList:  list.New([]list.Item{}, ui.ItemDelegate{}, 0, 0),
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
				m.selectedCard = m.cardList.SelectedItem().(*ui.CardListItem).Card
				m.infoTable = m.setInfoTable()
				m.infoTable.SetStyles(ui.TableStyle())
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
		m.cardList = list.New(msg.cards, ui.ItemDelegate{}, 20, 14)
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
