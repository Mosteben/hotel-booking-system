import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Building2, MapPin, Star } from "lucide-react";
import type { Hotel } from "@/types/hotel";
import { getHotelImageUrl } from "@/utils/hotelImages";
import { getHotelStartingPrice } from "@/api/roomApi";
import { FavoriteButton } from "@/components/common/FavoriteButton";

// A wide, horizontal card for the Favorites list - deliberately different
// from the compact grid HotelCard used on Home/Search. Favorites is a
// short, personal list rather than a browse-many grid, so each entry gets
// the full container width and room for its image, details, and actions to
// breathe, instead of being squeezed into a 3-up grid cell.
export function FavoriteHotelCard({
  hotel,
  onFavoriteChange,
}: {
  hotel: Hotel;
  onFavoriteChange: (hotelId: number, isFavorite: boolean) => void;
}) {
  const imageUrl = getHotelImageUrl(hotel);
  const [startingPrice, setStartingPrice] = useState<number | null>(null);

  useEffect(() => {
    let cancelled = false;
    getHotelStartingPrice(hotel.id).then((price) => {
      if (!cancelled) setStartingPrice(price);
    });
    return () => {
      cancelled = true;
    };
  }, [hotel.id]);

  return (
    <article className="flex flex-col gap-6 rounded-card border border-line bg-white p-7 shadow-[0_8px_24px_rgba(16,24,40,0.06)] transition-shadow hover:shadow-[0_16px_36px_rgba(16,24,40,0.1)] sm:p-8 lg:grid lg:grid-cols-[220px_1fr_180px] lg:items-stretch lg:gap-8">
      {/* Zone 1: image */}
      <Link
        to={`/hotels/${hotel.id}`}
        className="relative block h-48 w-full shrink-0 overflow-hidden rounded-[16px] bg-cream lg:h-auto"
      >
        {imageUrl ? (
          <img src={imageUrl} alt={hotel.name} className="h-full w-full object-cover" />
        ) : (
          <div className="flex h-full w-full flex-col items-center justify-center gap-2">
            <span className="flex h-11 w-11 items-center justify-center rounded-full bg-white text-teal shadow-sm">
              <Building2 size={20} />
            </span>
            <p className="text-xs font-medium text-muted">No photo yet</p>
          </div>
        )}
        {hotel.stars > 0 && (
          <div className="absolute left-3 top-3 flex items-center gap-1 rounded-pill bg-white/95 px-2.5 py-1 text-xs font-semibold text-ink shadow-sm">
            <Star size={12} className="fill-teal text-teal" />
            {hotel.stars}
          </div>
        )}
      </Link>

      {/* Zone 2: hotel information, separated by a divider on desktop */}
      <div className="flex flex-col justify-center gap-4 lg:border-l lg:border-line lg:px-8">
        <Link to={`/hotels/${hotel.id}`}>
          <h2 className="font-display text-lg font-semibold text-ink">{hotel.name}</h2>
          <p className="mt-1.5 flex items-center gap-1.5 text-sm text-muted">
            <MapPin size={14} className="shrink-0" />
            {hotel.city}
            {hotel.country ? `, ${hotel.country}` : ""}
          </p>
          {hotel.description && (
            <p className="mt-3 max-w-lg text-sm leading-relaxed text-muted line-clamp-2">
              {hotel.description}
            </p>
          )}
        </Link>
      </div>

      {/* Zone 3: price + actions, separated by its own divider on desktop */}
      <div className="flex flex-row items-center justify-between gap-4 border-t border-line pt-6 lg:flex-col lg:items-end lg:justify-center lg:border-l lg:border-t-0 lg:pl-8 lg:pt-0 lg:text-right">
        <div>
          <p className="text-xs font-medium uppercase tracking-wide text-muted">
            {startingPrice !== null ? "From" : "Price"}
          </p>
          {startingPrice !== null ? (
            <p className="mt-1 font-display text-2xl font-bold text-ink">
              ${startingPrice.toFixed(0)}
              <span className="text-sm font-normal text-muted"> / night</span>
            </p>
          ) : (
            <p className="mt-1 text-sm font-semibold text-muted">On request</p>
          )}
        </div>

        <div className="flex items-center gap-2.5 lg:w-full lg:flex-col">
          <Link
            to={`/hotels/${hotel.id}`}
            className="rounded-pill bg-teal px-5 py-2.5 text-center text-sm font-semibold text-white transition-colors hover:bg-teal-dark lg:w-full"
          >
            View details
          </Link>
          <FavoriteButton
            hotelId={hotel.id}
            isFavorite
            onChange={(next) => onFavoriteChange(hotel.id, next)}
          />
        </div>
      </div>
    </article>
  );
}

export function FavoriteHotelCardSkeleton() {
  return (
    <div className="flex flex-col gap-6 rounded-card border border-line bg-white p-7 sm:p-8 lg:grid lg:grid-cols-[220px_1fr_180px] lg:gap-8">
      <div className="h-48 w-full animate-pulse rounded-[16px] bg-line lg:h-auto" />
      <div className="flex flex-col justify-center gap-3 lg:border-l lg:border-line lg:px-8">
        <div className="h-5 w-48 animate-pulse rounded bg-line" />
        <div className="h-3 w-32 animate-pulse rounded bg-line" />
        <div className="h-3 w-full max-w-xs animate-pulse rounded bg-line" />
      </div>
      <div className="flex flex-col justify-center gap-4 lg:border-l lg:border-line lg:pl-8">
        <div className="h-8 w-20 animate-pulse rounded bg-line" />
        <div className="h-10 w-full animate-pulse rounded-pill bg-line" />
      </div>
    </div>
  );
}
