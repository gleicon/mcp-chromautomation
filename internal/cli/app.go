package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gleicon/mcp-chromautomation/internal/browser"
	"github.com/gleicon/mcp-chromautomation/internal/storage"
)

// Screen types
type screen int

const (
	menuScreen screen = iota
	navigateScreen
	sessionScreen
	logsScreen
	settingsScreen
)

// App represents the main application state
type App struct {
	screen      screen
	width       int
	height      int
	browser     *browser.Manager
	storage     *storage.Manager
	
	// UI Components
	menu        list.Model
	urlInput    textinput.Model
	sessionList list.Model
	logList     list.Model
	help        help.Model
	
	// State
	message     string
	lastResult  interface{}
	keys        keyMap
}

// keyMap defines keyboard shortcuts
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
	Help   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Quit, k.Help},
	}
}

// MenuItem represents a menu item
type MenuItem struct {
	title       string
	description string
	action      string
}

func (i MenuItem) FilterValue() string { return i.title }
func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }

// Start starts the CLI application
func Start() error {
	app := &App{
		screen: menuScreen,
		keys: keyMap{
			Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
			Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
			Left:  key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "left")),
			Right: key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "right")),
			Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
			Back:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
			Quit:  key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
			Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		},
		help: help.New(),
	}
	
	// Initialize components
	if err := app.init(); err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}
	
	p := tea.NewProgram(app, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

// init initializes the application
func (a *App) init() error {
	// Initialize browser and storage
	a.browser = browser.New()
	a.storage = storage.New()
	
	if err := a.browser.Init(); err != nil {
		return fmt.Errorf("failed to initialize browser: %w", err)
	}
	
	if err := a.storage.Init(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	
	// Initialize menu
	menuItems := []list.Item{
		MenuItem{title: "🌐 Navigate", description: "Navigate to a website", action: "navigate"},
		MenuItem{title: "🎯 Click Element", description: "Click on a page element", action: "click"},
		MenuItem{title: "📝 Fill Form", description: "Fill out form fields", action: "form"},
		MenuItem{title: "📄 Extract Text", description: "Extract text from elements", action: "text"},
		MenuItem{title: "💾 Sessions", description: "Manage browser sessions", action: "session"},
		MenuItem{title: "📊 Logs", description: "View request logs", action: "logs"},
		MenuItem{title: "⚙️  Settings", description: "Application settings", action: "settings"},
		MenuItem{title: "❌ Quit", description: "Exit application", action: "quit"},
	}
	
	a.menu = list.New(menuItems, a.itemDelegate(), 0, 0)
	a.menu.Title = "🤖 Chrome Automation MCP Service"
	a.menu.SetShowStatusBar(false)
	a.menu.SetFilteringEnabled(false)
	a.menu.Styles.Title = titleStyle
	
	// Initialize URL input
	a.urlInput = textinput.New()
	a.urlInput.Placeholder = "Enter URL (e.g., https://example.com)"
	a.urlInput.Focus()
	a.urlInput.CharLimit = 500
	a.urlInput.Width = 60
	
	return nil
}

// itemDelegate creates a custom item delegate for the list
func (a *App) itemDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = selectedItemStyle
	d.Styles.SelectedDesc = selectedDescStyle
	return d
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.menu.SetWidth(msg.Width)
		a.menu.SetHeight(msg.Height - 4)
		return a, nil
		
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, a.keys.Quit):
			a.cleanup()
			return a, tea.Quit
			
		case key.Matches(msg, a.keys.Back):
			if a.screen != menuScreen {
				a.screen = menuScreen
				a.message = ""
				return a, nil
			}
			
		case key.Matches(msg, a.keys.Enter):
			return a.handleEnter()
		}
	}
	
	// Update components based on current screen
	var cmd tea.Cmd
	switch a.screen {
	case menuScreen:
		a.menu, cmd = a.menu.Update(msg)
	case navigateScreen:
		a.urlInput, cmd = a.urlInput.Update(msg)
	}
	
	return a, cmd
}

// View implements tea.Model
func (a *App) View() string {
	switch a.screen {
	case menuScreen:
		return a.menuView()
	case navigateScreen:
		return a.navigateView()
	case sessionScreen:
		return a.sessionView()
	case logsScreen:
		return a.logsView()
	case settingsScreen:
		return a.settingsView()
	default:
		return a.menuView()
	}
}

// handleEnter handles the enter key press
func (a *App) handleEnter() (tea.Model, tea.Cmd) {
	switch a.screen {
	case menuScreen:
		if item, ok := a.menu.SelectedItem().(MenuItem); ok {
			return a.handleMenuAction(item.action)
		}
	case navigateScreen:
		return a.handleNavigate()
	}
	return a, nil
}

// handleMenuAction handles menu item selection
func (a *App) handleMenuAction(action string) (tea.Model, tea.Cmd) {
	switch action {
	case "navigate":
		a.screen = navigateScreen
		a.urlInput.SetValue("")
		a.urlInput.Focus()
		return a, nil
	case "session":
		a.screen = sessionScreen
		return a, tea.Cmd(a.loadSessions)
	case "logs":
		a.screen = logsScreen
		return a, tea.Cmd(a.loadLogs)
	case "settings":
		a.screen = settingsScreen
		return a, nil
	case "quit":
		a.cleanup()
		return a, tea.Quit
	}
	return a, nil
}

// handleNavigate handles URL navigation
func (a *App) handleNavigate() (tea.Model, tea.Cmd) {
	url := a.urlInput.Value()
	if url == "" {
		a.message = "❌ Please enter a URL"
		return a, nil
	}
	
	// Navigate using browser manager
	result, err := a.browser.Navigate(url, "", true)
	if err != nil {
		a.message = fmt.Sprintf("❌ Navigation failed: %v", err)
	} else {
		if result.Success {
			a.message = fmt.Sprintf("✅ Successfully navigated to %s", result.URL)
			if result.Title != "" {
				a.message += fmt.Sprintf("\n📄 Page title: %s", result.Title)
			}
		} else {
			a.message = fmt.Sprintf("❌ %s", result.Message)
		}
	}
	
	a.lastResult = result
	return a, nil
}

// View methods
func (a *App) menuView() string {
	var b strings.Builder
	
	b.WriteString(a.menu.View())
	
	if a.message != "" {
		b.WriteString("\n\n")
		b.WriteString(messageStyle.Render(a.message))
	}
	
	b.WriteString("\n\n")
	b.WriteString(a.help.View(a.keys))
	
	return b.String()
}

func (a *App) navigateView() string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("🌐 Navigate to Website"))
	b.WriteString("\n\n")
	b.WriteString("Enter the URL you want to visit:\n\n")
	b.WriteString(a.urlInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press Enter to navigate • Esc to go back"))
	
	if a.message != "" {
		b.WriteString("\n\n")
		b.WriteString(messageStyle.Render(a.message))
	}
	
	return b.String()
}

func (a *App) sessionView() string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("💾 Browser Sessions"))
	b.WriteString("\n\n")
	
	sessions, err := a.storage.ListSessions()
	if err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error loading sessions: %v", err)))
	} else if len(sessions) == 0 {
		b.WriteString(helpStyle.Render("No saved sessions found"))
	} else {
		for _, session := range sessions {
			sessionInfo := fmt.Sprintf("📁 %s\n   URL: %s\n   Updated: %s\n", 
				session.Name, 
				session.Data.URL, 
				session.UpdatedAt.Format("2006-01-02 15:04:05"))
			b.WriteString(sessionItemStyle.Render(sessionInfo))
			b.WriteString("\n")
		}
	}
	
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Esc to go back"))
	
	return b.String()
}

func (a *App) logsView() string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("📊 Request Logs"))
	b.WriteString("\n\n")
	
	logs, err := a.storage.GetRequestLogs(20)
	if err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error loading logs: %v", err)))
	} else if len(logs) == 0 {
		b.WriteString(helpStyle.Render("No request logs found"))
	} else {
		for _, log := range logs {
			status := "✅"
			if log.Status >= 400 {
				status = "❌"
			} else if log.Status >= 300 {
				status = "⚠️"
			}
			
			logInfo := fmt.Sprintf("%s %s %s (%d) - %dms\n   %s\n", 
				status, log.Method, log.URL, log.Status, log.Duration, 
				log.Timestamp.Format("2006-01-02 15:04:05"))
			b.WriteString(logItemStyle.Render(logInfo))
			b.WriteString("\n")
		}
	}
	
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Esc to go back"))
	
	return b.String()
}

func (a *App) settingsView() string {
	var b strings.Builder
	
	b.WriteString(headerStyle.Render("⚙️ Settings"))
	b.WriteString("\n\n")
	
	stats, err := a.storage.GetStats()
	if err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error loading stats: %v", err)))
	} else {
		b.WriteString("📊 Database Statistics:\n\n")
		for key, value := range stats {
			b.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
		}
	}
	
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Esc to go back"))
	
	return b.String()
}

// Helper methods
func (a *App) loadSessions() tea.Msg {
	// This would be used for async loading if needed
	return nil
}

func (a *App) loadLogs() tea.Msg {
	// This would be used for async loading if needed
	return nil
}

func (a *App) cleanup() {
	if a.browser != nil {
		a.browser.Close()
	}
	if a.storage != nil {
		a.storage.Close()
	}
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		MarginLeft(2)
	
	headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#2563EB")).
		MarginBottom(1)
	
	selectedItemStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED"))
	
	selectedDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B"))
	
	messageStyle = lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#10B981"))
	
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#EF4444")).
		Bold(true)
	
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B")).
		Italic(true)
	
	sessionItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#3B82F6"))
	
	logItemStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color("#8B5CF6"))
)