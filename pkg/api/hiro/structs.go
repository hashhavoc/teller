package hiro

import "time"

type TransactionsResponse struct {
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
	Total   int           `json:"total"`
	Results []Transaction `json:"results"`
}

type Transaction struct {
	Tx          Tx     `json:"tx,omitempty"`
	StxSent     string `json:"stx_sent,omitempty"`
	StxReceived string `json:"stx_received,omitempty"`
	Events      Events `json:"events,omitempty"`
}

type TxResult struct {
	Hex  string `json:"hex,omitempty"`
	Repr string `json:"repr,omitempty"`
}

type TokenTransfer struct {
	RecipientAddress string `json:"recipient_address,omitempty"`
	Amount           string `json:"amount,omitempty"`
	Memo             string `json:"memo,omitempty"`
}

type ContractCall struct {
	ContractId        string         `json:"contract_id,omitempty"`
	FunctionName      string         `json:"function_name,omitempty"`
	FunctionSignature string         `json:"function_signature,omitempty"`
	FunctionArgs      []ClarityValue `json:"function_args,omitempty"`
}

type ClarityValue struct {
	Hex  string `json:"hex,omitempty"`
	Repr string `json:"repr,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type Tx struct {
	TxID                     string        `json:"tx_id,omitempty"`
	Nonce                    int           `json:"nonce,omitempty"`
	FeeRate                  string        `json:"fee_rate,omitempty"`
	SenderAddress            string        `json:"sender_address,omitempty"`
	Sponsored                bool          `json:"sponsored,omitempty"`
	PostConditionMode        string        `json:"post_condition_mode,omitempty"`
	PostConditions           []any         `json:"post_conditions,omitempty"`
	AnchorMode               string        `json:"anchor_mode,omitempty"`
	IsUnanchored             bool          `json:"is_unanchored,omitempty"`
	BlockHash                string        `json:"block_hash,omitempty"`
	ParentBlockHash          string        `json:"parent_block_hash,omitempty"`
	BlockHeight              int           `json:"block_height,omitempty"`
	BlockTime                int           `json:"block_time,omitempty"`
	BlockTimeIso             time.Time     `json:"block_time_iso,omitempty"`
	BurnBlockTime            int           `json:"burn_block_time,omitempty"`
	BurnBlockTimeIso         time.Time     `json:"burn_block_time_iso,omitempty"`
	ParentBurnBlockTime      int           `json:"parent_burn_block_time,omitempty"`
	ParentBurnBlockTimeIso   time.Time     `json:"parent_burn_block_time_iso,omitempty"`
	Canonical                bool          `json:"canonical,omitempty"`
	ContractCall             ContractCall  `json:"contract_call,omitempty"`
	TxIndex                  int           `json:"tx_index,omitempty"`
	TxStatus                 string        `json:"tx_status,omitempty"`
	TxResult                 TxResult      `json:"tx_result,omitempty"`
	MicroblockHash           string        `json:"microblock_hash,omitempty"`
	MicroblockSequence       int64         `json:"microblock_sequence,omitempty"`
	MicroblockCanonical      bool          `json:"microblock_canonical,omitempty"`
	EventCount               int           `json:"event_count,omitempty"`
	Events                   []Event       `json:"events,omitempty"`
	ExecutionCostReadCount   int           `json:"execution_cost_read_count,omitempty"`
	ExecutionCostReadLength  int           `json:"execution_cost_read_length,omitempty"`
	ExecutionCostRuntime     int           `json:"execution_cost_runtime,omitempty"`
	ExecutionCostWriteCount  int           `json:"execution_cost_write_count,omitempty"`
	ExecutionCostWriteLength int           `json:"execution_cost_write_length,omitempty"`
	TxType                   string        `json:"tx_type,omitempty"`
	TokenTransfer            TokenTransfer `json:"token_transfer,omitempty"`
}

type Stx struct {
	Transfer int `json:"transfer,omitempty"`
	Mint     int `json:"mint,omitempty"`
	Burn     int `json:"burn,omitempty"`
}

type Ft struct {
	Transfer int `json:"transfer,omitempty"`
	Mint     int `json:"mint,omitempty"`
	Burn     int `json:"burn,omitempty"`
}

type Nft struct {
	Transfer int `json:"transfer,omitempty"`
	Mint     int `json:"mint,omitempty"`
	Burn     int `json:"burn,omitempty"`
}

type Events struct {
	Stx Stx `json:"stx,omitempty"`
	Ft  Ft  `json:"ft,omitempty"`
	Nft Nft `json:"nft,omitempty"`
}

type Event struct {
	EventIndex  int         `json:"event_index,omitempty"`
	EventType   string      `json:"event_type,omitempty"`
	TxId        string      `json:"tx_id,omitempty"`
	Asset       Asset       `json:"asset,omitempty"`
	ContractLog ContractLog `json:"contract_log,omitempty"`
}

type Asset struct {
	AssetEventType string       `json:"asset_event_type,omitempty"`
	Sender         string       `json:"sender,omitempty"`
	Recipient      string       `json:"recipient,omitempty"`
	Amount         int          `json:"amount,omitempty"`
	AssetId        string       `json:"asset_id,omitempty"`
	Value          ClarityValue `json:"value,omitempty"`
}

type ContractLog struct {
	ContractId string       `json:"contract_id,omitempty"`
	Topic      string       `json:"topic,omitempty"`
	Value      ClarityValue `json:"value,omitempty"`
}

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
