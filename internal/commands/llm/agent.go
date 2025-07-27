package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashhavoc/teller/internal/commands/props"
	"github.com/sashabaranov/go-openai"
)

// StacksAgent handles Stacks blockchain queries based on LLM requests
type StacksAgent struct {
	props  *props.AppProps
	client *openai.Client
	memory *Memory
}

// NewStacksAgent creates a new Stacks agent
func NewStacksAgent(props *props.AppProps) (*StacksAgent, error) {
	apiKey := props.Config.OpenAI.APIKey
	if apiKey == "" {
		// Try to get from environment variable
		apiKey = os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OpenAI API key not found in config or OPENAI_API_KEY environment variable")
		}
	}

	// Create client configuration with custom base URL if specified
	config := openai.DefaultConfig(apiKey)
	if props.Config.OpenAI.BaseURL != "" && props.Config.OpenAI.BaseURL != "https://api.openai.com/v1" {
		config.BaseURL = props.Config.OpenAI.BaseURL
	}

	client := openai.NewClientWithConfig(config)
	memory := NewMemory(20) // Keep last 20 messages

	// Add system prompt to set context
	systemPrompt := `You are a helpful AI assistant specialized in Stacks blockchain data analysis. 
You have access to various Stacks APIs and can help users with:
- Account balances and transaction history
- Token information (fungible and non-fungible)
- Contract details and function calls
- Network status and block information
- DEX data from Alex protocol
- Names (BNS) system queries
- Ordinals and runes data

When a user asks for specific blockchain data, I will query the appropriate APIs and provide accurate, up-to-date information.

If you need to query specific data, let me know what type of information you need and any specific addresses, contract names, or other identifiers.

Available data types:
- Account balances: Provide a Stacks address
- Transactions: Provide a principal/address  
- Tokens: General token information or specific token details
- Contracts: Contract addresses and function calls
- DEX prices: Token pricing from Alex DEX
- Names: BNS name lookups
- Network status: Current blockchain state

Always be helpful and provide context about the data you return.`

	memory.AddToMemory(openai.ChatMessageRoleSystem, systemPrompt)

	return &StacksAgent{
		props:  props,
		client: client,
		memory: memory,
	}, nil
}

// ProcessQuery processes a user query and returns a response
func (sa *StacksAgent) ProcessQuery(userInput string) (string, error) {
	sa.memory.AddToMemory(openai.ChatMessageRoleUser, userInput)

	// Check if the query needs specific blockchain data
	blockchainData := sa.queryStacksData(userInput)

	// Add blockchain data to context if found
	var contextualInput string
	if blockchainData != "" {
		contextualInput = fmt.Sprintf("User query: %s\n\nBlockchain data retrieved:\n%s\n\nPlease analyze this data and provide a helpful response to the user.", userInput, blockchainData)
	} else {
		contextualInput = userInput
	}

	// Get LLM response
	response, err := sa.getLLMResponse(contextualInput)
	if err != nil {
		return "", err
	}

	sa.memory.AddToMemory(openai.ChatMessageRoleAssistant, response)
	return response, nil
}

// queryStacksData attempts to identify and fetch relevant blockchain data
func (sa *StacksAgent) queryStacksData(input string) string {
	input = strings.ToLower(input)

	// Check for address patterns (Stacks addresses start with 'SP' or 'ST')
	if strings.Contains(input, "sp") || strings.Contains(input, "st") {
		return sa.queryAddressData(input)
	}

	// Check for token-related queries
	if strings.Contains(input, "token") || strings.Contains(input, "balance") {
		return sa.queryTokenData(input)
	}

	// Check for transaction queries
	if strings.Contains(input, "transaction") || strings.Contains(input, "tx") {
		return sa.queryTransactionData(input)
	}

	// Check for network status queries
	if strings.Contains(input, "network") || strings.Contains(input, "status") || strings.Contains(input, "block") {
		return sa.queryNetworkStatus()
	}

	// Check for DEX/price queries
	if strings.Contains(input, "price") || strings.Contains(input, "dex") || strings.Contains(input, "alex") {
		return sa.queryDEXData()
	}

	return ""
}

// queryAddressData extracts address from input and fetches balance data
func (sa *StacksAgent) queryAddressData(input string) string {
	// Extract potential Stacks address (basic pattern matching)
	words := strings.Fields(strings.ToUpper(input))
	for _, word := range words {
		if strings.HasPrefix(word, "SP") || strings.HasPrefix(word, "ST") {
			if len(word) > 20 { // Basic length check for Stacks address
				balance, err := sa.props.HeroClient.GetAccountBalance(word, 0)
				if err == nil {
					data, _ := json.MarshalIndent(balance, "", "  ")
					return fmt.Sprintf("Account balance for %s:\n%s", word, string(data))
				}
			}
		}
	}
	return ""
}

// queryTokenData fetches general token information
func (sa *StacksAgent) queryTokenData(input string) string {
	tokens, err := sa.props.StxToolsClient.GetAllTokens()
	if err == nil && len(tokens) > 0 {
		// Return first few tokens as sample
		sample := tokens
		if len(sample) > 5 {
			sample = tokens[:5]
		}
		data, _ := json.MarshalIndent(sample, "", "  ")
		return fmt.Sprintf("Sample token data:\n%s", string(data))
	}
	return ""
}

// queryTransactionData would fetch transaction information
func (sa *StacksAgent) queryTransactionData(input string) string {
	// This could be enhanced to extract transaction IDs or addresses from input
	return ""
}

// queryNetworkStatus fetches current network status
func (sa *StacksAgent) queryNetworkStatus() string {
	// This would require implementing network status endpoints
	return "Network status functionality available - specific implementation would depend on available API endpoints"
}

// queryDEXData fetches DEX pricing information
func (sa *StacksAgent) queryDEXData() string {
	prices, err := sa.props.AlexClient.FetchLatestPrices()
	if err == nil {
		data, _ := json.MarshalIndent(prices, "", "  ")
		return fmt.Sprintf("Latest DEX prices from Alex:\n%s", string(data))
	}
	return ""
}

// getLLMResponse gets a response from the LLM
func (sa *StacksAgent) getLLMResponse(input string) (string, error) {
	// Create a temporary message list for this specific query
	messages := sa.memory.GetMemory()

	// Replace the last user message with our contextual input if we have blockchain data
	if len(messages) > 0 && messages[len(messages)-1].Role == openai.ChatMessageRoleUser {
		messages[len(messages)-1].Content = input
	}

	resp, err := sa.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       sa.props.Config.OpenAI.Model,
			Messages:    messages,
			MaxTokens:   500,
			Temperature: 0.7,
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

// ClearMemory clears the conversation history but keeps system prompt
func (sa *StacksAgent) ClearMemory() {
	sa.memory.ClearMemory()
	// Re-add system prompt
	systemPrompt := `You are a helpful AI assistant specialized in Stacks blockchain data analysis...` // Same as above
	sa.memory.AddToMemory(openai.ChatMessageRoleSystem, systemPrompt)
}
