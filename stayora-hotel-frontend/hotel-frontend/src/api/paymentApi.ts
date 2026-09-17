import { apiClient } from "./client";
import type { ApiResponse } from "@/types/api";
import type {
  AdminPaymentSummary,
  Payment,
  PaymentMethod,
  PaymentStatus,
} from "@/types/payment";

export async function createPayment(
  bookingId: number | string,
  paymentMethod: PaymentMethod
): Promise<ApiResponse<Payment>> {
  const { data } = await apiClient.post<ApiResponse<Payment>>(
    `/bookings/${bookingId}/payment`,
    { payment_method: paymentMethod }
  );
  return data;
}

export async function getMyPayments(): Promise<ApiResponse<Payment[]>> {
  const { data } = await apiClient.get<ApiResponse<Payment[]>>("/payments/my");
  return data;
}

export async function getPaymentByID(
  id: number | string
): Promise<ApiResponse<Payment>> {
  const { data } = await apiClient.get<ApiResponse<Payment>>(`/payments/${id}`);
  return data;
}

// Admin/Manager only - requires AuthMiddleware + RequireRoles on the backend.
// Unlike every other payment endpoint, this one returns each payment with
// its payer/booking/hotel/room already resolved server-side (see
// AdminPaymentSummary).

export async function getAllPayments(): Promise<ApiResponse<AdminPaymentSummary[]>> {
  const { data } = await apiClient.get<ApiResponse<AdminPaymentSummary[]>>("/payments");
  return data;
}

// internal/payment/service UpdatePaymentStatus only accepts "paid",
// "failed", or "refunded", and enforces its own transition rules
// (e.g. only a paid payment can be refunded) - invalid transitions come
// back as a normal API error, not a crash.
export async function updatePaymentStatus(
  id: number | string,
  status: Extract<PaymentStatus, "paid" | "failed" | "refunded">
): Promise<ApiResponse<Payment>> {
  const { data } = await apiClient.patch<ApiResponse<Payment>>(
    `/payments/${id}/status`,
    { status }
  );
  return data;
}
