# Chrome Setup for MCP Automation

This guide explains how to set up Chrome to work with the MCP Chrome Automation service while **preserving your existing sessions, cookies, and logged-in accounts**.

## Why Connect to Existing Chrome?

Instead of spawning a new Chrome instance, this service connects to your existing Chrome browser, which means:

✅ **Your sessions are preserved** - Stay logged into Gmail, GitHub, etc.  
✅ **Cookies remain intact** - No need to re-authenticate  
✅ **Extensions work** - Your ad blockers, password managers stay active  
✅ **Real-world testing** - Automation happens in your actual browsing environment  
✅ **Seamless experience** - No separate browser windows to manage  

## Quick Setup

### Option 1: Automated Setup (Recommended)

```bash
# Run the setup script
./start_chrome.sh
```

This script will:
1. Detect your Chrome installation
2. Stop existing Chrome processes safely
3. Restart Chrome with debugging enabled
4. Verify the connection is working

### Option 2: Manual Setup

**macOS:**
```bash
# Close Chrome completely first
pkill -f "Google Chrome"

# Start Chrome with debugging
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222 \
  --remote-debugging-address=127.0.0.1 \
  --user-data-dir="$HOME/.chrome-mcp-debug"
```

**Linux:**
```bash
# Close Chrome completely first
pkill -f chrome

# Start Chrome with debugging
google-chrome \
  --remote-debugging-port=9222 \
  --remote-debugging-address=127.0.0.1 \
  --user-data-dir="$HOME/.chrome-mcp-debug"
```

**Windows:**
```cmd
REM Close Chrome completely first
taskkill /f /im chrome.exe

REM Start Chrome with debugging
"C:\Program Files\Google\Chrome\Application\chrome.exe" ^
  --remote-debugging-port=9222 ^
  --remote-debugging-address=127.0.0.1 ^
  --user-data-dir="%USERPROFILE%\.chrome-mcp-debug"
```

## 🔧 What the Flags Do

- `--remote-debugging-port=9222` - Opens debugging API on port 9222
- `--remote-debugging-address=127.0.0.1` - Only allows local connections
- `--user-data-dir` - Uses a separate profile (optional, preserves your main profile)

## Verify Setup

After starting Chrome, verify the debugging port is active:

```bash
# Should return JSON with browser info
curl http://localhost:9222/json

# Or check in your browser
open http://localhost:9222
```

## Usage Patterns

### Pattern 1: Daily Automation
```bash
# Start Chrome with debugging (once per day)
./start_chrome.sh

# Use MCP service throughout the day
./mcp-chromautomation server
```

### Pattern 2: Session-Aware Automation
```bash
# 1. Start Chrome and log into your accounts normally
./start_chrome.sh

# 2. Browse and log into Gmail, GitHub, etc.
# 3. Now run automation that uses your logged-in sessions
./mcp-chromautomation server
```

### Pattern 3: Development Workflow
```bash
# Keep Chrome with debugging always running
./start_chrome.sh

# Run multiple automation sessions
./mcp-chromautomation ui        # Interactive testing
./mcp-chromautomation server    # MCP server for clients
```

## Security Considerations

**The debugging port is powerful** - it allows full control over Chrome:

**Safe practices:**
- Only bind to localhost (127.0.0.1)
- Use specific user data directory
- Close debugging when not needed

**Avoid:**
- Exposing debugging port to network
- Running debugging mode on production machines
- Leaving debugging enabled permanently

## Troubleshooting

### Chrome Won't Start with Debugging
```bash
# Check if port is already in use
lsof -i :9222

# Kill any processes using the port
kill $(lsof -t -i:9222)

# Try starting Chrome again
./start_chrome.sh
```

### MCP Service Can't Connect
```bash
# Verify Chrome debugging is working
curl http://localhost:9222/json

# Check Chrome is running with correct flags
ps aux | grep chrome | grep remote-debugging
```

### Multiple Chrome Instances
```bash
# Close all Chrome instances
pkill -f chrome

# Start fresh
./start_chrome.sh
```

### Extensions Not Working
If you need your extensions, remove the `--user-data-dir` flag to use your main Chrome profile:

```bash
google-chrome \
  --remote-debugging-port=9222 \
  --remote-debugging-address=127.0.0.1
  # No --user-data-dir flag
```

## Integration with MCP Service

Once Chrome is running with debugging, the MCP service will:

1. **Auto-detect** Chrome on port 9222
2. **Connect seamlessly** to existing tabs
3. **Preserve sessions** - no re-authentication needed
4. **Respect your setup** - works with your extensions and settings

## Pro Tips

### Permanent Setup
Add to your shell profile (`.bashrc`, `.zshrc`):
```bash
alias chrome-debug='./start_chrome.sh'
alias mcp-chrome='./mcp-chromautomation server'
```

### Multiple Debugging Ports
Run multiple Chrome instances on different ports:
```bash
# Instance 1 (personal)
google-chrome --remote-debugging-port=9222 --user-data-dir="$HOME/.chrome-personal"

# Instance 2 (work)  
google-chrome --remote-debugging-port=9223 --user-data-dir="$HOME/.chrome-work"
```

### Automation-Friendly Bookmarks
Create bookmarks for commonly automated sites:
- Testing environments
- Admin panels  
- Forms you frequently fill

This makes navigation faster and more reliable for automation scripts.

---

** You're now ready for session-aware Chrome automation!**

Your MCP service will connect to your existing Chrome instance, preserving all your logged-in accounts and settings while providing powerful browser automation capabilities.
