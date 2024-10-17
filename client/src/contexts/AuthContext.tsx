import { createContext, ReactNode, useCallback, useEffect, useState } from "react";
import { ApiError } from "../api/error.ts";

const SESSION_STORE_KEY = "session";

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

const userFromSession = (session: Session): User | null => {
  if (!session) return null;
  return {
    id: session.user.id,
    email: session.user.email,
    name: session.user.user_metadata.name ?? "unknown",
  };
};

const AuthContext = createContext<AuthContextProps | undefined>(undefined);

const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [loading, setLoading] = useState<boolean>(false);
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

  const logout = useCallback(async () => {
    setLoading(true);

    try {
      const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/logout`;
      const resp = await fetch(endpoint, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session?.access_token}`,
          "Content-Type": "application/json",
        },
      });

      if (!resp.ok) {
        const errorBody = await resp.json();
        const message =
          errorBody?.message || `Error making request. Status: ${resp.status}, ${resp.statusText}`;
        throw new ApiError(resp.status, resp.statusText, message);
      }
    } catch (error) {
      console.error("logout error", error);
    } finally {
      setSession(null);
      setLoading(false);
    }
  }, [session?.access_token]);

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

  useEffect(() => {
    const refreshToken = async () => {
      console.log("Refreshing token");
      const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/refresh`;

      try {
        const resp = await fetch(endpoint, {
          method: "POST",
          headers: {
            Authorization: `Bearer ${session?.access_token}`,
            "Content-Type": "application/json",
          },
        });

        if (!resp.ok) {
          const errorBody = await resp.json();
          const message =
            errorBody?.message || `Error making request. Status: ${resp.status}, ${resp.statusText}`;
          throw new ApiError(resp.status, resp.statusText, message);
        }

        const newSession: Session = await resp.json();
        if (!newSession?.access_token || !newSession?.refresh_token || !newSession?.user) {
          throw new Error("Invalid session data");
        }

        setSession(newSession);
      } catch (error) {
        console.error("refresh error", error);
      }
    };

    if (!session) return;
    const TOKEN_REFRESH_THRESHOLD = 5 * 60; // Refresh 5 minutes before expiration
    const currentTime = Math.floor(Date.now() / 1000);

    // Log out if the token is already expired.
    if (session.expires_at <= currentTime) {
      logout().catch((err) => console.error(err));
      return;
    }

    const refreshTime = session.expires_at - TOKEN_REFRESH_THRESHOLD;
    const timeoutId = setTimeout(refreshToken, (refreshTime - currentTime) * 1000);

    return () => clearTimeout(timeoutId);
  }, [session, logout]);

  const login = async (credentials: LoginParams) => {
    setLoading(true);
    const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/login`;

    try {
      const resp = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(credentials),
      });

      if (!resp.ok) {
        const errorBody = await resp.json();
        const message =
          errorBody?.message || `Error making request. Status: ${resp.status}, ${resp.statusText}`;
        throw new ApiError(resp.status, resp.statusText, message);
      }

      const session: Session = await resp.json();
      if (!session?.access_token || !session?.refresh_token || !session?.user) {
        throw new Error("Invalid session data");
      }

      setSession(session);
    } catch (error) {
      console.error("login error", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthContext.Provider value={{ loggedIn, user, loading, login, logout }}>{children}</AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };
