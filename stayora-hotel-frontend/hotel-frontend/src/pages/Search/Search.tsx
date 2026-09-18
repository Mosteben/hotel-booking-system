import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { SearchX, SlidersHorizontal } from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { PageHeading } from "@/components/common/PageHeading";
import { EmptyState } from "@/components/common/EmptyState";
import { ErrorState } from "@/components/common/ErrorState";
import { SearchBar } from "@/components/SearchBar/SearchBar";
import { HotelCard, HotelCardSkeleton } from "@/components/HotelCard/HotelCard";
import { searchHotels } from "@/api/hotelApi";
import { getMyFavorites } from "@/api/favoriteApi";
import { extractErrorMessage } from "@/api/client";
import { useAuth } from "@/context/AuthContext";
import type { Hotel, HotelSearchFilters } from "@/types/hotel";

type Status = "loading" | "success" | "empty" | "error";
type SortOption = "relevance" | "stars_desc";

export function Search() {
  const { isAuthenticated } = useAuth();
  const [searchParams] = useSearchParams();
  const [hotels, setHotels] = useState<Hotel[]>([]);
  const [favoriteIds, setFavoriteIds] = useState<Set<number> | null>(null);
  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");

  const [minStars, setMinStars] = useState(0);
  const [maxPrice, setMaxPrice] = useState("");
  const [sort, setSort] = useState<SortOption>("relevance");

  const city = searchParams.get("city") || "";
  const minCapacity = searchParams.get("min_capacity");
  const checkIn = searchParams.get("check_in");
  const checkOut = searchParams.get("check_out");

  async function load() {
    setStatus("loading");
    try {
      const filters: HotelSearchFilters = {};
      if (city) filters.city = city;
      if (minCapacity) filters.min_capacity = Number(minCapacity);
      if (maxPrice) filters.max_price = Number(maxPrice);
      if (minStars > 0) filters.stars = minStars;

      const response = await searchHotels(filters);
      if (response.success) {
        if (response.data.length === 0) {
          setStatus("empty");
        } else {
          setHotels(response.data);
          setStatus("success");
        }
      } else {
        setErrorMessage(response.message || "Couldn't search hotels right now.");
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
  }, [searchParams, minStars, maxPrice]);

  // Sorting reorders the real results already returned by the backend - it
  // never adds, removes, or invents data.
  const sortedHotels = useMemo(() => {
    if (sort === "stars_desc") {
      return [...hotels].sort((a, b) => b.stars - a.stars);
    }
    return hotels;
  }, [hotels, sort]);

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="bg-cream px-4 pb-8 pt-8 sm:px-6">
        <div className="mx-auto max-w-6xl">
          <SearchBar />
        </div>
      </div>

      <div className="page-fade-in mx-auto max-w-6xl px-4 py-10 sm:px-6">
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
          <PageHeading
            eyebrow="Hotel search"
            title={city ? `Stays in ${city}` : "All stays"}
            subtitle={
              checkIn || checkOut
                ? `${checkIn ? new Date(checkIn).toLocaleDateString() : "Any date"} to ${
                    checkOut ? new Date(checkOut).toLocaleDateString() : "Any date"
                  }`
                : undefined
            }
          />

          {status === "success" && (
            <label className="flex items-center gap-2 text-sm text-muted">
              Sort by
              <select
                value={sort}
                onChange={(e) => setSort(e.target.value as SortOption)}
                className="rounded-pill border border-line bg-white px-3.5 py-2 text-sm font-medium text-ink outline-none focus:border-teal cursor-pointer"
              >
                <option value="relevance">Relevance</option>
                <option value="stars_desc">Highest rated</option>
              </select>
            </label>
          )}
        </div>

        {/* Filters */}
        <div className="mt-8 flex flex-wrap items-center gap-3 rounded-card border border-line bg-white p-5">
          <span className="flex items-center gap-1.5 text-xs font-semibold uppercase tracking-wider text-muted">
            <SlidersHorizontal size={14} />
            Filters
          </span>

          <div className="flex items-center gap-1.5">
            {[0, 3, 4, 5].map((value) => (
              <button
                key={value}
                onClick={() => setMinStars(value)}
                className={`rounded-pill px-3.5 py-1.5 text-xs font-semibold transition-colors cursor-pointer ${
                  minStars === value
                    ? "bg-teal text-white"
                    : "bg-bg text-muted hover:text-ink"
                }`}
              >
                {value === 0 ? "Any rating" : `${value}+ stars`}
              </button>
            ))}
          </div>

          <div className="flex items-center gap-2">
            <label className="text-xs font-medium text-muted">
              Max price/night
            </label>
            <input
              type="number"
              min={0}
              placeholder="Any"
              value={maxPrice}
              onChange={(e) => setMaxPrice(e.target.value)}
              className="w-24 rounded-pill border border-line bg-white px-3 py-1.5 text-xs font-medium text-ink outline-none focus:border-teal"
            />
          </div>
        </div>

        {status === "loading" && (
          <div className="mt-10 grid grid-cols-1 gap-8 sm:grid-cols-2 sm:gap-y-10 lg:grid-cols-3 lg:gap-10">
            {Array.from({ length: 6 }).map((_, i) => (
              <HotelCardSkeleton key={i} />
            ))}
          </div>
        )}

        {status === "success" && (
          <div className="mt-10 grid grid-cols-1 gap-8 sm:grid-cols-2 sm:gap-y-10 lg:grid-cols-3 lg:gap-10">
            {sortedHotels.map((hotel, i) => (
              <div
                key={hotel.id}
                className="card-enter"
                style={{ animationDelay: `${Math.min(i, 11) * 50}ms` }}
              >
                <HotelCard
                  hotel={hotel}
                  favoriteKnown={favoriteIds !== null}
                  initialFavorite={favoriteIds?.has(hotel.id) ?? false}
                />
              </div>
            ))}
          </div>
        )}

        {status === "empty" && (
          <div className="mt-8">
            <EmptyState
              icon={SearchX}
              title="No stays match your search"
              description="Try a different destination, fewer guests, or a wider price range."
            />
          </div>
        )}

        {status === "error" && (
          <div className="mt-8">
            <ErrorState
              title="Couldn't search hotels"
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
