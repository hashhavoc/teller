module github.com/hashhavoc/teller

go 1.22.0

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/charmbracelet/bubbles v0.18.0
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/charmbracelet/lipgloss v0.10.0
	github.com/hashhavoccat/stacks-go v0.0.0-20240404042855-623fe7696f55
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/phuslu/log v1.0.89
	github.com/pkg/errors v0.9.1
	github.com/urfave/cli/v2 v2.27.1
	golang.design/x/clipboard v0.7.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/btcsuite/btcutil v1.0.2 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/go-openapi/errors v0.21.0 // indirect
	github.com/go-openapi/strfmt v0.22.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/okx/go-wallet-sdk/crypto v0.0.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 // indirect
	golang.org/x/image v0.10.0 // indirect
	golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace github.com/hashhavoc/teller/internal/commands/transactions => github.com/whoabuddy/teller/internal/commands/transactions feat/download-transactions
