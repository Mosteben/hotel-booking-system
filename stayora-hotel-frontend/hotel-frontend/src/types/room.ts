// Mirrors the real GET /rooms/... response fields exactly (internal/room/model/room.go).

export type RoomStatus = "available" | "occupied" | "maintenance";

export interface RoomImage {
  id: number;
  room_id: number;
  url: string;
  is_main: boolean;
  created_at: string;
}

export interface Room {
  id: number;
  hotel_id: number;
  room_number: string;
  type: string;
  description: string;
  price_per_night: number;
  capacity: number;
  status: RoomStatus;
  images: RoomImage[] | null;
  image_url: string | null;
  created_at: string;
  updated_at: string;
}

// POST /hotels/:hotel_id/rooms and PUT /rooms/:id request body.
export interface RoomRequest {
  room_number: string;
  type: string;
  description?: string;
  price_per_night: number;
  capacity: number;
  status?: RoomStatus;
}
