import { createContext, ReactNode, useEffect, useState } from "react";
import { ApiError } from "../api/error.ts";

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
  login: (params: LoginParams) => Promise<void>;
}

const userFromSession = (session: Session): User | null => {
  if (!session) return null;
  return {
    id: session.user.id,
    email: session.user.email,
    name: session.user.user_metadata.name ?? "unknown",
  }
}

const SESSION_STORE_KEY = "session";

export const AuthContext = createContext<AuthContextProps | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [loggedIn, setLoggedIn] = useState<boolean>(false);
  const [session, setSession] = useState<Session | null>(() => {
    const storedSession = localStorage.getItem(SESSION_STORE_KEY);
    return storedSession ? JSON.parse(storedSession) : null;
  });
  const [user, setUser] = useState<User | null>(() => {
    const storedSession = localStorage.getItem(SESSION_STORE_KEY);
    if (!storedSession) return null;
    const session: Session = JSON.parse(storedSession);
    return userFromSession(session);
  });

  useEffect(() => {
    if (session) {
      setLoggedIn(true);
      setUser(userFromSession(session));
      localStorage.setItem(SESSION_STORE_KEY, JSON.stringify(session));
    } else {
      setLoggedIn(false);
      setUser(null);
      localStorage.removeItem(SESSION_STORE_KEY);
    }
  }, [session]);

  const login = async (credentials: LoginParams) => {
    const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/login`;
    const resp = await fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(credentials),
    });

    if (!resp.ok) {
      const errorBody = await resp.json();
      const message = errorBody?.message || `Error making request. Status: ${resp.status}, ${resp.statusText}`;
      throw new ApiError(resp.status, resp.statusText, message);
    }

    const session: Session = await resp.json();
    if (!session?.access_token || !session?.refresh_token || !session?.user) {
      throw new Error("Invalid session data");
    }

    setSession(session);
  }

  return (
    <AuthContext.Provider value={{ loggedIn, user, login }}>
      {children}
    </AuthContext.Provider>
  )
}