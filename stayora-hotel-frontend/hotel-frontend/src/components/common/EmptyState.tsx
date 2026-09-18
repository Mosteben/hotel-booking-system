import type { LucideIcon } from "lucide-react";

// Customer-side counterpart to the Admin area's AdminEmptyState - same
// visual language (dashed panel, icon badge, title, optional description),
// used by every list-type customer page (Search, HotelDetails rooms/
// reviews, MyBookings, Favorites, HotelSection) instead of each hand-
// rolling its own empty-state markup.
export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
}: {
  icon: LucideIcon;
  title: string;
  description?: string;
  action?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-panel border border-dashed border-line bg-white px-6 py-16 text-center">
      <span className="flex h-12 w-12 items-center justify-center rounded-full bg-bg text-teal">
        <Icon size={22} />
      </span>
      <p className="font-display text-lg font-semibold text-ink">{title}</p>
      {description && <p className="max-w-sm text-sm text-muted">{description}</p>}
      {action}
    </div>
  );
}
