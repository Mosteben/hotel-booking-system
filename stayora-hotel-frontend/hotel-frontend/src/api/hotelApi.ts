import { apiClient } from "./client";
import type { ApiResponse, ApiSuccessNoData } from "@/types/api";
import type {
  Hotel,
  HotelDetails,
  HotelRequest,
  HotelSearchFilters,
} from "@/types/hotel";

export async function getHotels(): Promise<ApiResponse<Hotel[]>> {
  const { data } = await apiClient.get<ApiResponse<Hotel[]>>("/hotels");
  return data;
}

export async function searchHotels(
  filters: HotelSearchFilters
): Promise<ApiResponse<Hotel[]>> {
  const { data } = await apiClient.get<ApiResponse<Hotel[]>>("/hotels/search", {
    params: filters,
  });
  return data;
}

export async function getHotelByID(
  id: number | string
): Promise<ApiResponse<Hotel>> {
  const { data } = await apiClient.get<ApiResponse<Hotel>>(`/hotels/${id}`);
  return data;
}

export async function getHotelDetails(
  id: number | string
): Promise<ApiResponse<HotelDetails>> {
  const { data } = await apiClient.get<ApiResponse<HotelDetails>>(
    `/hotels/${id}/details`
  );
  return data;
}

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.

export async function createHotel(
  hotel: HotelRequest
): Promise<ApiResponse<Hotel>> {
  const { data } = await apiClient.post<ApiResponse<Hotel>>("/hotels", hotel);
  return data;
}

export async function updateHotel(
  id: number | string,
  hotel: HotelRequest
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.put<
    ApiSuccessNoData | { success: false; message: string }
  >(`/hotels/${id}`, hotel);
  return data;
}

export async function deleteHotel(
  id: number | string
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.delete<
    ApiSuccessNoData | { success: false; message: string }
  >(`/hotels/${id}`);
  return data;
}
