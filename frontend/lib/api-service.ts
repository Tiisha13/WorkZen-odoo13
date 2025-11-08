import { API_ENDPOINTS } from "./config";
import type {
  LoginRequest,
  LoginResponse,
  SignupRequest,
  ChangePasswordRequest,
  ApiResponse,
  User,
  Company,
} from "./types";

class ApiService {
  private getHeaders(includeAuth = true): HeadersInit {
    const headers: HeadersInit = { "Content-Type": "application/json" };
    if (includeAuth) {
      const token = this.getToken();
      if (token) headers["Authorization"] = `Bearer ${token}`;
    }
    return headers;
  }

  private getToken(): string | null {
    return typeof window !== "undefined" ? localStorage.getItem("token") : null;
  }

  setToken(token: string): void {
    if (typeof window !== "undefined") localStorage.setItem("token", token);
  }

  removeToken(): void {
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      localStorage.removeItem("company");
    }
  }

  setUser(user: User): void {
    if (typeof window !== "undefined")
      localStorage.setItem("user", JSON.stringify(user));
  }

  getUser(): User | null {
    if (typeof window !== "undefined") {
      const user = localStorage.getItem("user");
      return user ? JSON.parse(user) : null;
    }
    return null;
  }

  setCompany(company: Company): void {
    if (typeof window !== "undefined")
      localStorage.setItem("company", JSON.stringify(company));
  }

  getCompany(): Company | null {
    if (typeof window !== "undefined") {
      const company = localStorage.getItem("company");
      return company ? JSON.parse(company) : null;
    }
    return null;
  }

  async handleResponse<T>(response: Response): Promise<T> {
    try {
      const contentType = response.headers.get("content-type");

      // Handle non-JSON responses
      if (!contentType || !contentType.includes("application/json")) {
        if (!response.ok) {
          throw new Error(
            `Server error: ${response.status} ${response.statusText}`
          );
        }
        throw new Error("Invalid response format from server");
      }

      const data = await response.json();

      if (!response.ok) {
        // Handle 401 Unauthorized
        if (response.status === 401) {
          this.removeToken();
          if (typeof window !== "undefined") {
            window.location.href = "/login";
          }
          throw new Error("Session expired. Please login again.");
        }

        // Handle 404 Not Found
        if (response.status === 404) {
          throw new Error(data.message || "Resource not found");
        }

        // Handle other errors
        throw new Error(
          data.message ||
            data.error ||
            data.msg ||
            response.statusText ||
            "Request failed"
        );
      }

      return data;
    } catch (error) {
      if (error instanceof SyntaxError) {
        throw new Error(`Invalid JSON response: ${response.statusText}`);
      }
      if (error instanceof TypeError && error.message.includes("fetch")) {
        throw new Error("Network error. Please check your connection.");
      }
      throw error;
    }
  }

  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await fetch(API_ENDPOINTS.LOGIN, {
      method: "POST",
      headers: this.getHeaders(false),
      body: JSON.stringify(credentials),
    });
    const data = await this.handleResponse<ApiResponse<LoginResponse>>(
      response
    );
    if (data.success && data.data) {
      this.setToken(data.data.token);
      this.setUser(data.data.user);
      if (data.data.company) this.setCompany(data.data.company);
      return data.data;
    }
    throw new Error(data.message || "Login failed");
  }

  async signup(data: SignupRequest): Promise<ApiResponse> {
    const response = await fetch(API_ENDPOINTS.SIGNUP, {
      method: "POST",
      headers: this.getHeaders(false),
      body: JSON.stringify(data),
    });
    return this.handleResponse<ApiResponse>(response);
  }

  async getMe(): Promise<User> {
    const response = await fetch(API_ENDPOINTS.ME, {
      method: "GET",
      headers: this.getHeaders(),
    });
    const data = await this.handleResponse<ApiResponse<User>>(response);
    if (data.success && data.data) {
      this.setUser(data.data);
      return data.data;
    }
    throw new Error(data.message || "Failed to fetch user data");
  }

  async changePassword(data: ChangePasswordRequest): Promise<ApiResponse> {
    const response = await fetch(API_ENDPOINTS.CHANGE_PASSWORD, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });
    return this.handleResponse<ApiResponse>(response);
  }

  async verifyEmail(token: string): Promise<ApiResponse> {
    const response = await fetch(
      `${API_ENDPOINTS.VERIFY_EMAIL}?token=${token}`,
      {
        method: "GET",
        headers: this.getHeaders(false),
      }
    );
    return this.handleResponse<ApiResponse>(response);
  }

  async resendVerification(email: string): Promise<ApiResponse> {
    const response = await fetch(API_ENDPOINTS.RESEND_VERIFICATION, {
      method: "POST",
      headers: this.getHeaders(false),
      body: JSON.stringify({ email }),
    });
    return this.handleResponse<ApiResponse>(response);
  }

  logout(): void {
    this.removeToken();
  }

  async get<T>(url: string): Promise<T> {
    const response = await fetch(url, {
      method: "GET",
      headers: this.getHeaders(),
    });
    return this.handleResponse<T>(response);
  }

  async post<T>(url: string, data: unknown): Promise<T> {
    const response = await fetch(url, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });
    return this.handleResponse<T>(response);
  }

  async put<T>(url: string, data: unknown): Promise<T> {
    const response = await fetch(url, {
      method: "PUT",
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });
    return this.handleResponse<T>(response);
  }

  async patch<T>(url: string, data: unknown): Promise<T> {
    const response = await fetch(url, {
      method: "PATCH",
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });
    return this.handleResponse<T>(response);
  }

  async delete<T>(url: string): Promise<T> {
    const response = await fetch(url, {
      method: "DELETE",
      headers: this.getHeaders(),
    });
    return this.handleResponse<T>(response);
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }
}

export const apiService = new ApiService();
