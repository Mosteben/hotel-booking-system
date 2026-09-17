// Mirrors the real /payments response fields exactly (internal/payment/model/payment.go).

export type PaymentStatus = "pending" | "paid" | "failed" | "refunded";
export type PaymentMethod = "cash" | "card";

export interface Payment {
  id: number;
  booking_id: number;
  user_id: string;
  amount: number;
  payment_method: PaymentMethod;
  status: PaymentStatus;
  transaction_id: string | null;
  created_at: string;
  updated_at: string;
}

// Mirrors GET /payments exactly (internal/payment/service AdminPaymentSummary).
// Admin/Manager only - payer/booking/hotel/room resolved server-side. Any
// of them can be null when the reference doesn't resolve to a real record
// - never fabricated.
export interface AdminPaymentUser {
  id: string;
  name: string;
  email: string;
  phone: string;
}

export interface AdminPaymentBooking {
  id: number;
  check_in: string;
  check_out: string;
  status: string;
}

export interface AdminPaymentHotel {
  id: number;
  name: string;
}

export interface AdminPaymentRoom {
  id: number;
  room_number: string;
  type: string;
}

export interface AdminPaymentSummary {
  id: number;
  booking_id: number;
  amount: number;
  payment_method: PaymentMethod;
  status: PaymentStatus;
  transaction_id: string | null;
  created_at: string;
  updated_at: string;
  user: AdminPaymentUser | null;
  booking: AdminPaymentBooking | null;
  hotel: AdminPaymentHotel | null;
  room: AdminPaymentRoom | null;
}
