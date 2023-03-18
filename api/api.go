package api

import (
	"encoding/json"
	"net/http"

	"github.com/Morphclue/ygo-bubble-tea/entity"
)

var baseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php?fname="

func GetCards(cardName string) ([]entity.Card, error) {
	url := baseURL + cardName
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body = resp.Body
	var data struct {
		Data []entity.Card `json:"data"`
	}
	err = json.NewDecoder(body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
