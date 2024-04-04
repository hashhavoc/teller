package alex

type CurrencyPair struct {
	TickerID       string  `json:"ticker_id"`
	PoolID         string  `json:"pool_id"`
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Base           string  `json:"base"`
	Target         string  `json:"target"`
	LastPrice      float64 `json:"last_price"`
	BaseVolume     float64 `json:"base_volume"`
	TargetVolume   float64 `json:"target_volume"`
	LiquidityInUSD float64 `json:"liquidity_in_usd"`
}

type TokenPriceResponse struct {
	Data struct {
		LaplaceCurrentTokenPrice []struct {
			AvgPriceUSD float64 `json:"avg_price_usd"`
			Token       string  `json:"token"`
		} `json:"laplace_current_token_price"`
	} `json:"data"`
}
