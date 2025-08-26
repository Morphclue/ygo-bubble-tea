package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Morphclue/ygo-bubble-tea/entity"
)

var baseURL = "https://db.ygoprodeck.com/api/v7/cardinfo.php?fname="

func GetCards(cardName string) ([]entity.Card, error) {
	endpoint := baseURL + url.QueryEscape(cardName)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Data []entity.Card `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Data, nil
}
