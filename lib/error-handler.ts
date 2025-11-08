import { toast } from "sonner";

export class AppError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public code?: string
  ) {
    super(message);
    this.name = "AppError";
  }
}

export function handleApiError(
  error: unknown,
  defaultMessage = "An error occurred"
): void {
  console.error("API Error:", error);

  if (error instanceof AppError) {
    toast.error(error.message);
    return;
  }

  if (error instanceof Error) {
    toast.error(error.message || defaultMessage);
    return;
  }

  if (typeof error === "string") {
    toast.error(error);
    return;
  }

  toast.error(defaultMessage);
}

export function getErrorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }
  if (typeof error === "string") {
    return error;
  }
  return "An unexpected error occurred";
}

export function isNetworkError(error: unknown): boolean {
  if (error instanceof Error) {
    return (
      error.message.includes("fetch") ||
      error.message.includes("network") ||
      error.message.includes("Network") ||
      error.message.includes("Failed to fetch")
    );
  }
  return false;
}

export function handleNetworkError(): void {
  toast.error("Network error. Please check your connection and try again.");
}

export function handleAuthError(): void {
  toast.error("Session expired. Please login again.");
  // Redirect to login will be handled by auth context
}
