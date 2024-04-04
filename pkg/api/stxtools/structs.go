package stxtools

import "time"

type Token struct {
	ContractID        string  `json:"contract_id"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	Decimals          int     `json:"decimals"`
	CirculatingSupply string  `json:"circulating_supply"`
	TotalSupply       string  `json:"total_supply"`
	ImageURL          string  `json:"image_url"`
	Enabled           bool    `json:"enabled"`
	WrappedToken      string  `json:"wrapped_token"`
	Metrics           Metrics `json:"metrics"`
}

type Metrics struct {
	ContractID     string  `json:"contract_id"`
	HolderCount    int     `json:"holder_count"`
	SwapCount      int     `json:"swap_count"`
	TransferCount  int     `json:"transfer_count"`
	PriceUSD       float64 `json:"price_usd"`
	PriceChange1D  float64 `json:"price_change_1d"`
	PriceChange7D  float64 `json:"price_change_7d"`
	PriceChange30D float64 `json:"price_change_30d"`
	LiquidityUSD   float64 `json:"liquidity_usd"`
}

type HoldersResponse struct {
	Data HoldersData `json:"data"`
	Page Page        `json:"page"`
}

type HoldersData struct {
	TokenInfo  TokenInfo   `json:"token_info"`
	TopHolders []TopHolder `json:"top_holders"`
}

type TokenInfo struct {
	Decimals    int    `json:"decimals"`
	TotalSupply string `json:"total_supply"`
}

type TopHolder struct {
	WalletAddress string `json:"wallet_address"`
	TokenBalance  string `json:"token_balance"`
	TotalSent     string `json:"total_sent"`
	TotalReceived string `json:"total_received"`
	Wallet        Wallet `json:"wallet"`
	Rank          int    `json:"rank"`
}

type Wallet struct {
	Address    string      `json:"address"`
	StxBalance string      `json:"stx_balance"`
	WalletTags []WalletTag `json:"wallet_tags"`
}

type WalletTag struct {
	Tag string `json:"tag"`
}

type SwapsResponse struct {
	Data []SwapsData `json:"data"`
	Page Page        `json:"page"`
}

type SwapsData struct {
	TxID          string    `json:"tx_id"`
	PoolID        string    `json:"pool_id"`
	SenderAddress string    `json:"sender_address"`
	TokenXAmount  string    `json:"token_x_amount"`
	TokenYAmount  string    `json:"token_y_amount"`
	BurnBlockTime time.Time `json:"burn_block_time"`
	TokenX        SwapToken `json:"token_x"`
	TokenY        SwapToken `json:"token_y"`
}

type SwapToken struct {
	ContractID string `json:"contract_id"`
	Decimals   int    `json:"decimals"`
	ImageURL   string `json:"image_url"`
	Symbol     string `json:"symbol"`
}

type Page struct {
	Size          int    `json:"size"`
	TotalPages    int    `json:"totalPages"`
	TotalElements int    `json:"totalElements"`
	PageNumber    int    `json:"pageNumber"`
	SortDirection string `json:"sortDirection"`
}

type TransfersResponse struct {
	Data []Transaction `json:"data"`
	Page Page          `json:"page"`
}

type Transaction struct {
	TxID             string `json:"tx_id"`
	SenderAddress    string `json:"sender_address"`
	Amount           string `json:"amount"`
	RecipientAddress string `json:"recipient_address"`
	ContractID       string `json:"contract_id"`
	BurnBlockTime    string `json:"burn_block_time"`
	BlockHeight      int    `json:"block_height"`
	Token            Token  `json:"token"`
}

type TransfersToken struct {
	Symbol      string `json:"symbol"`
	ImageURL    string `json:"image_url"`
	Decimals    int    `json:"decimals"`
	TotalSupply string `json:"total_supply"`
}
