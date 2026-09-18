import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Building2, MapPin, Star } from "lucide-react";
import type { Hotel } from "@/types/hotel";
import { getHotelImageUrl } from "@/utils/hotelImages";
import { getHotelStartingPrice } from "@/api/roomApi";
import { getMyFavorites } from "@/api/favoriteApi";
import { useAuth } from "@/context/AuthContext";
import { FavoriteButton } from "@/components/common/FavoriteButton";
import { Skeleton } from "@/components/common/Skeleton";

function ImagePlaceholder() {
  return (
    <div className="flex h-full w-full flex-col items-center justify-center gap-2 bg-cream">
      <span className="flex h-11 w-11 items-center justify-center rounded-full bg-white text-teal shadow-sm">
        <Building2 size={20} />
      </span>
      <p className="text-xs font-medium text-muted">No photo yet</p>
    </div>
  );
}

// Fetches each hotel's own real starting price (min room price) client-side,
// since the Hotel model itself has no price field - price only lives on
// Room. `initialFavorite` lets a page that already knows the favorite state
// (e.g. Favorites list) skip the extra lookup.
export function HotelCard({
  hotel,
  initialFavorite = false,
  favoriteKnown = false,
  showFavorite = true,
  onFavoriteChange,
}: {
  hotel: Hotel;
  initialFavorite?: boolean;
  // Set when the parent already fetched /favorites for the whole list, so
  // this card doesn't issue its own redundant lookup per card.
  favoriteKnown?: boolean;
  showFavorite?: boolean;
  // Lets a page (e.g. Favorites) react to a confirmed backend change, such
  // as removing the card from view the moment it's actually unfavorited.
  onFavoriteChange?: (hotelId: number, isFavorite: boolean) => void;
}) {
  const { isAuthenticated } = useAuth();
  const imageUrl = getHotelImageUrl(hotel);

  const [startingPrice, setStartingPrice] = useState<number | null>(null);
  const [isFavorite, setIsFavorite] = useState(initialFavorite);

  useEffect(() => {
    let cancelled = false;
    getHotelStartingPrice(hotel.id).then((price) => {
      if (!cancelled) setStartingPrice(price);
    });
    return () => {
      cancelled = true;
    };
  }, [hotel.id]);

  useEffect(() => {
    if (favoriteKnown || !isAuthenticated || !showFavorite) return;
    let cancelled = false;
    getMyFavorites()
      .then((res) => {
        if (!cancelled && res.success) {
          setIsFavorite(res.data.some((f) => f.hotel_id === hotel.id));
        }
      })
      .catch(() => {});
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hotel.id, isAuthenticated]);

  return (
    <article className="group relative flex flex-col overflow-hidden rounded-card bg-white shadow-[var(--shadow-card)] transition-all duration-300 hover:-translate-y-1 hover:shadow-[var(--shadow-hover)]">
      <Link
        to={`/hotels/${hotel.id}`}
        aria-label={`View ${hotel.name}`}
        className="absolute inset-0 z-0"
      />

      <div className="relative h-52 w-full overflow-hidden">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={hotel.name}
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
        ) : (
          <ImagePlaceholder />
        )}

        {hotel.stars > 0 && (
          <div className="pointer-events-none absolute left-3 top-3 flex items-center gap-1 rounded-pill bg-white/95 px-2.5 py-1 text-xs font-semibold text-ink shadow-sm">
            <Star size={12} className="fill-gold text-gold" />
            {hotel.stars}
          </div>
        )}

        {showFavorite && (
          <div className="absolute right-3 top-3 z-10">
            <FavoriteButton
              hotelId={hotel.id}
              isFavorite={isFavorite}
              onChange={(next) => {
                setIsFavorite(next);
                onFavoriteChange?.(hotel.id, next);
              }}
              size="compact"
            />
          </div>
        )}
      </div>

      <div className="flex flex-1 flex-col gap-2.5 p-6">
        <div className="flex items-start justify-between gap-3">
          <h3 className="font-display text-base font-semibold leading-snug text-ink">
            {hotel.name}
          </h3>
        </div>

        <p className="flex items-center gap-1.5 text-sm text-muted">
          <MapPin size={14} className="shrink-0" />
          <span className="truncate">
            {hotel.city}
            {hotel.country ? `, ${hotel.country}` : ""}
          </span>
        </p>

        {hotel.description && (
          <p className="line-clamp-2 text-sm leading-relaxed text-muted">
            {hotel.description}
          </p>
        )}

        <div className="mt-auto flex items-end justify-between pt-2">
          <div>
            {startingPrice !== null ? (
              <p className="font-display text-lg font-bold text-gold">
                ${startingPrice.toFixed(0)}
                <span className="text-xs font-normal text-muted"> / night</span>
              </p>
            ) : (
              <p className="text-xs font-medium text-muted">Price on request</p>
            )}
          </div>
          {/* Purely decorative label, not a second interactive element -
              the whole card is already one real, accessible <Link> (see
              the full-card overlay above, aria-labelled "View {hotel}").
              This used to be `relative z-10`, which put it visually
              *above* that Link in stacking order; with no click handler
              or href of its own, it silently absorbed every click aimed
              at this pill instead of letting the real Link underneath
              receive it - that was the actual cause of "View details"
              being unclickable. pointer-events-none lets clicks (and
              hover, for the group-hover color swap below) pass straight
              through to the real Link. */}
          <span className="pointer-events-none rounded-pill bg-teal/10 px-3.5 py-2 text-xs font-semibold text-teal transition-colors group-hover:bg-teal group-hover:text-white">
            View details
          </span>
        </div>
      </div>
    </article>
  );
}

export function HotelCardSkeleton() {
  return (
    <div className="flex flex-col overflow-hidden rounded-card bg-white shadow-[var(--shadow-card)]">
      <Skeleton className="h-52 w-full rounded-none" />
      <div className="flex flex-col gap-3 p-6">
        <Skeleton className="h-4 w-2/3" />
        <Skeleton className="h-3 w-1/2" />
        <Skeleton className="h-3 w-full" />
        <div className="mt-2 flex items-center justify-between">
          <Skeleton className="h-5 w-16" />
          <Skeleton className="h-8 w-24 rounded-pill" />
        </div>
      </div>
    </div>
  );
}
