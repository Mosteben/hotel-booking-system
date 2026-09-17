import { apiClient } from "./client";
import type { ApiResponse } from "@/types/api";
import type { AdminUserSummary } from "@/types/user";

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.
export async function getAllUsers(): Promise<ApiResponse<AdminUserSummary[]>> {
  const { data } = await apiClient.get<ApiResponse<AdminUserSummary[]>>("/users");
  return data;
}
