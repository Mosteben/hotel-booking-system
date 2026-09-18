import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Heart, Loader2, AlertCircle } from "lucide-react";
import { addFavorite, removeFavorite } from "@/api/favoriteApi";
import { useAuth } from "@/context/AuthContext";

export function FavoriteButton({
  hotelId,
  isFavorite,
  onChange,
  size = "default",
}: {
  hotelId: number;
  isFavorite: boolean;
  onChange: (next: boolean) => void;
  size?: "default" | "compact";
}) {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [hasError, setHasError] = useState(false);
  const [justToggled, setJustToggled] = useState(false);

  async function toggle(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();

    if (!isAuthenticated) {
      navigate("/login");
      return;
    }
    setIsSubmitting(true);
    setHasError(false);
    try {
      const response = isFavorite
        ? await removeFavorite(hotelId)
        : await addFavorite(hotelId);
      if (response.success) {
        onChange(!isFavorite);
        setJustToggled(true);
        setTimeout(() => setJustToggled(false), 300);
      } else {
        // The request completed but the backend rejected it - surface that
        // instead of pretending the toggle worked.
        setHasError(true);
        setTimeout(() => setHasError(false), 2500);
      }
    } catch {
      setHasError(true);
      setTimeout(() => setHasError(false), 2500);
    } finally {
      setIsSubmitting(false);
    }
  }

  const dimension = size === "compact" ? "h-9 w-9" : "h-11 w-11";
  const iconSize = size === "compact" ? 16 : 18;

  return (
    <button
      onClick={toggle}
      disabled={isSubmitting}
      aria-pressed={isFavorite}
      aria-label={isFavorite ? "Remove from favorites" : "Add to favorites"}
      title={hasError ? "Couldn't update favorites - try again" : undefined}
      className={`flex ${dimension} shrink-0 items-center justify-center rounded-full border shadow-sm transition-all duration-150 hover:scale-105 active:scale-95 cursor-pointer disabled:cursor-not-allowed disabled:hover:scale-100 ${
        justToggled ? "favorite-pop" : ""
      } ${
        hasError
          ? "border-amber-200 bg-amber-50 text-amber-600"
          : isFavorite
            ? "border-red-200 bg-red-50 text-red-500"
            : "border-line bg-white/95 text-muted hover:text-red-500"
      }`}
    >
      {isSubmitting ? (
        <Loader2 size={iconSize} className="animate-spin" />
      ) : hasError ? (
        <AlertCircle size={iconSize} />
      ) : (
        <Heart size={iconSize} className={isFavorite ? "fill-current" : ""} />
      )}
    </button>
  );
}
