#!/usr/bin/env node

/**
 * Updated MCP client example for Chrome Automation Service with mark3labs/mcp-go
 * 
 * This demonstrates how to use the enhanced MCP Chrome Automation service
 * from a JavaScript/Node.js client application.
 */

const { spawn } = require('child_process');
const readline = require('readline');

class MCPChromeClient {
    constructor() {
        this.process = null;
        this.requestId = 1;
        this.pendingRequests = new Map();
    }

    async start() {
        console.log('🚀 Starting Enhanced MCP Chrome Automation Service...');
        
        this.process = spawn('./mcp-chromautomation', ['server'], {
            stdio: ['pipe', 'pipe', 'inherit']
        });

        this.process.stdout.setEncoding('utf8');
        this.process.stdout.on('data', (data) => {
            this.handleResponse(data);
        });

        this.process.on('error', (error) => {
            console.error('❌ Process error:', error);
        });

        // Initialize the MCP session
        await this.initialize();
        console.log('✅ Enhanced MCP Chrome Automation Service initialized');
    }

    async stop() {
        if (this.process) {
            this.process.kill();
            console.log('🛑 MCP service stopped');
        }
    }

    sendRequest(method, params = {}) {
        return new Promise((resolve, reject) => {
            const id = this.requestId++;
            const request = {
                jsonrpc: '2.0',
                id: id,
                method: method,
                params: params
            };

            this.pendingRequests.set(id, { resolve, reject });
            
            const requestStr = JSON.stringify(request) + '\n';
            this.process.stdin.write(requestStr);
        });
    }

    handleResponse(data) {
        const lines = data.trim().split('\n');
        
        for (const line of lines) {
            try {
                const response = JSON.parse(line);
                
                if (response.id && this.pendingRequests.has(response.id)) {
                    const { resolve, reject } = this.pendingRequests.get(response.id);
                    this.pendingRequests.delete(response.id);
                    
                    if (response.error) {
                        reject(new Error(`${response.error.message}: ${response.error.data || ''}`));
                    } else {
                        resolve(response.result);
                    }
                }
            } catch (error) {
                console.error('❌ Failed to parse response:', error);
            }
        }
    }

    async initialize() {
        return this.sendRequest('initialize', {
            protocolVersion: '2024-11-05',
            capabilities: {},
            clientInfo: {
                name: 'mcp-chrome-example-enhanced',
                version: '1.0.0'
            }
        });
    }

    async listTools() {
        return this.sendRequest('tools/list');
    }

    async callTool(name, arguments) {
        return this.sendRequest('tools/call', {
            name: name,
            arguments: arguments
        });
    }

    // Enhanced convenience methods for new tools
    async navigate(url, options = {}) {
        return this.callTool('chrome_navigate', {
            url: url,
            wait_for: options.waitFor,
            screenshot: options.screenshot || false
        });
    }

    async click(selector, options = {}) {
        return this.callTool('chrome_click', {
            selector: selector,
            screenshot: options.screenshot || false
        });
    }

    async extractText(selector) {
        return this.callTool('chrome_extract_text', {
            selector: selector
        });
    }

    async fillForm(fields, submit = false) {
        return this.callTool('chrome_fill_form', {
            fields: fields,
            submit: submit
        });
    }

    async takeScreenshot(filename) {
        return this.callTool('chrome_screenshot', {
            filename: filename
        });
    }

    async waitForElement(selector, timeout = 10) {
        return this.callTool('chrome_wait_for_element', {
            selector: selector,
            timeout: timeout
        });
    }

    async saveSession(name) {
        return this.callTool('session_save', {
            name: name
        });
    }

    async loadSession(name) {
        return this.callTool('session_load', {
            name: name
        });
    }
}

// Enhanced example usage with new features
async function runEnhancedExamples() {
    const client = new MCPChromeClient();
    
    try {
        await client.start();
        
        console.log('\n📋 Available tools:');
        const tools = await client.listTools();
        tools.tools.forEach(tool => {
            console.log(`  🔧 ${tool.name}: ${tool.description}`);
        });

        console.log('\n🌐 Example 1: Navigate with element waiting');
        const navResult = await client.navigate('https://httpbin.org/html', {
            screenshot: true,
            waitFor: 'h1'
        });
        console.log('✅ Navigation with waiting result:', JSON.parse(navResult.content[0].text));

        console.log('\n⏱️ Example 2: Wait for specific element');
        const waitResult = await client.waitForElement('h1', 5);
        console.log('✅ Element wait result:', JSON.parse(waitResult.content[0].text));

        console.log('\n📸 Example 3: Take a screenshot');
        const screenshotResult = await client.takeScreenshot('example_page.png');
        console.log('✅ Screenshot result:', JSON.parse(screenshotResult.content[0].text));

        console.log('\n📄 Example 4: Extract multiple text elements');
        const textResult = await client.extractText('h1, p');
        console.log('✅ Extracted text:', JSON.parse(textResult.content[0].text));

        console.log('\n📝 Example 5: Navigate to form and interact');
        await client.navigate('https://httpbin.org/forms/post');
        
        // Wait for form to load
        await client.waitForElement('form', 3);
        
        const formResult = await client.fillForm({
            'input[name="custname"]': 'John Doe Enhanced',
            'input[name="custtel"]': '+1234567890',
            'input[name="custemail"]': 'john.enhanced@example.com'
        });
        console.log('✅ Enhanced form filled:', JSON.parse(formResult.content[0].text));

        console.log('\n💾 Example 6: Enhanced session management');
        const sessionResult = await client.saveSession('enhanced_session_' + Date.now());
        console.log('✅ Enhanced session saved:', JSON.parse(sessionResult.content[0].text));

        console.log('\n🎉 All enhanced examples completed successfully!');

    } catch (error) {
        console.error('❌ Error:', error.message);
    } finally {
        await client.stop();
    }
}

// Handle graceful shutdown
process.on('SIGINT', () => {
    console.log('\n👋 Shutting down gracefully...');
    process.exit(0);
});

// Run the enhanced examples
if (require.main === module) {
    runEnhancedExamples().catch(console.error);
}

module.exports = MCPChromeClient;