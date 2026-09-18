import { useMemo, useState } from "react";
import {
  Banknote,
  CreditCard,
  Loader2,
  RefreshCw,
  Search,
  Undo2,
  UserRound,
  Wallet,
  WalletCards,
  XCircle,
} from "lucide-react";
import { useAdminData } from "@/context/AdminDataContext";
import { updatePaymentStatus } from "@/api/paymentApi";
import { extractErrorMessage } from "@/api/client";
import { StatusBadge } from "@/components/common/StatusBadge";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import {
  AdminEmptyState,
  AdminErrorState,
  AdminLoadingRows,
} from "@/components/Admin/AdminStates";
import { ConfirmDialog } from "@/components/Admin/ConfirmDialog";
import { useConfirmDialog } from "@/hooks/useConfirmDialog";
import type { PaymentMethod, PaymentStatus } from "@/types/payment";

type Action = Extract<PaymentStatus, "paid" | "failed" | "refunded">;

const METHOD_ICON: Record<PaymentMethod, typeof CreditCard> = {
  card: CreditCard,
  cash: Wallet,
};

const STATUS_FILTERS: { label: string; value: PaymentStatus | "all" }[] = [
  { label: "All", value: "all" },
  { label: "Pending", value: "pending" },
  { label: "Paid", value: "paid" },
  { label: "Failed", value: "failed" },
  { label: "Refunded", value: "refunded" },
];

export function PaymentList() {
  const { payments, setPayments, isLoading, errors, refresh } = useAdminData();
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<PaymentStatus | "all">("all");
  const [mutatingId, setMutatingId] = useState<number | null>(null);
  const [mutationError, setMutationError] = useState("");

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();

    return payments.filter((payment) => {
      if (statusFilter !== "all" && payment.status !== statusFilter) return false;
      if (!query) return true;

      const haystack = [
        String(payment.id),
        String(payment.booking_id),
        payment.payment_method,
        payment.user?.name ?? "",
        payment.user?.email ?? "",
        payment.user?.phone ?? "",
        payment.hotel?.name ?? "",
      ]
        .join(" ")
        .toLowerCase();

      return haystack.includes(query);
    });
  }, [payments, search, statusFilter]);

  // The backend owns the actual transition rules (e.g. only a paid payment
  // can be refunded) - marking as refunded is confirmed since it's
  // typically irreversible for real money; paid/failed are quick one-click
  // actions consistent with how bookings confirm/cancel work.
  async function applyStatus(paymentId: number, next: Action) {
    if (mutatingId !== null) return;
    setMutatingId(paymentId);
    setMutationError("");
    try {
      const res = await updatePaymentStatus(paymentId, next);
      if (res.success) {
        // The status-update endpoint returns the bare payment (no resolved
        // user/booking/hotel/room) - merge just the fields that can
        // change, keep the already-resolved relationships from the list.
        const updated = res.data;
        setPayments((prev) =>
          prev.map((p) =>
            p.id === paymentId
              ? {
                  ...p,
                  status: updated.status,
                  transaction_id: updated.transaction_id,
                  updated_at: updated.updated_at,
                }
              : p
          )
        );
      } else {
        setMutationError(res.message || "Couldn't update this payment.");
      }
    } catch (err) {
      setMutationError(extractErrorMessage(err));
    } finally {
      setMutatingId(null);
    }
  }

  async function handleRefund(paymentId: number) {
    if (mutatingId !== null) return;
    const ok = await confirm({
      title: `Refund payment #${paymentId}?`,
      description: "This marks the payment as refunded. This cannot be undone.",
      confirmLabel: "Refund payment",
      destructive: true,
    });
    if (!ok) return;
    await applyStatus(paymentId, "refunded");
  }

  return (
    <div className="mx-auto max-w-6xl">
      <AdminPageHeader
        title="Payments"
        subtitle="Mark payments as paid, failed, or refunded."
        actions={
          <button
            onClick={refresh}
            disabled={isLoading}
            title="Refresh"
            aria-label="Refresh payments"
            className="flex h-10 w-10 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
          >
            <RefreshCw size={16} className={isLoading ? "animate-spin" : ""} />
          </button>
        }
      />

      <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative w-full sm:max-w-xs">
          <Search
            size={15}
            className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-muted"
          />
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by payment ID, booking ID, or guest..."
            className="w-full rounded-pill border border-line bg-white py-2.5 pl-10 pr-4 text-sm text-ink outline-none focus:border-teal focus:ring-4 focus:ring-teal/10"
          />
        </div>

        <div className="flex flex-wrap gap-1.5">
          {STATUS_FILTERS.map(({ label, value }) => (
            <button
              key={value}
              onClick={() => setStatusFilter(value)}
              className={`rounded-pill px-3.5 py-1.5 text-xs font-semibold transition-colors cursor-pointer ${
                statusFilter === value
                  ? "bg-ink text-white"
                  : "bg-white text-muted hover:text-ink border border-line"
              }`}
            >
              {label}
            </button>
          ))}
        </div>
      </div>

      {mutationError && (
        <p
          role="alert"
          className="mt-4 rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
        >
          {mutationError}
        </p>
      )}

      <div className="mt-6">
        {isLoading && <AdminLoadingRows count={5} />}

        {!isLoading && errors.payments && (
          <AdminErrorState
            message={`Failed to load payments: ${errors.payments}`}
            onRetry={refresh}
          />
        )}

        {!isLoading && !errors.payments && filtered.length === 0 && (
          <AdminEmptyState
            icon={WalletCards}
            title="No payments found"
            description={
              payments.length === 0
                ? "No payments have been recorded yet."
                : "Try a different search or filter."
            }
          />
        )}

        {!isLoading && !errors.payments && filtered.length > 0 && (
          <div className="flex flex-col gap-2.5">
            {filtered.map((payment) => {
              const isBusy = mutatingId === payment.id;
              const MethodIcon = METHOD_ICON[payment.payment_method] ?? Wallet;

              return (
                <div
                  key={payment.id}
                  className="flex flex-col gap-4 rounded-[14px] border border-line bg-white p-4 transition-shadow duration-200 hover:shadow-[var(--shadow-card)] sm:flex-row sm:items-center"
                >
                  <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-[11px] bg-cream text-teal">
                    <MethodIcon size={17} />
                  </div>

                  <div className="min-w-0 flex-1">
                    <p className="flex items-baseline gap-2">
                      <span className="font-display text-lg font-bold text-gold">
                        ${payment.amount.toFixed(2)}
                      </span>
                      <span className="text-xs font-medium capitalize text-muted">
                        {payment.payment_method}
                      </span>
                    </p>
                    <p className="mt-0.5 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted">
                      <span>Payment #{payment.id}</span>
                      <span>Booking #{payment.booking_id}</span>
                      <span className="flex items-center gap-1">
                        <UserRound size={12} />
                        {payment.user ? payment.user.name : "Unknown guest"}
                      </span>
                      {payment.hotel && (
                        <span>
                          {payment.hotel.name}
                          {payment.room ? ` · ${payment.room.room_number}` : ""}
                        </span>
                      )}
                      <span>{new Date(payment.created_at).toLocaleDateString()}</span>
                    </p>
                  </div>

                  <div className="flex shrink-0 flex-wrap items-center gap-2">
                    <StatusBadge status={payment.status} />

                    <button
                      onClick={() => applyStatus(payment.id, "paid")}
                      disabled={payment.status === "paid" || isBusy}
                      title="Mark as paid"
                      aria-label="Mark as paid"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-emerald-300 hover:text-emerald-600 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      {isBusy ? <Loader2 size={15} className="animate-spin" /> : <Banknote size={15} />}
                    </button>

                    <button
                      onClick={() => applyStatus(payment.id, "failed")}
                      disabled={payment.status === "failed" || isBusy}
                      title="Mark as failed"
                      aria-label="Mark as failed"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-red-300 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      {isBusy ? <Loader2 size={15} className="animate-spin" /> : <XCircle size={15} />}
                    </button>

                    <button
                      onClick={() => handleRefund(payment.id)}
                      disabled={payment.status === "refunded" || isBusy}
                      title="Refund payment"
                      aria-label="Refund payment"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      {isBusy ? <Loader2 size={15} className="animate-spin" /> : <Undo2 size={15} />}
                    </button>
                  </div>
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
