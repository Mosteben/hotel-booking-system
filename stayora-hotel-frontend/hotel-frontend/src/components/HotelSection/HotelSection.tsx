import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { ArrowRight, BuildingIcon, RefreshCw, TriangleAlert } from "lucide-react";
import { getHotels } from "@/api/hotelApi";
import { getMyFavorites } from "@/api/favoriteApi";
import { extractErrorMessage } from "@/api/client";
import { HotelCard, HotelCardSkeleton } from "@/components/HotelCard/HotelCard";
import { useAuth } from "@/context/AuthContext";
import type { Hotel } from "@/types/hotel";

type Status = "loading" | "success" | "empty" | "error";

export function HotelSection() {
  const { isAuthenticated } = useAuth();
  const [hotels, setHotels] = useState<Hotel[]>([]);
  const [favoriteIds, setFavoriteIds] = useState<Set<number> | null>(null);
  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");

  async function load() {
    setStatus("loading");
    try {
      const response = await getHotels();
      if (response.success) {
        if (response.data.length === 0) {
          setStatus("empty");
        } else {
          setHotels(response.data);
          setStatus("success");
        }
      } else {
        setErrorMessage(response.message || "Couldn't load hotels right now.");
        setStatus("error");
      }
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }

    if (isAuthenticated) {
      try {
        const favRes = await getMyFavorites();
        if (favRes.success) {
          setFavoriteIds(new Set(favRes.data.map((f) => f.hotel_id)));
        }
      } catch {
        setFavoriteIds(new Set());
      }
    } else {
      setFavoriteIds(new Set());
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <section id="hotels" className="px-4 py-16 sm:px-6 sm:py-20">
      <div className="mx-auto max-w-6xl">
        <div className="mb-10 flex flex-col gap-2 sm:mb-12 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <p className="text-sm font-semibold tracking-wide text-teal">
              Available now
            </p>
            <h2 className="mt-1.5 font-display text-2xl font-bold text-ink sm:text-3xl">
              Stays ready to book
            </h2>
          </div>
          <Link
            to="/search"
            className="inline-flex items-center gap-1.5 text-sm font-semibold text-teal transition-colors hover:text-teal-dark"
          >
            Browse all hotels
            <ArrowRight size={15} />
          </Link>
        </div>

        {status === "loading" && (
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 sm:gap-y-10 lg:grid-cols-3 lg:gap-10">
            {Array.from({ length: 6 }).map((_, i) => (
              <HotelCardSkeleton key={i} />
            ))}
          </div>
        )}

        {status === "success" && (
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 sm:gap-y-10 lg:grid-cols-3 lg:gap-10">
            {hotels.slice(0, 6).map((hotel) => (
              <HotelCard
                key={hotel.id}
                hotel={hotel}
                favoriteKnown={favoriteIds !== null}
                initialFavorite={favoriteIds?.has(hotel.id) ?? false}
              />
            ))}
          </div>
        )}

        {status === "empty" && (
          <div className="flex flex-col items-center gap-3 rounded-panel border border-dashed border-line bg-white px-6 py-16 text-center">
            <span className="flex h-12 w-12 items-center justify-center rounded-full bg-bg text-teal">
              <BuildingIcon size={22} />
            </span>
            <p className="font-display text-lg font-semibold text-ink">
              No stays listed yet
            </p>
            <p className="max-w-sm text-sm text-muted">
              Hotels will appear here as soon as they're added to the
              platform.
            </p>
          </div>
        )}

        {status === "error" && (
          <div className="flex flex-col items-center gap-3 rounded-panel border border-line bg-white px-6 py-16 text-center">
            <span className="flex h-12 w-12 items-center justify-center rounded-full bg-red-50 text-red-500">
              <TriangleAlert size={22} />
            </span>
            <p className="font-display text-lg font-semibold text-ink">
              Couldn't load hotels
            </p>
            <p className="max-w-sm text-sm text-muted">{errorMessage}</p>
            <button
              onClick={load}
              className="mt-2 inline-flex items-center gap-2 rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark cursor-pointer"
            >
              <RefreshCw size={15} />
              Try again
            </button>
          </div>
        )}
      </div>
    </section>
  );
}
