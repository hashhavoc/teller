# Teller

Teller is a sophisticated Command Line Interface (CLI) tool designed to facilitate developers' interactions with STX contracts. It streamlines the process of exploring source code and leveraging various analytics tools, thereby significantly enhancing the developer experience within the STX ecosystem. With its intuitive design, Teller ensures that integrating STX contracts into projects is both straightforward and efficient.

## Installing

Homebrew is the recommended way to install Teller:

```sh
brew install hashhavoc/tap/teller
```

## Upgrade

```sh
brew upgrade hashhavoc/tap/teller
```

## Configuration

There is a configuration file that contains various configuration values. You can find an example at `config/config.yaml.example`.

```sh
teller conf init
```

This will create a new configuration file at `~/.teller.yaml` with the default values. You can then edit this file to your liking. There is not currently a way to specify the configuration file location. Not all of the endpoints are avaliable publicly, so you may need to specify your own endpoints.

The configuration file supports:

- **API Endpoints**: Configure custom endpoints for Hiro, Alex DEX, STXTools, BOB, and Ord
- **Wallet Addresses**: Pre-configure wallet addresses for quick access
- **OpenAI Integration**: Set up API keys, models, and custom base URLs for LLM functionality

See the [LLM Assistant](#llm-assistant) section for detailed OpenAI configuration options.

## Source

### Building

To use Teller, follow these simple steps:

1. Clone the repository:

    ```sh
    git clone https://github.com/hashhavoc/teller.git
    ```

2. Navigate to the project directory:

    ```sh
    cd teller
    ```

3. Build the executable:

    ```sh
    go build -o teller cmd/teller/main.go
    ```

4. Run Teller:

    ```sh
    ./teller
    ```

## Command Line Interface

Upon running Teller, you'll encounter the following command line interface:

```sh
➜  teller
NAME:
   teller - interact with the stx blockchain

USAGE:
   teller [global options] command [command options]

VERSION:
   v0.0.1

COMMANDS:
   config         Commands to manage the configuration file
   bob            Provides interactions with bob chain
   contracts      Provides interactions with contracts
   token          Provides interactions with tokens
   wallet         Provides interactions with wallets
   dex            Provides interactions with multiple dex
   transactions   Provides interactions with transactions
   ordinals, ord  Provides interactions with ordinals
   llm            Interactive AI assistant for Stacks blockchain queries
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Commands

Teller offers the following commands:

- **contracts**: Provides interactions with contracts.
- **token**: Provides interactions with tokens.
- **wallet**: Provides interactions with wallets.
- **dex**: Provides interactions with multiple decentralized exchanges.
- **transactions**: Provides interactions with transactions.
- **ordinals**: Provides interactions with ordinals on bitcoin.
- **llm**: Interactive AI assistant for Stacks blockchain queries with natural language interface.
- **help**: Shows a list of commands or help for one command.

## LLM Assistant

Teller includes an intelligent AI assistant that provides a conversational interface for exploring Stacks blockchain data. The assistant can automatically fetch and analyze blockchain information based on natural language queries.

### Features

- **Natural Language Queries**: Ask questions about blockchain data in plain English
- **Automatic Data Fetching**: Intelligent detection and retrieval of relevant blockchain information
- **Interactive Chat Interface**: Beautiful terminal UI similar to modern AI chat applications
- **Conversation Memory**: Maintains context across multiple interactions
- **Multi-API Integration**: Seamlessly queries Hiro, Alex DEX, STXTools, and other blockchain APIs

### Capabilities

The AI assistant can help with:

- **Account Information**: Check balances and transaction history for Stacks addresses
- **Token Analysis**: Explore fungible and non-fungible token data
- **DEX Data**: Get current prices and trading information from Alex protocol
- **Contract Exploration**: Analyze smart contracts and their functions
- **Network Status**: Check current blockchain state and network information
- **BNS Lookups**: Search and analyze Stacks naming system data
- **Ordinals & Runes**: Explore Bitcoin ordinals and runes data

### Usage

Start an interactive chat session:

```sh
teller llm chat
```

Use a specific model:

```sh
teller llm chat --model gpt-3.5-turbo
```

### Example Queries

```
> What's the balance for SP1ABC123...?
> Show me recent token data from STXTools
> What are the current DEX prices on Alex?
> Tell me about contract SP456DEF...
> How is the Stacks network performing?
> Look up the BNS name "example.btc"
```

### Configuration

Add OpenAI configuration to your `~/.teller.yaml`:

```yaml
openai:
  # Get your API key from https://platform.openai.com/api-keys
  # You can also set the OPENAI_API_KEY environment variable instead
  api_key: "your-openai-api-key"
  
  # Available models: gpt-4, gpt-4-turbo, gpt-3.5-turbo, etc.
  model: "gpt-4"
  
  # Base URL for OpenAI API - useful for self-hosted or alternative LLM endpoints
  # Default: https://api.openai.com/v1
  # Examples for other providers:
  # - OpenAI-compatible local models: http://localhost:1234/v1
  # - Azure OpenAI: https://your-resource.openai.azure.com/openai/deployments/your-deployment
  # - Anthropic Claude (via proxy): https://your-proxy/v1
  base_url: "https://api.openai.com/v1"
```

Alternatively, set the API key via environment variable:

```sh
export OPENAI_API_KEY="your-openai-api-key"
teller llm chat
```

## Support

If you encounter any issues or have suggestions for improvement, please feel free to open an issue on [GitHub](https://github.com/hashhavoc/teller/issues). Your feedback is highly appreciated!

## License

This project is licensed under the GNU General Public License v3.0 License - see the [LICENSE](LICENSE) file for details.
