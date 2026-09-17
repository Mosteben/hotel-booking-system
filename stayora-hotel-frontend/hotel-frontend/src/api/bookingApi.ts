// Unlike every other module, the booking handlers return the raw Booking
// (or Booking[]) on success and { error: string } on failure - there is no
// {success, message, data} envelope here. See internal/booking/handler.

import { apiClient } from "./client";
import type {
  AdminBookingSummary,
  Booking,
  BookingStatus,
  CreateBookingRequest,
} from "@/types/booking";

// internal/booking/handler CreateBooking/UpdateBooking bind the JSON body
// straight into model.Booking, whose CheckIn/CheckOut are time.Time - Go's
// encoding/json requires RFC3339 for those ("2026-09-20" alone fails to
// unmarshal with "cannot parse \"\" as \"T\"", confirmed by reproducing the
// exact unmarshal call). HTML <input type="date"> only ever gives bare
// "YYYY-MM-DD", so every date sent to this endpoint must be normalized to
// midnight UTC RFC3339 here before the request leaves the browser.
function toRFC3339Midnight(dateOnly: string): string {
  return `${dateOnly}T00:00:00Z`;
}

export async function createBooking(
  payload: CreateBookingRequest
): Promise<Booking> {
  const { data } = await apiClient.post<Booking>("/bookings", {
    ...payload,
    check_in: toRFC3339Midnight(payload.check_in),
    check_out: toRFC3339Midnight(payload.check_out),
  });
  return data;
}

export async function getMyBookings(): Promise<Booking[]> {
  const { data } = await apiClient.get<Booking[]>("/bookings/my");
  return data;
}

export async function getBookingByID(id: number | string): Promise<Booking> {
  const { data } = await apiClient.get<Booking>(`/bookings/${id}`);
  return data;
}

export async function cancelBooking(
  id: number | string
): Promise<{ message: string }> {
  const { data } = await apiClient.delete<{ message: string }>(
    `/bookings/${id}`
  );
  return data;
}

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.
// Unlike every other booking endpoint, this one returns each booking with
// its user/hotel/room already resolved server-side (see AdminBookingSummary)
// - the frontend never has to guess relationships from bare IDs.

export async function getAllBookings(): Promise<AdminBookingSummary[]> {
  const { data } = await apiClient.get<AdminBookingSummary[]>("/bookings");
  return data;
}

// internal/booking/service UpdateBookingStatus only accepts "confirmed" or
// "cancelled" - it rejects anything else (including "pending").
export async function updateBookingStatus(
  id: number | string,
  status: Extract<BookingStatus, "confirmed" | "cancelled">
): Promise<Booking> {
  const { data } = await apiClient.patch<Booking>(
    `/bookings/${id}/status`,
    { status }
  );
  return data;
}
