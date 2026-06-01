const axios = require('axios');

const API_BASE_URL = 'http://localhost:3000/api/v1';

async function signupAndGetToken() {
    // Randomize data to avoid unique constraint errors (email/nickname)
    const timestamp = Date.now();
    const userData = {
        email: `user_${timestamp}@example.com`,
        password: 'password123',
        first_name: 'Test',
        last_name: 'User',
        nickname: `tester_${timestamp}`
    };

    console.log("Attempting to signup with:", userData.email);

    try {
        // 1. Signup
        const signupResponse = await axios.post(`${API_BASE_URL}/users/signup`, userData);
        console.log("Signup Response Data:", JSON.stringify(signupResponse.data, null, 2));
        
        // Expected: { message: "...", user: { user: { ID: ... }, auth_token: "..." } }
        const signupData = signupResponse.data.user;
        const userObj = signupData.user;
        const authToken = signupData.auth_token;
        
        console.log("\n--- Success (from Signup)! ---");
        console.log("User ID:", userObj?.ID);
        console.log("Auth Token:", authToken);
        console.log("----------------\n");
        
        return authToken;
    } catch (err) {
        console.error("Error during authentication flow:");
        if (err.response) {
            console.error("Status:", err.response.status);
            console.error("Data:", err.response.data);
        } else {
            console.error(err.message);
        }
    }
}

signupAndGetToken();
