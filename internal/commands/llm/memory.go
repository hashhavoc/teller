package llm

import (
	"github.com/sashabaranov/go-openai"
)

// Message represents a single message in the conversation
type Message struct {
	Role    string // Role can be "user", "assistant", or "system"
	Content string // The content of the message
}

// Memory structure to store conversation history
type Memory struct {
	conversationHistory []Message
	maxMessages         int // Limit to prevent context overflow
}

// NewMemory creates a new memory instance with specified max messages
func NewMemory(maxMessages int) *Memory {
	return &Memory{
		conversationHistory: make([]Message, 0),
		maxMessages:         maxMessages,
	}
}

// AddToMemory method to add a new message to memory
func (m *Memory) AddToMemory(role, content string) {
	m.conversationHistory = append(m.conversationHistory, Message{
		Role:    role,
		Content: content,
	})

	// Keep only the last N messages to prevent context overflow
	if len(m.conversationHistory) > m.maxMessages {
		m.conversationHistory = m.conversationHistory[len(m.conversationHistory)-m.maxMessages:]
	}
}

// GetMemory method to retrieve the conversation history as a slice of ChatCompletionMessage
func (m *Memory) GetMemory() []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage
	for _, msg := range m.conversationHistory {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return messages
}

// ClearMemory clears the conversation history
func (m *Memory) ClearMemory() {
	m.conversationHistory = make([]Message, 0)
}

// GetLastUserMessage returns the last user message, if any
func (m *Memory) GetLastUserMessage() string {
	for i := len(m.conversationHistory) - 1; i >= 0; i-- {
		if m.conversationHistory[i].Role == openai.ChatMessageRoleUser {
			return m.conversationHistory[i].Content
		}
	}
	return ""
} 