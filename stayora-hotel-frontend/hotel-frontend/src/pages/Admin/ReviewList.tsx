import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { Loader2, MessageSquareText, RefreshCw, Star, Trash2, UserRound } from "lucide-react";
import { useAdminData } from "@/context/AdminDataContext";
import { deleteReview } from "@/api/reviewApi";
import { extractErrorMessage } from "@/api/client";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import {
  AdminEmptyState,
  AdminErrorState,
  AdminLoadingRows,
} from "@/components/Admin/AdminStates";
import { ConfirmDialog } from "@/components/Admin/ConfirmDialog";
import { useConfirmDialog } from "@/hooks/useConfirmDialog";

export function ReviewList() {
  const { reviews, setReviews, isLoading, errors, refresh } = useAdminData();
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();

  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [actionError, setActionError] = useState("");
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());

  function toggleExpanded(reviewId: number) {
    setExpandedIds((prev) => {
      const next = new Set(prev);
      if (next.has(reviewId)) next.delete(reviewId);
      else next.add(reviewId);
      return next;
    });
  }

  const sorted = useMemo(
    () =>
      [...reviews].sort(
        (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      ),
    [reviews]
  );

  async function handleDelete(reviewId: number) {
    const ok = await confirm({
      title: "Delete this review?",
      description: "This permanently removes the review. This cannot be undone.",
      confirmLabel: "Delete review",
      destructive: true,
    });
    if (!ok) return;

    setDeletingId(reviewId);
    setActionError("");
    try {
      const res = await deleteReview(reviewId);
      if (res.success) {
        setReviews((prev) => prev.filter((r) => r.id !== reviewId));
      } else {
        setActionError(res.message || "Couldn't delete this review.");
      }
    } catch (err) {
      setActionError(extractErrorMessage(err));
    } finally {
      setDeletingId(null);
    }
  }

  return (
    <div className="mx-auto max-w-4xl">
      <AdminPageHeader
        title="Reviews"
        subtitle="Moderate guest reviews across every hotel."
        actions={
          <button
            onClick={refresh}
            disabled={isLoading}
            title="Refresh"
            aria-label="Refresh reviews"
            className="flex h-10 w-10 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
          >
            <RefreshCw size={16} className={isLoading ? "animate-spin" : ""} />
          </button>
        }
      />

      {actionError && (
        <p
          role="alert"
          className="mt-4 rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
        >
          {actionError}
        </p>
      )}

      <div className="mt-6">
        {isLoading && <AdminLoadingRows count={4} />}

        {!isLoading && errors.reviews && (
          <AdminErrorState
            message={`Failed to load reviews: ${errors.reviews}`}
            onRetry={refresh}
          />
        )}

        {!isLoading && !errors.reviews && sorted.length === 0 && (
          <AdminEmptyState
            icon={MessageSquareText}
            title="No reviews yet"
            description="Guest reviews will show up here once submitted."
          />
        )}

        {!isLoading && !errors.reviews && sorted.length > 0 && (
          <div className="flex flex-col gap-2.5">
            {sorted.map((review) => {
              const isBusy = deletingId === review.id;
              const isLong = review.comment.length > 220;
              const isExpanded = expandedIds.has(review.id);

              return (
                <div
                  key={review.id}
                  className="flex flex-col gap-3 rounded-[14px] border border-line bg-white p-4 sm:flex-row sm:items-start sm:justify-between"
                >
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <div className="flex items-center gap-0.5">
                        {Array.from({ length: 5 }).map((_, i) => (
                          <Star
                            key={i}
                            size={13}
                            className={
                              i < review.rating ? "fill-teal text-teal" : "text-line"
                            }
                          />
                        ))}
                      </div>
                      {review.hotel ? (
                        <Link
                          to={`/admin/hotels/${review.hotel.id}/edit`}
                          className="text-xs font-semibold text-ink hover:text-teal"
                        >
                          {review.hotel.name}
                        </Link>
                      ) : (
                        <span className="text-xs font-semibold text-muted">
                          Unknown hotel
                        </span>
                      )}
                      <span className="flex items-center gap-1 text-xs text-muted">
                        <UserRound size={11} />
                        {review.user ? review.user.name : "Unknown guest"}
                      </span>
                      <span className="text-xs text-muted">
                        {new Date(review.created_at).toLocaleDateString()}
                      </span>
                    </div>
                    {review.comment && (
                      <>
                        <p
                          className={`mt-2 text-sm leading-relaxed text-ink ${
                            isLong && !isExpanded ? "line-clamp-3" : ""
                          }`}
                        >
                          {review.comment}
                        </p>
                        {isLong && (
                          <button
                            onClick={() => toggleExpanded(review.id)}
                            className="mt-1 text-xs font-semibold text-teal hover:text-teal-dark cursor-pointer"
                          >
                            {isExpanded ? "Show less" : "Show more"}
                          </button>
                        )}
                      </>
                    )}
                  </div>

                  <button
                    onClick={() => handleDelete(review.id)}
                    disabled={isBusy}
                    title="Delete review"
                    aria-label="Delete review"
                    className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-red-300 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                  >
                    {isBusy ? (
                      <Loader2 size={14} className="animate-spin" />
                    ) : (
                      <Trash2 size={14} />
                    )}
                  </button>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <ConfirmDialog
        open={!!dialogState}
        title={dialogState?.title ?? ""}
        description={dialogState?.description ?? ""}
        confirmLabel={dialogState?.confirmLabel}
        cancelLabel={dialogState?.cancelLabel}
        destructive={dialogState?.destructive}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
      />
    </div>
  );
}
