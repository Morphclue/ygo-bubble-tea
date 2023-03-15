package ui

import (
	"fmt"
	"io"

	"github.com/Morphclue/ygo-bubble-tea/entity"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ItemDelegate struct{}

func (d ItemDelegate) Height() int                         { return 1 }
func (d ItemDelegate) Spacing() int                        { return 0 }
func (d ItemDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*CardListItem)
	if !ok || i == nil {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Card.Name)

	fn := ItemStyle.Render
	if index == m.Index() {
		fn = func(strings ...string) string {
			args := make([]interface{}, len(strings)-1)
			for i, arg := range strings[1:] {
				args[i] = arg
			}
			return SelectedItemStyle.Render("> " + fmt.Sprintf(strings[0], args...))
		}
	}

	_, err := fmt.Fprint(w, fn(str))
	if err != nil {
		return
	}
}

type CardListItem struct {
	Card entity.Card
}

func (c CardListItem) FilterValue() string {
	return c.Card.Name
}
