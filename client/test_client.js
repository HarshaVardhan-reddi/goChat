const WebSocket = require('ws');
const axios = require('axios');
const readline = require('readline');

const API_BASE_URL = 'http://localhost:3000/api/v1';
const WS_BASE_URL = 'ws://localhost:3000/api/v1/ws';

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

const question = (query) => new Promise((resolve) => rl.question(query, resolve));

async function main() {
    console.log("--- goChat Test Client ---");
    const mode = await question("Choose mode: (1) Signup (2) Login: ");

    let user;
    let token;

    if (mode === '1') {
        const email = await question("Email: ");
        const password = await question("Password: ");
        const firstName = await question("First Name: ");
        const nickname = await question("Nickname: ");

        try {
            const response = await axios.post(`${API_BASE_URL}/users/signup`, {
                email,
                password,
                first_name: firstName,
                nickname
            });
            console.log("Signup successful!");
            // After signup, we need to login to get the token
            const loginResponse = await axios.post(`${API_BASE_URL}/users/login`, { email, password });
            user = loginResponse.data.user.User;
            token = loginResponse.data.user.AuthToken;
        } catch (err) {
            console.error("Signup/Login failed:", err.response?.data || err.message);
            process.exit(1);
        }
    } else {
        const email = await question("Email: ");
        const password = await question("Password: ");

        try {
            const response = await axios.post(`${API_BASE_URL}/users/login`, { email, password });
            console.log("Login successful!");
            user = response.data.user.User;
            token = response.data.user.AuthToken;
        } catch (err) {
            console.error("Login failed:", err.response?.data || err.message);
            process.exit(1);
        }
    }

    console.log(`\n--- Session Started ---`);
    console.log(`User: ${user.first_name}`);
    console.log(`ID: ${user.ID}`);
    console.log(`-----------------------\n`);
    
    const ws = new WebSocket(`${WS_BASE_URL}/startchat`, {
        headers: {
            'Authorization': token
        }
    });

    ws.on('open', () => {
        console.log("WebSocket connection established!");
        startMessaging(ws, user);
    });

    ws.on('message', (data) => {
        // The server sends back raw text for the message content based on messaging_service.go
        console.log(`\n[Incoming Message]: ${data.toString()}`);
    });

    ws.on('error', (err) => {
        console.error("WebSocket error:", err.message);
    });

    ws.on('close', () => {
        console.log("WebSocket connection closed.");
        process.exit(0);
    });
}

async function startMessaging(ws, currentUser) {
    while (true) {
        const targetId = await question("To User ID (or 'exit'): ");
        if (targetId.toLowerCase() === 'exit') {
            ws.close();
            break;
        }

        const messageText = await question("Message: ");
        
        const chatMsg = {
            sessionId: "test-session",
            timestamp: Math.floor(Date.now() / 1000),
            from: { id: currentUser.ID },
            to: { id: parseInt(targetId) },
            message: {
                text: messageText,
                attachments: []
            },
            token: ""
        };

        ws.send(JSON.stringify(chatMsg));
        console.log(`[Sent to ${targetId}]: ${messageText}`);
    }
}

main();
