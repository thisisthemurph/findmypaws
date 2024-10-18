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

export interface LogInParams {
  email: string;
  password: string;
}

export interface SignUpParams extends LogInParams {
  username: string;
}

export interface AuthContextProps {
  user: User | null;
  session: Session | null;
  signup: (params: SignUpParams) => Promise<void>;
  login: (params: LogInParams) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: () => boolean;
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

  const sessionHasExpired = useCallback((): boolean => {
    const currentTime = Math.floor(Date.now() / 1000);
    return !session || session.expires_at <= currentTime;
  }, [session]);

  const logout = useCallback(async () => {
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
    } finally {
      setSession(null);
    }
  }, [session?.access_token]);

  useEffect(() => {
    if (session) {
      setUser(userFromSession(session));
      localStorage.setItem(SESSION_STORE_KEY, JSON.stringify(session));
    } else {
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
            "X-Refresh-Token": `Bearer ${session?.refresh_token}`,
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
    if (sessionHasExpired()) {
      logout().catch((err) => console.error(err));
      return;
    }

    const refreshTime = session.expires_at - TOKEN_REFRESH_THRESHOLD;
    const timeoutId = setTimeout(refreshToken, (refreshTime - currentTime) * 1000);

    return () => clearTimeout(timeoutId);
  }, [session, logout, sessionHasExpired]);

  const login = async (credentials: LogInParams) => {
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
      const message =
        errorBody?.message || `Error making request. Status: ${resp.status}, ${resp.statusText}`;
      throw new ApiError(resp.status, resp.statusText, message);
    }

    const session: Session = await resp.json();
    if (!session?.access_token || !session?.refresh_token || !session?.user) {
      throw new Error("Invalid session data");
    }

    setSession(session);
  };

  const signup = async (values: SignUpParams) => {
    const endpoint = `${import.meta.env.VITE_API_BASE_URL}/auth/signup`;

    const resp = await fetch(endpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(values),
    });

    if (!resp.ok) {
      const errorBody = await resp.json();
      const message =
        errorBody?.message || `Error making request. Status: ${resp.status}, ${resp.statusText}`;
      throw new ApiError(resp.status, resp.statusText, message);
    }
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: () => !sessionHasExpired(),
        user,
        session,
        signup,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export { AuthContext, AuthProvider };
