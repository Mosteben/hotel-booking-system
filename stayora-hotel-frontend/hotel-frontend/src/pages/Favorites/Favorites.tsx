import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Heart } from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { PageHeading } from "@/components/common/PageHeading";
import { EmptyState } from "@/components/common/EmptyState";
import { ErrorState } from "@/components/common/ErrorState";
import {
  FavoriteHotelCard,
  FavoriteHotelCardSkeleton,
} from "@/components/HotelCard/FavoriteHotelCard";
import { getMyFavorites } from "@/api/favoriteApi";
import { getHotelByID } from "@/api/hotelApi";
import { extractErrorMessage } from "@/api/client";
import type { Hotel } from "@/types/hotel";

type Status = "loading" | "success" | "empty" | "error";

export function Favorites() {
  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [hotels, setHotels] = useState<Hotel[]>([]);

  async function load() {
    setStatus("loading");
    try {
      const favRes = await getMyFavorites();
      if (!favRes.success) {
        setErrorMessage(favRes.message || "Couldn't load your favorites.");
        setStatus("error");
        return;
      }

      if (favRes.data.length === 0) {
        setHotels([]);
        setStatus("empty");
        return;
      }

      const results = await Promise.all(
        favRes.data.map(async (fav) => {
          try {
            const hotelRes = await getHotelByID(fav.hotel_id);
            return hotelRes.success ? hotelRes.data : null;
          } catch {
            return null;
          }
        })
      );

      const validHotels = results.filter((h): h is Hotel => h !== null);
      if (validHotels.length === 0) {
        setHotels([]);
        setStatus("empty");
      } else {
        setHotels(validHotels);
        setStatus("success");
      }
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }
  }

  useEffect(() => {
    load();
  }, []);

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-5xl px-4 py-10 sm:px-6 sm:py-14">
        <PageHeading
          eyebrow="Saved stays"
          title="Your favorites"
          subtitle="Hotels you've saved for later."
        />

        {status === "loading" && (
          <div className="mt-10 flex flex-col gap-6">
            {Array.from({ length: 3 }).map((_, i) => (
              <FavoriteHotelCardSkeleton key={i} />
            ))}
          </div>
        )}

        {status === "success" && (
          <div className="mt-10 flex flex-col gap-6">
            {hotels.map((hotel, i) => (
              <div
                key={hotel.id}
                className="card-enter"
                style={{ animationDelay: `${Math.min(i, 8) * 60}ms` }}
              >
                <FavoriteHotelCard
                  hotel={hotel}
                  onFavoriteChange={(hotelId, isFavorite) => {
                    if (isFavorite) return;
                    setHotels((prev) => {
                      const next = prev.filter((h) => h.id !== hotelId);
                      if (next.length === 0) setStatus("empty");
                      return next;
                    });
                  }}
                />
              </div>
            ))}
          </div>
        )}

        {status === "empty" && (
          <div className="mt-10">
            <EmptyState
              icon={Heart}
              title="No favorites yet"
              description="Tap the heart icon on any hotel to save it here."
              action={
                <Link
                  to="/"
                  className="mt-2 rounded-pill bg-teal px-6 py-3 text-sm font-semibold text-white hover:bg-teal-dark"
                >
                  Browse hotels
                </Link>
              }
            />
          </div>
        )}

        {status === "error" && (
          <div className="mt-10">
            <ErrorState
              title="Couldn't load your favorites"
              message={errorMessage}
              onRetry={load}
            />
          </div>
        )}
      </div>

      <Footer />
    </div>
  );
}
