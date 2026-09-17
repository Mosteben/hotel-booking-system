// Mirrors the real /favorites response fields exactly (internal/favorite/model/favorite.go).

export interface Favorite {
  id: number;
  user_id: string;
  hotel_id: number;
  created_at: string;
}
