// Mirrors the real /hotels/:id/reviews response fields exactly
// (internal/review/model/review.go). No update/delete - the backend
// doesn't support them yet.

export interface Review {
  id: number;
  hotel_id: number;
  user_id: string;
  rating: number;
  comment: string;
  created_at: string;
  updated_at: string;
}

export interface CreateReviewRequest {
  rating: number;
  comment: string;
}

// Mirrors GET /reviews exactly (internal/review/service AdminReviewSummary).
// Admin/Manager only - reviewer/hotel resolved server-side. Either can be
// null when the reference doesn't resolve to a real record - never
// fabricated.
export interface AdminReviewUser {
  id: string;
  name: string;
  email: string;
  phone: string;
}

export interface AdminReviewHotel {
  id: number;
  name: string;
}

export interface AdminReviewSummary {
  id: number;
  rating: number;
  comment: string;
  created_at: string;
  updated_at: string;
  user: AdminReviewUser | null;
  hotel: AdminReviewHotel | null;
}
