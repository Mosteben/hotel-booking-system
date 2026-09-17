import { apiClient } from "./client";
import type { ApiResponse, ApiSuccessNoData } from "@/types/api";
import type { RoomImage } from "@/types/room";

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the
// backend. Mirrors hotelImageApi.ts exactly - same Cloudinary-backed
// storage, separate folder server-side.

export async function uploadRoomImages(
  roomId: number | string,
  files: File[]
): Promise<ApiResponse<RoomImage[]>> {
  const form = new FormData();
  files.forEach((file) => form.append("images", file));

  const { data } = await apiClient.post<ApiResponse<RoomImage[]>>(
    `/rooms/${roomId}/images`,
    form,
    { headers: { "Content-Type": "multipart/form-data" } }
  );
  return data;
}

export async function deleteRoomImage(
  roomId: number | string,
  imageId: number
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.delete<
    ApiSuccessNoData | { success: false; message: string }
  >(`/rooms/${roomId}/images/${imageId}`);
  return data;
}

export async function setMainRoomImage(
  roomId: number | string,
  imageId: number
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.patch<
    ApiSuccessNoData | { success: false; message: string }
  >(`/rooms/${roomId}/images/${imageId}/main`);
  return data;
}
