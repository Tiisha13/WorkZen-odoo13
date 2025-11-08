"use client";

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
} from "react";
import { useRouter } from "next/navigation";
import { apiService } from "@/lib/api-service";
import type {
  User,
  Company,
  LoginRequest,
  SignupRequest,
  ChangePasswordRequest,
} from "@/lib/types";
import { toast } from "sonner";

interface AuthContextType {
  user: User | null;
  company: Company | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (credentials: LoginRequest) => Promise<void>;
  signup: (data: SignupRequest) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
  changePassword: (data: ChangePasswordRequest) => Promise<void>;
  hasRole: (roles: string | string[]) => boolean;
  isSuperAdmin: () => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [company, setCompany] = useState<Company | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  const isAuthenticated = !!user && apiService.isAuthenticated();

  const loadUser = useCallback(async () => {
    try {
      if (apiService.isAuthenticated()) {
        const userData = await apiService.getMe();
        setUser(userData);
        const companyData = apiService.getCompany();
        setCompany(companyData);
      } else {
        // No token, clear everything
        setUser(null);
        setCompany(null);
      }
    } catch (error) {
      console.error("Failed to load user:", error);
      // Clear auth on error
      apiService.logout();
      setUser(null);
      setCompany(null);

      // Only redirect if we're not already on login/signup pages
      if (typeof window !== "undefined") {
        const path = window.location.pathname;
        if (!path.includes("/login") && !path.includes("/signup")) {
          router.push("/login");
        }
      }
    } finally {
      setIsLoading(false);
    }
  }, [router]);

  useEffect(() => {
    loadUser();
  }, [loadUser]);

  const login = async (credentials: LoginRequest) => {
    try {
      setIsLoading(true);
      const response = await apiService.login(credentials);
      setUser(response.user);
      if (response.company) setCompany(response.company);
      toast.success("Login successful!");
      router.push("/dashboard");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Login failed");
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const signup = async (data: SignupRequest) => {
    try {
      setIsLoading(true);
      await apiService.signup(data);
      toast.success(
        "Signup successful! Please check your email for verification."
      );
      router.push("/login");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Signup failed");
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = useCallback(() => {
    // Clear API service storage
    apiService.logout();

    // Clear state
    setUser(null);
    setCompany(null);

    // Show toast
    toast.info("Logged out successfully");

    // Redirect to login
    router.push("/login");

    // Force reload to clear any cached state
    if (typeof window !== "undefined") {
      setTimeout(() => {
        window.location.href = "/login";
      }, 100);
    }
  }, [router]);

  const refreshUser = async () => {
    await loadUser();
  };

  const changePassword = async (data: ChangePasswordRequest) => {
    try {
      setIsLoading(true);
      await apiService.changePassword(data);
      toast.success("Password changed successfully!");
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : "Failed to change password"
      );
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const hasRole = (roles: string | string[]): boolean => {
    if (!user) return false;
    const roleArray = Array.isArray(roles) ? roles : [roles];
    return roleArray.includes(user.role);
  };

  const isSuperAdmin = (): boolean => {
    return user?.is_super_admin === true;
  };

  const value: AuthContextType = {
    user,
    company,
    isAuthenticated,
    isLoading,
    login,
    signup,
    logout,
    refreshUser,
    changePassword,
    hasRole,
    isSuperAdmin,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
