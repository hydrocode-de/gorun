export interface User {
    id: string;
    email: string;
    isAdmin: boolean;
}

export interface AuthResponse {
    user: User;
    access_token: string;
    expires_at: Date;
    refresh_token: string;
}

export const config = $state({
    apiServer: 'http://localhost:8080',
    refreshToken: '',
    auth: {} as AuthResponse
})