// Success/error use the approved brand hex values (--color-success/-error);
// pending stays a neutral amber (not part of the brand palette, but the
// universal "waiting" signal) rather than being forced into orange, which
// would make it read as a brand accent/CTA instead of a status.
const STYLES: Record<string, string> = {
  pending: "bg-amber-50 text-amber-600",
  confirmed: "bg-success-soft text-success",
  paid: "bg-success-soft text-success",
  cancelled: "bg-error-soft text-error",
  failed: "bg-error-soft text-error",
  refunded: "bg-bg text-muted",
};

export function StatusBadge({ status }: { status: string }) {
  const style = STYLES[status] ?? "bg-bg text-muted";
  return (
    <span
      className={`inline-flex items-center rounded-pill px-2.5 py-1 text-xs font-semibold capitalize transition-colors duration-200 ${style}`}
    >
      {status}
    </span>
  );
}
