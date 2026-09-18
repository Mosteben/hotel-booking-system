// Low-level shimmer-block primitive. Card-specific skeletons (HotelCardSkeleton,
// FavoriteHotelCardSkeleton, booking-card skeletons, etc.) compose this
// instead of each hand-rolling their own loading-block divs, so every
// loading state in the app shares one visual rhythm. Uses the
// `skeleton-shimmer` keyframe (a moving highlight sweep, defined in
// index.css) rather than Tailwind's flat opacity-pulse, and is fully
// covered by the global prefers-reduced-motion rule in index.css.

export function Skeleton({ className = "" }: { className?: string }) {
  return (
    <div
      className={`skeleton-shimmer overflow-hidden rounded-[10px] ${className}`}
    />
  );
}
