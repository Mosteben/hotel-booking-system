import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { getCurrentUser } from "@/api/authApi";
import { clearToken, getToken, setToken as persistToken } from "@/api/client";
import type { AuthenticatedUser } from "@/types/auth";

interface AuthContextValue {
  user: AuthenticatedUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  loginWithToken: (token: string) => Promise<AuthenticatedUser | null>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthenticatedUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const fetchUser = useCallback(async (): Promise<AuthenticatedUser | null> => {
    try {
      const response = await getCurrentUser();
      if (response.success) {
        setUser(response.data);
        return response.data;
      }
      setUser(null);
      return null;
    } catch {
      setUser(null);
      clearToken();
      return null;
    }
  }, []);

  useEffect(() => {
    const existingToken = getToken();
    if (!existingToken) {
      setIsLoading(false);
      return;
    }
    fetchUser().finally(() => setIsLoading(false));
  }, [fetchUser]);

  const loginWithToken = useCallback(
    async (token: string) => {
      persistToken(token);
      return fetchUser();
    },
    [fetchUser]
  );

  const logout = useCallback(() => {
    clearToken();
    setUser(null);
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      isAuthenticated: Boolean(user),
      isLoading,
      loginWithToken,
      logout,
      refreshUser: async () => {
        await fetchUser();
      },
    }),
    [user, isLoading, loginWithToken, logout, fetchUser]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
