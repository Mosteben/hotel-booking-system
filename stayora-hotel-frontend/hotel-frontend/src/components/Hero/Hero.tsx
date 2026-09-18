import { SearchBar } from "@/components/SearchBar/SearchBar";
import { getHotelImageUrl } from "@/utils/hotelImages";
import type { Hotel } from "@/types/hotel";

export function Hero({ featuredHotel }: { featuredHotel: Hotel | null }) {
  const imageUrl = featuredHotel ? getHotelImageUrl(featuredHotel) : null;

  return (
    <section className="relative min-h-[440px] overflow-hidden px-4 pb-4 pt-10 sm:min-h-[540px] sm:px-6 sm:pt-14 lg:pt-16">
      {imageUrl ? (
        // A real hotel photo from the live catalog, not a stock image -
        // falls back to the plain gradient/blob treatment below when no
        // hotel has an image yet.
        //
        // Deliberately NOT z-index'd (no -z-10/z-0 here). This <section>
        // is `relative` but never gets an explicit z-index itself, so it
        // never becomes its own stacking context - a negative z-index
        // child would escape it and can end up painting behind the page's
        // own solid background (a real bug hit here previously, which is
        // why this image never appeared). Instead, this div is simply
        // first in DOM order and `absolute` (stacks at the default
        // level), and the text content below it is `relative` and later
        // in DOM order, so it naturally paints on top - no z-index needed
        // at all.
        <div className="absolute inset-0">
          <img
            src={imageUrl}
            alt=""
            aria-hidden="true"
            className="h-full w-full object-cover"
          />
          <div className="pointer-events-none absolute inset-0 bg-gradient-to-b from-ink/80 via-ink/55 to-bg" />
        </div>
      ) : (
        <>
          <div className="auth-blob pointer-events-none absolute -left-32 -top-32 h-96 w-96 rounded-full bg-teal/10 blur-3xl" />
          <div className="auth-blob auth-blob-delay pointer-events-none absolute -right-24 top-10 h-80 w-80 rounded-full bg-teal-light/10 blur-3xl" />
        </>
      )}

      <div className="page-fade-in relative mx-auto max-w-6xl py-6 sm:py-10">
        <div className="mx-auto max-w-2xl text-center">
          <p
            className={`text-sm font-semibold tracking-wide ${
              imageUrl ? "text-white/90" : "text-teal"
            }`}
          >
            Real stays, verified by NileStay
          </p>
          <h1
            className={`mt-2 font-display text-4xl font-bold leading-[1.1] tracking-tight sm:text-5xl lg:text-[3.4rem] ${
              imageUrl ? "text-white" : "text-ink"
            }`}
          >
            Stay somewhere
            <br />
            worth remembering
          </h1>
          <p
            className={`mx-auto mt-5 max-w-md text-base leading-relaxed sm:text-lg ${
              imageUrl ? "text-white/85" : "text-muted"
            }`}
          >
            Search, compare and book hotels with transparent pricing - every
            listing pulled straight from our live catalog.
          </p>
        </div>

        <div className="mt-10 sm:mt-12">
          <SearchBar />
        </div>
      </div>
    </section>
  );
}
