#!/usr/bin/env node

/**
 * Example MCP client for Chrome Automation Service
 * 
 * This demonstrates how to use the MCP Chrome Automation service
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
        console.log('🚀 Starting MCP Chrome Automation Service...');
        
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
        console.log('✅ MCP Chrome Automation Service initialized');
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
                        reject(new Error(`${response.error.message}: ${response.error.data}`));
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
                name: 'mcp-chrome-example',
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

    // Convenience methods for common operations
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

// Example usage
async function runExamples() {
    const client = new MCPChromeClient();
    
    try {
        await client.start();
        
        console.log('\n📋 Available tools:');
        const tools = await client.listTools();
        tools.tools.forEach(tool => {
            console.log(`  🔧 ${tool.name}: ${tool.description}`);
        });

        console.log('\n🌐 Example 1: Navigate to website');
        const navResult = await client.navigate('https://httpbin.org', {
            screenshot: true,
            waitFor: 'body'
        });
        console.log('✅ Navigation result:', navResult);

        console.log('\n📄 Example 2: Extract text');
        const textResult = await client.extractText('h1, h2');
        console.log('✅ Extracted text:', textResult);

        console.log('\n📝 Example 3: Navigate to form and fill it');
        await client.navigate('https://httpbin.org/forms/post');
        
        const formResult = await client.fillForm({
            'input[name="custname"]': 'John Doe',
            'input[name="custtel"]': '+1234567890',
            'input[name="custemail"]': 'john@example.com'
        });
        console.log('✅ Form filled:', formResult);

        console.log('\n💾 Example 4: Save session');
        const sessionResult = await client.saveSession('example_session_' + Date.now());
        console.log('✅ Session saved:', sessionResult);

        console.log('\n🎉 All examples completed successfully!');

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

// Run the examples
if (require.main === module) {
    runExamples().catch(console.error);
}

module.exports = MCPChromeClient;