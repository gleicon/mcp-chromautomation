# MCP Chrome Automation Service

A comprehensive Model Context Protocol (MCP) service for Chrome browser automation with enhanced capabilities, built on the [browserhttp](https://github.com/gleicon/browserhttp) library and using the [mcp-go](https://github.com/mark3labs/mcp-go) framework.

## Features

### Enhanced MCP Server Capabilities (19 Tools)

**Core Browser Automation:**
- **Enhanced Navigation**: Navigate with performance tracking and advanced waiting
- **Advanced Interaction**: Click, type, select with improved reliability
- **Content Extraction**: Text, links, images, forms with metadata
- **Session Management**: Complete cookie and localStorage handling
- **Screenshot Capture**: Automated screenshot with custom naming

**New Analysis Tools:**
- **Performance Monitoring**: DOM load times, network requests, resource analysis
- **SEO Analysis**: Title, description, keywords, heading structure
- **Security Scanning**: SSL validation, CSP analysis, vulnerability detection
- **Content Intelligence**: Comprehensive link categorization and form analysis

**Advanced Features:**
- **JSON API Integration**: Send POST requests with browser context
- **Local Storage Management**: Read/write browser localStorage
- **Multi-condition Waiting**: Wait for elements, text, or navigation
- **Session Cleanup**: Complete cookie and storage clearing

### Beautiful CLI Interface
- **Interactive Menu**: Navigate through options with keyboard
- **Real-time Feedback**: See results immediately
- **Session Browser**: View and manage saved sessions
- **Request Logs**: Monitor all browser activity
- **Settings Panel**: Database statistics and configuration

### Enhanced Browser Automation
- Built on your proven `browserhttp` library
- Real Chrome browser via chromedp
- JavaScript rendering and form submission
- Screenshot capture and storage
- Persistent tab management

## Architecture

```
mcp-chromautomation/
├── main.go                     # Entry point with CLI commands
├── internal/
│   ├── server/                 # MCP protocol implementation
│   │   └── server.go          # JSON-RPC server with tools
│   ├── browser/               # Browser automation layer
│   │   └── manager.go         # Enhanced browserhttp integration
│   ├── storage/               # Local data persistence
│   │   └── manager.go         # SQLite database management
│   └── cli/                   # Terminal user interface
│       └── app.go             # Bubble Tea application
└── go.mod                     # Dependencies
```

## Installation

1. **Clone and setup**:
   ```bash
   git clone https://github.com/gleicon/mcp-chromautomation
   cd mcp-chromautomation
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   go build -o mcp-chromautomation
   ```

4. **Setup Chrome for automation** (preserves your existing sessions):
   ```bash
   # Start Chrome with debugging enabled (keeps your sessions & cookies)
   ./start_chrome.sh
   ```

## Usage

### MCP Server Mode

**First, ensure Chrome is running with debugging enabled:**
```bash
# This preserves your existing sessions and cookies
./start_chrome.sh
```

**Then start the MCP server:**
```bash
./mcp-chromautomation server
```

The service will connect to your existing Chrome instance, preserving all your:
- 🔐 Logged-in sessions
- 🍪 Cookies and authentication
- 📝 Form data and preferences
- 🔖 Bookmarks and extensions

The server provides 19 enhanced MCP tools:

**Core Automation:**
- `chrome_navigate` - Enhanced navigation with performance tracking
- `chrome_click` - Reliable element clicking with validation
- `chrome_extract_text` - Advanced text extraction with selectors
- `chrome_fill_form` - Smart form filling with validation
- `chrome_screenshot` - Screenshot capture with custom naming
- `chrome_wait_for_element` - Element waiting with timeout control

**Content Analysis:**
- `chrome_extract_links` - Extract and categorize all page links
- `chrome_extract_images` - Extract images with metadata
- `chrome_extract_forms` - Analyze form structures and fields
- `chrome_analyze_seo` - Comprehensive SEO analysis

**Performance & Security:**
- `chrome_get_performance` - Detailed performance metrics
- `chrome_check_security` - Security vulnerability scanning

**Advanced Interaction:**
- `chrome_post_json` - Send JSON data via browser
- `chrome_wait_advanced` - Multi-condition waiting

**Session Management:**
- `chrome_get_local_storage` - Access browser localStorage
- `chrome_set_local_storage` - Manage localStorage data
- `chrome_clear_session` - Complete session cleanup
- `session_save` - Enhanced session saving with full state
- `session_load` - Complete session restoration

### Interactive CLI Mode
Launch the beautiful terminal interface:

```bash
./mcp-chromautomation ui
```

Navigate with keyboard shortcuts:
- `↑/↓` or `j/k` - Move up/down
- `Enter` - Select item
- `Esc` - Go back
- `q` - Quit

## MCP Integration

Use with any MCP client by configuring the server:

```json
{
  "mcpServers": {
    "chromautomation": {
      "command": "/path/to/mcp-chromautomation",
      "args": ["server"]
    }
  }
}
```

Example tool calls:

```javascript
// Navigate to a website
await mcp.callTool("chrome_navigate", {
  url: "https://example.com",
  screenshot: true,
  wait_for: "#main-content"
});

// Click an element
await mcp.callTool("chrome_click", {
  selector: ".submit-button",
  screenshot: true
});

// Extract text content
const result = await mcp.callTool("chrome_extract_text", {
  selector: ".article-content p"
});
```

## Technical Architecture

### JavaScript Integration & Security Model

The service uses **legitimate browser automation** through Chrome DevTools Protocol (CDP), not malicious code injection.

#### How JavaScript Execution Works

**1. Chrome DevTools Protocol (CDP)**
```go
// Uses Chrome's official debugging protocol
import "github.com/chromedp/chromedp"
```

**2. Safe JavaScript Execution Types**

**DOM Queries (Read-only):**
```javascript
// Extract links using standard DOM API
Array.from(document.querySelectorAll('a[href]')).map(a => a.href)
```

**Performance Metrics (Browser APIs):**
```javascript
// Get performance data using standard Performance API
(() => {
    const perf = performance.getEntriesByType('navigation')[0];
    return {
        domContentLoaded: perf.domContentLoadedEventEnd - perf.domContentLoadedEventStart,
        loadComplete: perf.loadEventEnd - perf.loadEventStart
    };
})()
```

**Session Management (Standard Web APIs):**
```javascript
// Cookie management via document.cookie
document.cookie = 'name=value; path=/';

// localStorage access
localStorage.getItem('key');
localStorage.setItem('key', 'value');
```

#### Security Model

**What it IS:**
- Legitimate browser automation using Chrome's official debugging protocol
- Execution equivalent to manual browser DevTools console usage
- Read-only operations for most functions (link extraction, text parsing)
- Standard DOM/Web API usage only

**What it's NOT:**
- Not XSS injection into target websites
- Not malicious code execution
- Not breaking website security policies
- Not bypassing same-origin policies maliciously

**Flow Overview:**
```
1. Connect to existing Chrome via CDP
2. Navigate to page normally (like a user)
3. Wait for page to load completely
4. Execute JavaScript in browser context (like F12 console)
5. Extract data using standard DOM APIs
6. Return results safely
```

### Performance Analysis

The service provides detailed performance monitoring:

```javascript
// Real browser performance metrics
const metrics = await mcp.callTool("chrome_get_performance", {});
// Returns: DOM load times, network requests, resource sizes
```

### Content Intelligence

Advanced content extraction and analysis:

```javascript
// Extract and categorize all links
const links = await mcp.callTool("chrome_extract_links", {});
// Returns: internal, external, product, category links with counts

// SEO analysis
const seo = await mcp.callTool("chrome_analyze_seo", {});
// Returns: title, description, keywords, heading structure
```

### Security Scanning

Built-in security analysis tools:

```javascript
// Comprehensive security check
const security = await mcp.callTool("chrome_check_security", {});
// Returns: SSL status, CSP analysis, vulnerability scan results
```

## Practical Examples

### Complete Website Analysis

```javascript
// Comprehensive analysis workflow
const client = new MCPChromeClient();

// Navigate with performance tracking
const nav = await client.callTool('chrome_navigate', {
    url: 'https://example.com',
    track_performance: true,
    screenshot: true
});

// Extract all content
const links = await client.callTool('chrome_extract_links', {});
const seo = await client.callTool('chrome_analyze_seo', {});
const performance = await client.callTool('chrome_get_performance', {});
const security = await client.callTool('chrome_check_security', {});

console.log(`Found ${links.count} links`);
console.log(`SEO score: ${seo.seo.title ? 'Has title' : 'Missing title'}`);
console.log(`Load time: ${performance.performance.load_complete}ms`);
console.log(`SSL valid: ${security.security.ssl.valid}`);
```

### E-commerce Site Analysis

```javascript
// Example: Analyze Magazine Luiza
const result = await client.callTool('chrome_navigate', {
    url: 'https://magazineluiza.com.br'
});

const links = await client.callTool('chrome_extract_links', {});
const linkData = JSON.parse(links.content[0].text);

console.log('Link Analysis:');
console.log(`Total links: ${linkData.count}`);
console.log(`Product links: ${linkData.links.filter(l => 
    l.includes('/produto/') || l.includes('/p/')
).length}`);
console.log(`Internal links: ${linkData.links.filter(l => 
    l.includes('magazineluiza.com.br')
).length}`);
```

### Form Analysis and Interaction

```javascript
// Extract and analyze forms
const forms = await client.callTool('chrome_extract_forms', {});
const formData = JSON.parse(forms.content[0].text);

console.log(`Found ${formData.count} forms`);

// Fill a form intelligently
await client.callTool('chrome_fill_form', {
    fields: {
        '#email': 'user@example.com',
        '#password': 'secure_password',
        '#remember': 'true'
    },
    submit: true,
    validate_before_submit: true
});
```

### Performance Monitoring

```javascript
// Monitor page performance over time
const sites = ['https://example.com', 'https://google.com'];
const results = [];

for (const site of sites) {
    await client.callTool('chrome_navigate', { 
        url: site, 
        track_performance: true 
    });
    
    const perf = await client.callTool('chrome_get_performance', {});
    const perfData = JSON.parse(perf.content[0].text);
    
    results.push({
        site,
        loadTime: perfData.performance.load_complete,
        domReady: perfData.performance.dom_content_loaded,
        requests: perfData.performance.network_requests
    });
}

console.table(results);
```

## Data Storage

The service stores data locally in `~/.mcp-chromautomation/`:

- **Sessions**: Browser state including cookies and URLs
- **Request Logs**: Complete HTTP request/response history  
- **Screenshots**: Captured page screenshots with metadata
- **Performance Data**: Page load metrics and analysis
- **Settings**: User preferences and configuration

## Dependencies

### Core Libraries
- **[browserhttp](https://github.com/gleicon/browserhttp)** - Your excellent browser automation library
- **[mcp-go](https://github.com/mark3labs/mcp-go)** - Robust MCP protocol implementation
- **[chromedp](https://github.com/chromedp/chromedp)** - Chrome DevTools Protocol
- **[cobra](https://github.com/spf13/cobra)** - CLI framework
- **[sqlite](https://modernc.org/sqlite)** - Pure Go SQLite

### UI Framework
- **[bubbletea](https://github.com/charmbracelet/bubbletea)** - TUI framework
- **[bubbles](https://github.com/charmbracelet/bubbles)** - TUI components  
- **[lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Built on [browserhttp](https://github.com/gleicon/browserhttp) library
- Inspired by the Model Context Protocol specification
- UI powered by the amazing Charm libraries


