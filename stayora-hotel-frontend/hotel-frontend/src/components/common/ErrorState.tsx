import { RefreshCw, TriangleAlert } from "lucide-react";

// Customer-side counterpart to the Admin area's AdminErrorState. `onRetry`
// is optional since a few pages (Booking, Payment, BookingDetails) don't
// currently offer a retry action for their error state - kept that way
// rather than inventing one.
export function ErrorState({
  title = "Something went wrong",
  message,
  onRetry,
}: {
  title?: string;
  message: string;
  onRetry?: () => void;
}) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-panel border border-line bg-white px-6 py-16 text-center">
      <span className="flex h-12 w-12 items-center justify-center rounded-full bg-red-50 text-red-500">
        <TriangleAlert size={22} />
      </span>
      <p className="font-display text-lg font-semibold text-ink">{title}</p>
      <p className="max-w-sm text-sm text-muted">{message}</p>
      {onRetry && (
        <button
          onClick={onRetry}
          className="mt-2 inline-flex items-center gap-2 rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark cursor-pointer"
        >
          <RefreshCw size={15} />
          Try again
        </button>
      )}
    </div>
  );
}
