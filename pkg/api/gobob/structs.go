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

type TokenHoldersResponse struct {
	Items          []TokenHolderItem          `json:"items,omitempty"`
	NextPageParams TokenHoldersNextPageParams `json:"next_page_params,omitempty"`
}
type TokenHolderAddress struct {
	EnsDomainName      string   `json:"ens_domain_name,omitempty"`
	Hash               string   `json:"hash,omitempty"`
	ImplementationName string   `json:"implementation_name,omitempty"`
	IsContract         bool     `json:"is_contract,omitempty"`
	IsVerified         string   `json:"is_verified,omitempty"`
	Metadata           string   `json:"metadata,omitempty"`
	Name               string   `json:"name,omitempty"`
	PrivateTags        []string `json:"private_tags,omitempty"`
	PublicTags         []string `json:"public_tags,omitempty"`
	WatchlistNames     []string `json:"watchlist_names,omitempty"`
}
type TokenHolderToken struct {
	Address              string `json:"address,omitempty"`
	CirculatingMarketCap string `json:"circulating_market_cap,omitempty"`
	Decimals             string `json:"decimals,omitempty"`
	ExchangeRate         string `json:"exchange_rate,omitempty"`
	Holders              string `json:"holders,omitempty"`
	IconURL              string `json:"icon_url,omitempty"`
	Name                 string `json:"name,omitempty"`
	Symbol               string `json:"symbol,omitempty"`
	TotalSupply          string `json:"total_supply,omitempty"`
	Type                 string `json:"type,omitempty"`
	Volume24H            string `json:"volume_24h,omitempty"`
}

type TokenHolderItem struct {
	Address TokenHolderAddress `json:"address,omitempty"`
	Token   TokenHolderToken   `json:"token,omitempty"`
	TokenID string             `json:"token_id,omitempty"`
	Value   string             `json:"value,omitempty"`
}

type TokenHoldersNextPageParams struct {
	AddressHash string `json:"address_hash,omitempty" url:"address_hash,omitempty"`
	ItemsCount  int    `json:"items_count,omitempty" url:"items_count,omitempty"`
	Value       int64  `json:"value,omitempty" url:"value,omitempty"`
}
