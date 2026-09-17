import axios from "axios";

const baseURL = import.meta.env.VITE_API_URL;

if (!baseURL) {
  // Fail loudly in dev rather than silently hitting a wrong host.
  // eslint-disable-next-line no-console
  console.error(
    "VITE_API_URL is not set. Create a .env file (see .env.example)."
  );
}

export const apiClient = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

const TOKEN_KEY = "nilestay_token";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

apiClient.interceptors.request.use((config) => {
  const token = getToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Normalizes any error thrown by axios into a single friendly string,
// without ever leaking raw technical/network details to the UI.
//
// The booking module is the one exception to the app's {success, message,
// data} envelope - its handlers return { error: string } instead of
// { message: string } on failure (see internal/booking/handler). Both keys
// are checked here so booking/room-availability errors surface their real
// backend reason instead of falling through to a generic message.
export function extractErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data;
    const backendMessage =
      typeof data?.message === "string" && data.message.trim()
        ? data.message
        : typeof data?.error === "string" && data.error.trim()
          ? data.error
          : null;
    if (backendMessage) {
      return backendMessage;
    }
    if (error.code === "ERR_NETWORK" || !error.response) {
      return "We couldn't reach the server. Please check your connection and try again.";
    }
    if (error.response.status >= 500) {
      return "Something went wrong on our end. Please try again shortly.";
    }
  }
  return "Something went wrong. Please try again.";
}
