import { Link } from "react-router-dom";
import { ArrowRight, BuildingIcon } from "lucide-react";
import { HotelCard, HotelCardSkeleton } from "@/components/HotelCard/HotelCard";
import { EmptyState } from "@/components/common/EmptyState";
import { ErrorState } from "@/components/common/ErrorState";
import type { Hotel } from "@/types/hotel";
import type { HotelListStatus } from "@/pages/Home/Home";

// Data is fetched once by the parent (Home) and passed down - this section
// no longer issues its own /hotels request. See Home.tsx.
export function HotelSection({
  hotels,
  status,
  errorMessage,
  favoriteIds,
  onRetry,
}: {
  hotels: Hotel[];
  status: HotelListStatus;
  errorMessage: string;
  favoriteIds: Set<number> | null;
  onRetry: () => void;
}) {
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
            {hotels.slice(0, 6).map((hotel, i) => (
              <div
                key={hotel.id}
                className="card-enter"
                style={{ animationDelay: `${i * 60}ms` }}
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
          <EmptyState
            icon={BuildingIcon}
            title="No stays listed yet"
            description="Hotels will appear here as soon as they're added to the platform."
          />
        )}

        {status === "error" && (
          <ErrorState
            title="Couldn't load hotels"
            message={errorMessage}
            onRetry={onRetry}
          />
        )}
      </div>
    </section>
  );
}
