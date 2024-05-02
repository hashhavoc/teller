package gobob

type TokenResponse struct {
	Items          []TokenItems   `json:"items"`
	NextPageParams NextPageParams `json:"next_page_params,omitempty"`
}

type TokenItems struct {
	CirculatingMarketCap string `json:"circulating_market_cap"`
	IconURL              string `json:"icon_url"`
	Name                 string `json:"name"`
	Decimals             string `json:"decimals"`
	Symbol               string `json:"symbol"`
	Address              string `json:"address"`
	Type                 string `json:"type"`
	Holders              string `json:"holders"`
	ExchangeRate         string `json:"exchange_rate"`
	TotalSupply          string `json:"total_supply"`
}

type NextPageParams struct {
	ContractAddressHash string `json:"contract_address_hash" url:"contract_address_hash"`
	FiatValue           string `json:"fiat_value" url:"fiat_value"`
	HolderCount         int    `json:"holder_count" url:"holder_count"`
	IsNameNull          bool   `json:"is_name_null" url:"is_name_null"`
	ItemsCount          int    `json:"items_count" url:"items_count"`
	MarketCap           string `json:"market_cap" url:"market_cap"`
	Name                string `json:"name" url:"name"`
}
