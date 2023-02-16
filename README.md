# ygo-bubble-tea

This is a command-line application that allows you to search for Yu-Gi-Oh! cards using the [Ygoprodeck API](https://ygoprodeck.com/api-guide/).
It is written in Go and uses the [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea).
You can search for a card name, select a card from a list of cards and view detailed information about it.

### How to use
![Demo Usage](./assets/demo.gif)

1. Type in a card name in the text input field (e.g. "Dark Magician") and press enter.
2. The search results will be displayed below the input field. Use the up and down arrow keys to navigate through the
   results. 
3. Press enter to select a card.
4. The selected card's information will be displayed as a table.
