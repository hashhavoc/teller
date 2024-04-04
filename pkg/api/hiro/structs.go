package hiro

type TokenResult struct {
	Name              string `json:"name"`
	Symbol            string `json:"symbol"`
	Decimals          int    `json:"decimals"`
	TotalSupply       string `json:"total_supply"`
	TokenURI          string `json:"token_uri"`
	Description       string `json:"description"`
	ImageURI          string `json:"image_uri"`
	ImageCanonicalURI string `json:"image_canonical_uri"`
	TxID              string `json:"tx_id"`
	SenderAddress     string `json:"sender_address"`
	ContractPrincipal string `json:"contract_principal"`
}

type Response struct {
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
	Total   int           `json:"total"`
	Results []TokenResult `json:"results"`
}

type ContractSourceResponse struct {
	Source        string `json:"source"`
	PublishHeight int    `json:"publish_height"`
	Proof         string `json:"proof"`
}

type NFTHoldingResponse struct {
	Limit   int                         `json:"limit"`
	Offset  int                         `json:"offset"`
	Total   int                         `json:"total"`
	Results []NFTHoldingResponseResults `json:"results"`
}

type NFTHoldingResponseResultsValue struct {
	Hex  string `json:"hex"`
	Repr string `json:"repr"`
}

type NFTHoldingResponseResults struct {
	AssetIdentifier string                         `json:"asset_identifier"`
	Value           NFTHoldingResponseResultsValue `json:"value"`
	BlockHeight     int                            `json:"block_height"`
	TxID            string                         `json:"tx_id"`
}

type Balance struct {
	Balance       string `json:"balance"`
	TotalSent     string `json:"total_sent"`
	TotalReceived string `json:"total_received"`
}
type FungibleTokens map[string]Balance

type NonFungibleTokens map[string]struct {
	Count         string `json:"count"`
	TotalSent     string `json:"total_sent"`
	TotalReceived string `json:"total_received"`
}

type BalanceResponse struct {
	Stx               Balance           `json:"stx"`
	FungibleTokens    FungibleTokens    `json:"fungible_tokens"`
	NonFungibleTokens NonFungibleTokens `json:"non_fungible_tokens"`
}

type ReadOnlyPayload struct {
	Sender    string   `json:"sender"`
	Arguments []string `json:"arguments"`
}

type ReadOnlyResponse struct {
	Okay   bool   `json:"okay"`
	Result string `json:"result"`
}

type ContractDetailsResponse struct {
	TxID           string `json:"tx_id"`
	Canonical      bool   `json:"canonical"`
	ContractID     string `json:"contract_id"`
	BlockHeight    int    `json:"block_height"`
	ClarityVersion int    `json:"clarity_version"`
	SourceCode     string `json:"source_code"`
	ABI            string `json:"abi"`
}

type ContractHoldersResponse map[string]string
