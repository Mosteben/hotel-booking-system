import { useMemo } from "react";
import { Link } from "react-router-dom";
import { MapPin } from "lucide-react";
import { getHotelImageUrl } from "@/utils/hotelImages";
import type { Hotel } from "@/types/hotel";

// Derived entirely from the already-fetched hotel list (grouped by real
// `city` values, counted for real) - there's no "destination"/"category"
// concept in the API, so nothing here is invented. Renders nothing if
// there isn't enough distinct city data to make the section meaningful.
export function PopularDestinations({ hotels }: { hotels: Hotel[] }) {
  const destinations = useMemo(() => {
    const byCity = new Map<string, { count: number; hotel: Hotel }>();

    for (const hotel of hotels) {
      const city = hotel.city?.trim();
      if (!city) continue;

      const existing = byCity.get(city);
      if (existing) {
        existing.count += 1;
        if (!getHotelImageUrl(existing.hotel) && getHotelImageUrl(hotel)) {
          existing.hotel = hotel;
        }
      } else {
        byCity.set(city, { count: 1, hotel });
      }
    }

    return Array.from(byCity.entries())
      .map(([city, { count, hotel }]) => ({ city, count, hotel }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 4);
  }, [hotels]);

  if (destinations.length < 2) return null;

  return (
    <section className="px-4 py-16 sm:px-6 sm:py-20">
      <div className="mx-auto max-w-6xl">
        <div className="mb-10 sm:mb-12">
          <p className="text-sm font-semibold tracking-wide text-teal">
            Where to next
          </p>
          <h2 className="mt-1.5 font-display text-2xl font-bold text-ink sm:text-3xl">
            Popular destinations
          </h2>
        </div>

        <div className="grid grid-cols-2 gap-5 lg:grid-cols-4">
          {destinations.map(({ city, count, hotel }, i) => {
            const imageUrl = getHotelImageUrl(hotel);
            return (
              <Link
                key={city}
                to={`/search?city=${encodeURIComponent(city)}`}
                className="card-enter group relative flex h-52 flex-col justify-end overflow-hidden rounded-card bg-ink shadow-[var(--shadow-card)] transition-shadow duration-300 hover:shadow-[var(--shadow-hover)]"
                style={{ animationDelay: `${i * 70}ms` }}
              >
                {imageUrl && (
                  <img
                    src={imageUrl}
                    alt=""
                    aria-hidden="true"
                    className="absolute inset-0 h-full w-full object-cover opacity-70 transition-transform duration-500 group-hover:scale-105"
                  />
                )}
                <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-ink/90 via-ink/20 to-transparent" />
                <div className="relative p-5">
                  <p className="flex items-center gap-1.5 font-display text-base font-semibold text-white">
                    <MapPin size={14} className="shrink-0" />
                    {city}
                  </p>
                  <p className="mt-0.5 text-xs text-white/80">
                    {count} {count === 1 ? "stay" : "stays"}
                  </p>
                </div>
              </Link>
            );
          })}
        </div>
      </div>
    </section>
  );
}
