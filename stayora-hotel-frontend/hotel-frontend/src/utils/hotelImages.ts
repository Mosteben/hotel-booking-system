// Centralized hotel image resolution.
//
// The backend computes image_url itself (the hotel's main image, or its
// first uploaded image if none is explicitly marked main) on every hotel
// endpoint. This never falls back to a stock/unrelated photo - a hotel
// with no images yet returns null here, and callers should render a
// neutral placeholder instead of a picture that isn't actually that hotel.

import type { Hotel } from "@/types/hotel";

export function getHotelImageUrl(hotel: Hotel): string | null {
  return hotel.image_url ?? null;
}
