package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Morphclue/ygo-bubble-tea/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetCards(t *testing.T) {
	mockCards := []entity.Card{
		{
			Id:   31305911,
			Name: "Marshmallon",
			Type: "Effect Monster",
			Desc: "Cannot be destroyed by battle. After damage calculation, if this card was attacked, " +
				"and was face-down at the start of the Damage Step: The attacking player takes 1000 damage.",
			CardSets: []struct {
				SetCode       string `json:"set_code"`
				SetRarityCode string `json:"set_rarity_code"`
				SetPrice      string `json:"set_price"`
			}{
				{SetCode: "DPYG-EN015", SetRarityCode: "R", SetPrice: "1.79"},
				{SetCode: "LDK2-ENY20", SetRarityCode: "C", SetPrice: "1.51"},
				{SetCode: "SDLS-EN013", SetRarityCode: "C", SetPrice: "1.29"},
				{SetCode: "PP01-EN003", SetRarityCode: "ScR", SetPrice: "5.28"},
				{SetCode: "YS18-EN017", SetRarityCode: "C", SetPrice: "1.26"},
				{SetCode: "YS17-EN015", SetRarityCode: "C", SetPrice: "1.55"},
				{SetCode: "YGLD-ENC22", SetRarityCode: "C", SetPrice: "1.34"},
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(struct {
			Data []entity.Card `json:"data"`
		}{
			Data: mockCards,
		})
		if err != nil {
			return
		}
	}))
	defer ts.Close()
	baseURL = ts.URL + "/"

	cards, err := GetCards("Marshmallon")
	assert.NoError(t, err)
	assert.Equal(t, mockCards, cards)
}

func TestGetCards_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error"))
	}))
	defer ts.Close()
	baseURL = ts.URL + "/"

	_, err := GetCards("NonExistentCard")
	assert.Error(t, err)
}
