export interface Session {
  access_token: string;
  expires_at: number;
  expires_in: number;
  refresh_token: string;
  token_type: string;
  user: SessionUser;
}

export interface SessionUser {
  app_metadata: string;
  aud: string;
  confirmed_at: string;
  created_at: string;
  email: string;
  email_confirmed_at: string;
  id: string;
  last_sign_in_at: string;
  role: string;
  updated_at: string;
  user_metadata: SessionUserMetadata;
}

export interface User {
  id: string;
  email: string;
  name: string;
}

export interface SessionUserMetadata {
  email: string;
  email_verified: boolean;
  name: string;
  sub: string;
}

export interface LoginParams {
  email: string;
  password: string;
}

export interface AuthContextProps {
  loggedIn: boolean;
  user: User | null;
  loading: boolean;
  login: (params: LoginParams) => Promise<void>;
  logout: () => Promise<void>;
}
