package llm

import (
	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/urfave/cli/v2"
)

func CreateLLMCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "llm",
		Usage: "Interactive AI assistant for Stacks blockchain queries",
		Subcommands: []*cli.Command{
			createChatCommand(props),
		},
	}
}

func createChatCommand(props *props.AppProps) *cli.Command {
	return &cli.Command{
		Name:  "chat",
		Usage: "Start interactive chat with AI assistant for Stacks blockchain",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "OpenAI model to use (default from config)",
			},
		},
		Action: func(c *cli.Context) error {
			model := c.String("model")
			if model == "" {
				model = props.Config.OpenAI.Model
			}
			return startChatTUI(props, model)
		},
	}
} 