import axios from "axios";
import { jwtDecode } from "jwt-decode";

const baseURL = process.env.REACT_APP_API_URL || 'http://localhost:8081/api/v1';

const axiosInstance = axios.create({
    baseURL: baseURL,
});

axiosInstance.interceptors.request.use(async (req) => {
    const authTokens = localStorage.getItem('authTokens') ? JSON.parse(localStorage.getItem('authTokens')) : null;

    if (!authTokens) return req;

    const user = jwtDecode(authTokens.access_token);
    const isExpired = new Date(user.exp * 1000) < new Date();

    if (!isExpired) {
        req.headers.Authorization = `Bearer ${authTokens.access_token}`;
        return req;
    }

    try {
        const response = await axios.post(`${baseURL}/auth/refresh`, {
            refresh_token: authTokens.refresh_token,
        });

        localStorage.setItem('authTokens', JSON.stringify(response.data));
        req.headers.Authorization = `Bearer ${response.data.access_token}`;
        return req;
    } catch (error) {
        console.error("Could not refresh token:", error);
        localStorage.removeItem('authTokens');

        window.location.href = '/login'; // Redirect to login if token refresh fails
        return Promise.reject(error);
    }
});

export default axiosInstance;