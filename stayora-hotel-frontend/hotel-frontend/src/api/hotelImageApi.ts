import { apiClient } from "./client";
import type { ApiResponse, ApiSuccessNoData } from "@/types/api";
import type { HotelImage } from "@/types/hotel";

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.
// The upload actually reaches the backend, which validates the files and
// pushes them to real object storage - nothing here is mocked.

export async function uploadHotelImages(
  hotelId: number | string,
  files: File[]
): Promise<ApiResponse<HotelImage[]>> {
  const form = new FormData();
  files.forEach((file) => form.append("images", file));

  const { data } = await apiClient.post<ApiResponse<HotelImage[]>>(
    `/hotels/${hotelId}/images`,
    form,
    { headers: { "Content-Type": "multipart/form-data" } }
  );
  return data;
}

export async function deleteHotelImage(
  hotelId: number | string,
  imageId: number
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.delete<
    ApiSuccessNoData | { success: false; message: string }
  >(`/hotels/${hotelId}/images/${imageId}`);
  return data;
}

export async function setMainHotelImage(
  hotelId: number | string,
  imageId: number
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.patch<
    ApiSuccessNoData | { success: false; message: string }
  >(`/hotels/${hotelId}/images/${imageId}/main`);
  return data;
}
