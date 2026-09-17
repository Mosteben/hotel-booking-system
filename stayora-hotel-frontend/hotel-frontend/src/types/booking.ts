// Mirrors the real /bookings response fields exactly (internal/booking/model/booking.go).
//
// IMPORTANT: unlike every other module, the booking handlers do NOT use the
// {success, message, data} envelope - they return the raw Booking (or
// Booking[]) on success, and { error: string } on failure. See bookingApi.ts.

export type BookingStatus = "pending" | "confirmed" | "cancelled";

export interface Booking {
  id: number;
  user_id: string;
  room_id: number;
  check_in: string;
  check_out: string;
  guests: number;
  total_price: number;
  status: BookingStatus;
  // Snapshot of the room/hotel as they existed when this booking was made
  // (or last changed). Present on every booking created after this field
  // was added; empty/zero on older bookings that predate it - never a
  // sign the room/hotel itself is empty. Kept optional here since bookings
  // fetched before this field existed on the frontend build should still
  // type-check.
  hotel_id?: number;
  hotel_name?: string;
  room_number?: string;
  room_type?: string;
  price_per_night?: number;
  created_at: string;
  updated_at: string;
}

export interface CreateBookingRequest {
  room_id: number;
  check_in: string;
  check_out: string;
  guests: number;
}

export interface RoomAvailability {
  room_id: number;
  check_in: string;
  check_out: string;
  available: boolean;
}

// Mirrors GET /bookings exactly (internal/booking/service AdminBookingSummary).
// Admin/Manager only - the user/hotel/room relationships are resolved
// server-side. Any of them can be null: the booking's user_id/room_id may
// not match an existing record (e.g. stale seed data) - that's real, never
// papered over with a fabricated name.
export interface AdminBookingUser {
  id: string;
  name: string;
  email: string;
  phone: string;
}

export interface AdminBookingHotel {
  id: number;
  name: string;
}

export interface AdminBookingRoom {
  id: number;
  room_number: string;
  type: string;
}

export interface AdminBookingSummary {
  id: number;
  status: BookingStatus;
  check_in: string;
  check_out: string;
  guests: number;
  total_price: number;
  created_at: string;
  updated_at: string;
  user: AdminBookingUser | null;
  hotel: AdminBookingHotel | null;
  room: AdminBookingRoom | null;
  // Snapshot captured at booking time - unlike hotel/room above, this is
  // never null by nature, but on bookings from before this field existed
  // it's simply empty/zero. Use it as the fallback label when hotel/room
  // are null (deleted since booking), not as a substitute for them
  // otherwise.
  hotel_id_snapshot: number;
  hotel_name_snapshot: string;
  room_number_snapshot: string;
  room_type_snapshot: string;
  price_per_night_snapshot: number;
}
