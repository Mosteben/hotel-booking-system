const STYLES: Record<string, string> = {
  pending: "bg-amber-50 text-amber-600",
  confirmed: "bg-emerald-50 text-emerald-600",
  paid: "bg-emerald-50 text-emerald-600",
  cancelled: "bg-red-50 text-red-500",
  failed: "bg-red-50 text-red-500",
  refunded: "bg-bg text-muted",
};

export function StatusBadge({ status }: { status: string }) {
  const style = STYLES[status] ?? "bg-bg text-muted";
  return (
    <span
      className={`inline-flex items-center rounded-pill px-2.5 py-1 text-xs font-semibold capitalize ${style}`}
    >
      {status}
    </span>
  );
}
