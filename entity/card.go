package entity

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
