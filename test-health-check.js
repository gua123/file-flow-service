// 简单的健康检查测试脚本
const http = require('http');
const https = require('https');
const url = require('url');

// 测试 /health 接口
function testHealthCheck() {
    const options = {
        hostname: 'localhost',
        port: 8080,
        path: '/health',
        method: 'GET'
    };

    const req = http.request(options, (res) => {
        console.log(`状态码: ${res.statusCode}`);
        console.log(`响应头: ${JSON.stringify(res.headers)}`);
        
        res.on('data', (chunk) => {
            console.log(`响应数据: ${chunk}`);
        });
        
        res.on('end', () => {
            console.log('健康检查测试完成');
        });
    });

    req.on('error', (e) => {
        console.error(`请求出现问题: ${e.message}`);
    });

    req.end();
}

// 测试 /api/health 接口
function testAPIHealthCheck() {
    const options = {
        hostname: 'localhost',
        port: 8080,
        path: '/api/health',
        method: 'GET'
    };

    const req = http.request(options, (res) => {
        console.log(`API健康检查状态码: ${res.statusCode}`);
        console.log(`API健康检查响应头: ${JSON.stringify(res.headers)}`);
        
        res.on('data', (chunk) => {
            console.log(`API健康检查响应数据: ${chunk}`);
        });
        
        res.on('end', () => {
            console.log('API健康检查测试完成');
        });
    });

    req.on('error', (e) => {
        console.error(`API健康检查请求出现问题: ${e.message}`);
    });

    req.end();
}

console.log('开始测试健康检查接口...');
testHealthCheck();
testAPIHealthCheck();
