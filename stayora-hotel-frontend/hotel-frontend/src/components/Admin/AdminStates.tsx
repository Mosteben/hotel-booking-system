import { RefreshCw, TriangleAlert } from "lucide-react";
import type { LucideIcon } from "lucide-react";
import { Skeleton } from "@/components/common/Skeleton";

export function AdminLoadingRows({ count = 5 }: { count?: number }) {
  return (
    <div className="flex flex-col gap-3">
      {Array.from({ length: count }).map((_, i) => (
        <Skeleton key={i} className="h-16 w-full rounded-[14px]" />
      ))}
    </div>
  );
}

export function AdminErrorState({
  message,
  onRetry,
}: {
  message: string;
  onRetry: () => void;
}) {
  return (
    <div className="flex flex-col items-center gap-3 rounded-[14px] border border-dashed border-line bg-white px-6 py-16 text-center">
      <span className="flex h-11 w-11 items-center justify-center rounded-full bg-red-50 text-red-500">
        <TriangleAlert size={20} />
      </span>
      <p className="text-sm text-muted">{message}</p>
      <button
        onClick={onRetry}
        className="mt-1 flex items-center gap-2 rounded-pill bg-teal px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-teal-dark cursor-pointer"
      >
        <RefreshCw size={14} />
        Try again
      </button>
    </div>
  );
}

export function AdminEmptyState({
  icon: Icon,
  title,
  description,
}: {
  icon: LucideIcon;
  title: string;
  description?: string;
}) {
  return (
    <div className="flex flex-col items-center gap-2 rounded-[14px] border border-dashed border-line bg-white px-6 py-16 text-center">
      <span className="flex h-11 w-11 items-center justify-center rounded-full bg-bg text-muted">
        <Icon size={20} />
      </span>
      <p className="text-sm font-semibold text-ink">{title}</p>
      {description && <p className="max-w-sm text-sm text-muted">{description}</p>}
    </div>
  );
}
