import { apiClient } from "./client";
import type { ApiResponse, ApiSuccessNoData } from "@/types/api";
import type { Room, RoomRequest } from "@/types/room";
import type { RoomAvailability } from "@/types/booking";

export async function getRoomsByHotelID(
  hotelId: number | string
): Promise<ApiResponse<Room[]>> {
  const { data } = await apiClient.get<ApiResponse<Room[]>>(
    `/rooms/hotel/${hotelId}`
  );
  return data;
}

// The Hotel model has no price field - price only exists per Room. This
// derives a real "from $X/night" figure from that hotel's actual rooms
// instead of inventing one. Returns null (not 0) when there's genuinely no
// priced room yet, so callers can render "Price on request" instead of a
// fabricated number.
export async function getHotelStartingPrice(
  hotelId: number | string
): Promise<number | null> {
  try {
    const response = await getRoomsByHotelID(hotelId);
    if (!response.success || response.data.length === 0) return null;
    return Math.min(...response.data.map((room) => room.price_per_night));
  } catch {
    return null;
  }
}

export async function getRoomByID(
  roomId: number | string
): Promise<ApiResponse<Room>> {
  const { data } = await apiClient.get<ApiResponse<Room>>(`/rooms/${roomId}`);
  return data;
}

// Requires auth (see routes.go) even though it's a read - the backend
// registers this under AuthMiddleware.
export async function checkRoomAvailability(
  roomId: number | string,
  checkIn: string,
  checkOut: string
): Promise<RoomAvailability> {
  const { data } = await apiClient.get<RoomAvailability>(
    `/rooms/${roomId}/availability`,
    { params: { check_in: checkIn, check_out: checkOut } }
  );
  return data;
}

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.

export async function createRoom(
  hotelId: number | string,
  room: RoomRequest
): Promise<ApiResponse<Room>> {
  const { data } = await apiClient.post<ApiResponse<Room>>(
    `/rooms/hotel/${hotelId}`,
    room
  );
  return data;
}

export async function updateRoom(
  id: number | string,
  room: RoomRequest
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.put<
    ApiSuccessNoData | { success: false; message: string }
  >(`/rooms/${id}`, room);
  return data;
}

// Admin only - DeleteRoom is not gated for manager on the backend.
export async function deleteRoom(
  id: number | string
): Promise<ApiSuccessNoData | { success: false; message: string }> {
  const { data } = await apiClient.delete<
    ApiSuccessNoData | { success: false; message: string }
  >(`/rooms/${id}`);
  return data;
}
