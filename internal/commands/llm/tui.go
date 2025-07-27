package llm

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashhavoc/teller/internal/commands/props"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Padding(0, 1)

	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			PaddingLeft(1)

	assistantStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true).
				PaddingLeft(1)

	messageStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			PaddingLeft(1)
)

type chatMessage struct {
	role    string
	content string
}

type model struct {
	agent        *StacksAgent
	viewport     viewport.Model
	textarea     textarea.Model
	messages     []chatMessage
	err          error
	ready        bool
	processing   bool
}

type processedMsg struct {
	response string
	err      error
}

func startChatTUI(props *props.AppProps, modelName string) error {
	agent, err := NewStacksAgent(props)
	if err != nil {
		return fmt.Errorf("failed to initialize Stacks agent: %w", err)
	}

	// Override model if specified
	if modelName != "" {
		props.Config.OpenAI.Model = modelName
	}

	ta := textarea.New()
	ta.Placeholder = "Ask me anything about Stacks blockchain..."
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = 1000
	ta.SetWidth(80)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.SetContent("")

	m := model{
		agent:    agent,
		textarea: ta,
		viewport: vp,
		messages: []chatMessage{},
		ready:    false,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-6)
			m.textarea.SetWidth(msg.Width - 4)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 6
			m.textarea.SetWidth(msg.Width - 4)
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if !m.processing {
				userInput := strings.TrimSpace(m.textarea.Value())
				if userInput != "" {
					m.messages = append(m.messages, chatMessage{role: "user", content: userInput})
					m.textarea.Reset()
					m.processing = true
					m.updateViewport()
					return m, m.processQuery(userInput)
				}
			}
		}

	case processedMsg:
		m.processing = false
		if msg.err != nil {
			m.err = msg.err
			m.messages = append(m.messages, chatMessage{role: "error", content: fmt.Sprintf("Error: %s", msg.err.Error())})
		} else {
			m.messages = append(m.messages, chatMessage{role: "assistant", content: msg.response})
		}
		m.updateViewport()

	case error:
		m.err = msg
		m.processing = false
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *model) updateViewport() {
	var content strings.Builder
	
	content.WriteString(titleStyle.Render("🤖 Stacks AI Assistant"))
	content.WriteString("\n\n")
	
	for _, msg := range m.messages {
		switch msg.role {
		case "user":
			content.WriteString(userStyle.Render("You: "))
			content.WriteString(messageStyle.Render(msg.content))
		case "assistant":
			content.WriteString(assistantStyle.Render("AI: "))
			content.WriteString(messageStyle.Render(msg.content))
		case "error":
			content.WriteString(errorStyle.Render("Error: "))
			content.WriteString(messageStyle.Render(msg.content))
		}
		content.WriteString("\n\n")
	}

	if m.processing {
		content.WriteString(assistantStyle.Render("AI: "))
		content.WriteString(messageStyle.Render("🤔 Thinking..."))
		content.WriteString("\n\n")
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoBottom()
}

func (m model) processQuery(query string) tea.Cmd {
	return func() tea.Msg {
		response, err := m.agent.ProcessQuery(query)
		return processedMsg{response: response, err: err}
	}
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	help := helpStyle.Render("Press Ctrl+C to quit • Enter to send message • Use natural language to query Stacks blockchain")
	
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.viewport.View(),
		m.textarea.View(),
		help,
	)
} 