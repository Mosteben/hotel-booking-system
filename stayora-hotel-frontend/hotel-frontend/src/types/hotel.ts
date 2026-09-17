// Mirrors the real GET /hotels response fields exactly.

import type { Room } from "@/types/room";

export interface HotelImage {
  id: number;
  hotel_id: number;
  url: string;
  is_main: boolean;
  created_at: string;
}

export interface Hotel {
  id: number;
  name: string;
  description: string;
  address: string;
  city: string;
  country: string;
  phone: string;
  email: string;
  stars: number;
  rooms: unknown | null;
  images: HotelImage[] | null;
  image_url: string | null;
  created_at: string;
  updated_at: string;
}

// POST /hotels and PUT /hotels/:id request body.
export interface HotelRequest {
  name: string;
  description?: string;
  address: string;
  city: string;
  country: string;
  phone?: string;
  email?: string;
  stars: number;
}

// GET /hotels/:id/details preloads the hotel's rooms.
export interface HotelDetails extends Omit<Hotel, "rooms"> {
  rooms: Room[] | null;
}

export interface HotelSearchFilters {
  name?: string;
  city?: string;
  country?: string;
  stars?: number;
  min_price?: number;
  max_price?: number;
  room_type?: string;
  min_capacity?: number;
}
