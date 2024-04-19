package transactions

type TransactionsResponse struct {
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
	Total   int           `json:"total"`
	Results []Transaction `json:"results"`
}

type Transaction struct {
	TxID              string      `json:"tx_id"`
	TxStatus          string      `json:"tx_status"`
	TxType            string      `json:"tx_type"`
	FeeRate           string      `json:"fee_rate"`
	SenderAddress     string      `json:"sender_address"`
	Sponsored         bool        `json:"sponsored"`
	PostConditionMode string      `json:"post_condition_mode"`
	BlockHash         string      `json:"block_hash"`
	BlockHeight       int         `json:"block_height"`
	BurnBlockTime     int         `json:"burn_block_time"`
	Canonical         bool        `json:"canonical"`
	IsUnanchored      bool        `json:"is_unanchored"`
	MicroblockHash    string      `json:"microblock_hash"`
	MicroblockSeq     int         `json:"microblock_sequence"`
	MicroblockCanon   bool        `json:"microblock_canonical"`
	TxIndex           int         `json:"tx_index"`
	CoinbasePayload   interface{} `json:"coinbase_payload"`
}
