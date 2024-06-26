package wallet

type TokenBalanceInfo struct {
	STXBalance             int64
	FungibleTokenBalances  map[string]int64
	NonFungibleTokenCounts map[string]int64
}

type Address struct {
	PrivateKey string
	Address    string
}
