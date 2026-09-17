import { apiClient } from "./client";
import type { ApiResponse, ApiSuccessNoData, ApiError } from "@/types/api";
import type {
  AuthenticatedUser,
  ChangePasswordRequest,
  LoginRequest,
  LoginResponseData,
  RegisterRequest,
  UpdateProfileRequest,
} from "@/types/auth";

export async function login(
  payload: LoginRequest
): Promise<ApiResponse<LoginResponseData>> {
  const { data } = await apiClient.post<ApiResponse<LoginResponseData>>(
    "/auth/login",
    payload
  );
  return data;
}

export async function register(
  payload: RegisterRequest
): Promise<{ success: boolean; message: string }> {
  const { data } = await apiClient.post<{ success: boolean; message: string }>(
    "/auth/register",
    payload
  );
  return data;
}

export async function getCurrentUser(): Promise<
  ApiResponse<AuthenticatedUser>
> {
  const { data } = await apiClient.get<ApiResponse<AuthenticatedUser>>(
    "/auth/me"
  );
  return data;
}

export async function updateProfile(
  payload: UpdateProfileRequest
): Promise<ApiSuccessNoData | ApiError> {
  const { data } = await apiClient.put<ApiSuccessNoData | ApiError>(
    "/auth/profile",
    payload
  );
  return data;
}

export async function changePassword(
  payload: ChangePasswordRequest
): Promise<ApiSuccessNoData | ApiError> {
  const { data } = await apiClient.put<ApiSuccessNoData | ApiError>(
    "/auth/change-password",
    payload
  );
  return data;
}
