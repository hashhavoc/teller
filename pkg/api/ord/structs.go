package ord

// Define the structs according to the JSON structure
type Entry struct {
	ID      string  `json:"id"`
	Details Details `json:"details"`
}

type Details struct {
	Block        int64   `json:"block"`
	Burned       int64   `json:"burned"`
	Divisibility int64   `json:"divisibility"`
	Etching      string  `json:"etching"`
	Mints        int64   `json:"mints"`
	Number       int64   `json:"number"`
	Premine      float64 `json:"premine"`
	SpacedRune   string  `json:"spaced_rune"`
	Symbol       string  `json:"symbol"`
	Terms        Terms   `json:"terms"`
	Timestamp    int64   `json:"timestamp"`
	Turbo        bool    `json:"turbo"`
}

type Terms struct {
	Amount float64       `json:"amount"`
	Cap    uint64        `json:"cap"`
	Height []int         `json:"height"`
	Offset []interface{} `json:"offset"` // Use interface{} for mixed types (null, int)
}

type Response struct {
	Entries [][]interface{} `json:"entries"` // Use interface{} because the first element is string and second is Details
	More    bool            `json:"more"`
	Prev    interface{}     `json:"prev"` // Use interface{} because it can be null or other types
	Next    interface{}     `json:"next"`
}
