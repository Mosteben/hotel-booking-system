// Matches the backend's consistent response envelope.
// Success: { success: true, message: string, data?: T }
// Error:   { success: false, message: string }

export interface ApiSuccess<T> {
  success: true;
  message: string;
  data: T;
}

export interface ApiSuccessNoData {
  success: true;
  message: string;
}

export interface ApiError {
  success: false;
  message: string;
}

export type ApiResponse<T> = ApiSuccess<T> | ApiError;
