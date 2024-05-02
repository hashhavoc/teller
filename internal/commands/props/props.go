package props

import (
	"github.com/phuslu/log"

	"github.com/hashhavoc/teller/internal/config"
	"github.com/hashhavoc/teller/pkg/api/alex"
	"github.com/hashhavoc/teller/pkg/api/gobob"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/api/ord"
	"github.com/hashhavoc/teller/pkg/api/stxtools"
)

type AppProps struct {
	AlexClient     *alex.APIClient
	HeroClient     *hiro.APIClient
	StxToolsClient *stxtools.APIClient
	OrdClient      *ord.APIClient
	BobClient      *gobob.APIClient
	Config         *config.Config
	Logger         log.Logger
}

func NewAppProps() *AppProps {
	return &AppProps{}
}
