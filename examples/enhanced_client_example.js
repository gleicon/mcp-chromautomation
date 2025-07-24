#!/usr/bin/env node

/**
 * Enhanced MCP client example showcasing all new browserhttp capabilities
 * 
 * This demonstrates the full power of the enhanced MCP Chrome Automation service
 * with all the new features from the updated browserhttp library.
 */

const { spawn } = require('child_process');

class EnhancedMCPChromeClient {
    constructor() {
        this.process = null;
        this.requestId = 1;
        this.pendingRequests = new Map();
    }

    async start() {
        console.log('🚀 Starting Enhanced MCP Chrome Automation Service with Full Capabilities...');
        
        this.process = spawn('./build/mcp-chromautomation', ['server'], {
            stdio: ['pipe', 'pipe', 'inherit']
        });

        this.process.stdout.setEncoding('utf8');
        this.process.stdout.on('data', (data) => {
            this.handleResponse(data);
        });

        this.process.on('error', (error) => {
            console.error('❌ Process error:', error);
        });

        await this.initialize();
        console.log('✅ Enhanced MCP Chrome Automation Service fully initialized');
    }

    async stop() {
        if (this.process) {
            this.process.kill();
            console.log('🛑 Enhanced MCP service stopped');
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
                name: 'enhanced-mcp-chrome-client',
                version: '2.0.0'
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

    // ========== ENHANCED CORE METHODS ==========
    
    async navigateEnhanced(url, options = {}) {
        return this.callTool('chrome_navigate', {
            url: url,
            wait_for: options.waitFor,
            screenshot: options.screenshot || false,
            track_performance: options.trackPerformance || false
        });
    }

    async fillFormEnhanced(fields, options = {}) {
        return this.callTool('chrome_fill_form', {
            fields: fields,
            submit: options.submit || false,
            validate_before_submit: options.validate !== false
        });
    }

    // ========== NEW CONTENT ANALYSIS METHODS ==========
    
    async extractLinks(filter) {
        return this.callTool('chrome_extract_links', {
            filter: filter
        });
    }

    async extractImages(includeMetadata = true) {
        return this.callTool('chrome_extract_images', {
            include_metadata: includeMetadata
        });
    }

    async extractForms() {
        return this.callTool('chrome_extract_forms', {});
    }

    async analyzeSEO() {
        return this.callTool('chrome_analyze_seo', {});
    }

    // ========== NEW PERFORMANCE ANALYSIS METHODS ==========
    
    async getPerformanceMetrics() {
        return this.callTool('chrome_get_performance', {});
    }

    // ========== NEW SECURITY ANALYSIS METHODS ==========
    
    async checkSecurity() {
        return this.callTool('chrome_check_security', {});
    }

    // ========== NEW ADVANCED INTERACTION METHODS ==========
    
    async postJSON(url, data) {
        return this.callTool('chrome_post_json', {
            url: url,
            data: data
        });
    }

    async waitAdvanced(options = {}) {
        return this.callTool('chrome_wait_advanced', {
            selector: options.selector,
            text: options.text,
            timeout: options.timeout || 10
        });
    }

    // ========== NEW DATA MANAGEMENT METHODS ==========
    
    async getLocalStorage(key) {
        return this.callTool('chrome_get_local_storage', {
            key: key
        });
    }

    async setLocalStorage(key, value) {
        return this.callTool('chrome_set_local_storage', {
            key: key,
            value: value
        });
    }

    async clearSession() {
        return this.callTool('chrome_clear_session', {});
    }
}

// Enhanced comprehensive examples showcasing all new capabilities
async function runEnhancedExamples() {
    const client = new EnhancedMCPChromeClient();
    
    try {
        await client.start();
        
        console.log('\n📋 Available Enhanced Tools:');
        const tools = await client.listTools();
        tools.tools.forEach(tool => {
            console.log(`  🔧 ${tool.name}: ${tool.description}`);
        });

        console.log('\n🌐 Enhanced Example 1: Navigate with Performance Tracking');
        const navResult = await client.navigateEnhanced('https://example.com', {
            screenshot: true,
            waitFor: 'h1',
            trackPerformance: true
        });
        console.log('✅ Enhanced navigation:', JSON.parse(navResult.content[0].text));

        console.log('\n📊 Enhanced Example 2: Get Performance Metrics');
        const performance = await client.getPerformanceMetrics();
        console.log('✅ Performance metrics:', JSON.parse(performance.content[0].text));

        console.log('\n🔗 Enhanced Example 3: Extract All Links');
        const links = await client.extractLinks();
        const linkData = JSON.parse(links.content[0].text);
        console.log(`✅ Extracted ${linkData.count} links from the page`);

        console.log('\n🖼️ Enhanced Example 4: Extract Images with Metadata');
        const images = await client.extractImages(true);
        const imageData = JSON.parse(images.content[0].text);
        console.log(`✅ Extracted ${imageData.count} images with metadata`);

        console.log('\n🔍 Enhanced Example 5: SEO Analysis');
        const seo = await client.analyzeSEO();
        const seoData = JSON.parse(seo.content[0].text);
        console.log('✅ SEO Analysis:', {
            title: seoData.seo.title,
            headings: Object.keys(seoData.seo.headings || {}).length
        });

        console.log('\n🔒 Enhanced Example 6: Security Analysis');
        const security = await client.checkSecurity();
        const securityData = JSON.parse(security.content[0].text);
        console.log('✅ Security Analysis:', {
            ssl_valid: securityData.security.ssl.valid,
            vulnerabilities: securityData.security.vulnerabilities.length
        });

        console.log('\n💾 Enhanced Example 7: Local Storage Management');
        await client.setLocalStorage('test_key', 'enhanced_value');
        const storageResult = await client.getLocalStorage('test_key');
        const storageData = JSON.parse(storageResult.content[0].text);
        console.log('✅ Local Storage:', storageData);

        console.log('\n📤 Enhanced Example 8: JSON POST Request');
        const postResult = await client.postJSON('https://httpbin.org/post', {
            message: 'Enhanced MCP Test',
            timestamp: new Date().toISOString(),
            capabilities: ['performance', 'security', 'seo']
        });
        console.log('✅ JSON POST result:', JSON.parse(postResult.content[0].text));

        console.log('\n📋 Enhanced Example 9: Form Analysis');
        await client.navigateEnhanced('https://httpbin.org/forms/post');
        const forms = await client.extractForms();
        const formData = JSON.parse(forms.content[0].text);
        console.log(`✅ Found ${formData.count} forms on the page`);

        console.log('\n⚡ Enhanced Example 10: Advanced Waiting');
        const waitResult = await client.waitAdvanced({
            selector: 'form',
            timeout: 5
        });
        console.log('✅ Advanced wait completed:', JSON.parse(waitResult.content[0].text));

        console.log('\n🧹 Enhanced Example 11: Session Cleanup');
        const clearResult = await client.clearSession();
        console.log('✅ Session cleared:', JSON.parse(clearResult.content[0].text));

        console.log('\n🎉 All Enhanced Examples Completed Successfully!');
        console.log('\n🚀 Summary of Enhanced Capabilities:');
        console.log('   • Performance tracking and analysis');
        console.log('   • Comprehensive content extraction (links, images, forms)');
        console.log('   • SEO analysis and optimization insights');
        console.log('   • Security vulnerability scanning');
        console.log('   • Local storage management');
        console.log('   • JSON API interaction');
        console.log('   • Advanced waiting conditions');
        console.log('   • Session management and cleanup');

    } catch (error) {
        console.error('❌ Enhanced example error:', error.message);
    } finally {
        await client.stop();
    }
}

// Performance benchmark example
async function runPerformanceBenchmark() {
    const client = new EnhancedMCPChromeClient();
    
    try {
        await client.start();
        
        console.log('\n⚡ Performance Benchmark: Testing Enhanced Capabilities');
        
        const testSites = [
            'https://example.com',
            'https://httpbin.org',
            'https://google.com'
        ];

        for (const site of testSites) {
            console.log(`\n📊 Benchmarking: ${site}`);
            
            const start = Date.now();
            const result = await client.navigateEnhanced(site, { trackPerformance: true });
            const navTime = Date.now() - start;
            
            const performance = await client.getPerformanceMetrics();
            const perfData = JSON.parse(performance.content[0].text);
            
            console.log(`   Navigation time: ${navTime}ms`);
            console.log(`   DOM Content Loaded: ${perfData.performance.dom_content_loaded}ms`);
            console.log(`   Load Complete: ${perfData.performance.load_complete}ms`);
            console.log(`   Network Requests: ${perfData.performance.network_requests}`);
        }
        
    } catch (error) {
        console.error('❌ Benchmark error:', error.message);
    } finally {
        await client.stop();
    }
}

// Handle command line arguments
const args = process.argv.slice(2);

// Handle graceful shutdown
process.on('SIGINT', () => {
    console.log('\n👋 Shutting down enhanced client gracefully...');
    process.exit(0);
});

// Run examples based on arguments
if (require.main === module) {
    if (args.includes('--benchmark')) {
        runPerformanceBenchmark().catch(console.error);
    } else {
        runEnhancedExamples().catch(console.error);
    }
}

module.exports = EnhancedMCPChromeClient;