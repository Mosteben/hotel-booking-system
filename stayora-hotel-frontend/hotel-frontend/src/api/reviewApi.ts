import { apiClient } from "./client";
import type { ApiResponse, ApiSuccessNoData } from "@/types/api";
import type {
  AdminReviewSummary,
  CreateReviewRequest,
  Review,
} from "@/types/review";

export async function getHotelReviews(
  hotelId: number | string
): Promise<ApiResponse<Review[]>> {
  const { data } = await apiClient.get<ApiResponse<Review[]>>(
    `/hotels/${hotelId}/reviews`
  );
  return data;
}

export async function createReview(
  hotelId: number | string,
  payload: CreateReviewRequest
): Promise<ApiResponse<Review>> {
  const { data } = await apiClient.post<ApiResponse<Review>>(
    `/hotels/${hotelId}/reviews`,
    payload
  );
  return data;
}

export async function getHotelAverageRating(
  hotelId: number | string
): Promise<ApiResponse<{ hotel_id: number; average_rating: number }>> {
  const { data } = await apiClient.get<
    ApiResponse<{ hotel_id: number; average_rating: number }>
  >(`/hotels/${hotelId}/rating`);
  return data;
}

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.
// Unlike the per-hotel customer endpoint above, this one returns each
// review with its reviewer/hotel already resolved server-side (see
// AdminReviewSummary).

export async function getAllReviews(): Promise<ApiResponse<AdminReviewSummary[]>> {
  const { data } = await apiClient.get<ApiResponse<AdminReviewSummary[]>>("/reviews");
  return data;
}

export async function deleteReview(
  id: number | string
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.delete<
    ApiSuccessNoData | { success: false; message: string }
  >(`/reviews/${id}`);
  return data;
}
