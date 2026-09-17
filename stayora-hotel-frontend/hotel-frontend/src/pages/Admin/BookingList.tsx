import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  BedDouble,
  Calendar,
  Check,
  ClipboardList,
  Loader2,
  RefreshCw,
  Search,
  UserRound,
  X,
} from "lucide-react";
import { useAdminData } from "@/context/AdminDataContext";
import { updateBookingStatus } from "@/api/bookingApi";
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
import type { BookingStatus } from "@/types/booking";

const STATUS_FILTERS: { label: string; value: BookingStatus | "all" }[] = [
  { label: "All", value: "all" },
  { label: "Pending", value: "pending" },
  { label: "Confirmed", value: "confirmed" },
  { label: "Cancelled", value: "cancelled" },
];

export function BookingList() {
  const { bookings, setBookings, isLoading, errors, refresh } = useAdminData();
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();

  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<BookingStatus | "all">("all");
  const [mutatingId, setMutatingId] = useState<number | null>(null);
  const [mutationError, setMutationError] = useState("");

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();

    return bookings.filter((booking) => {
      if (statusFilter !== "all" && booking.status !== statusFilter) return false;
      if (!query) return true;

      const haystack = [
        String(booking.id),
        booking.hotel?.name ?? booking.hotel_name_snapshot,
        booking.room?.type ?? booking.room_type_snapshot,
        booking.room?.room_number ?? booking.room_number_snapshot,
        booking.user?.name ?? "",
        booking.user?.email ?? "",
        booking.user?.phone ?? "",
      ]
        .join(" ")
        .toLowerCase();

      return haystack.includes(query);
    });
  }, [bookings, search, statusFilter]);

  async function handleConfirmBooking(bookingId: number) {
    if (mutatingId !== null) return;
    setMutatingId(bookingId);
    setMutationError("");
    try {
      const updated = await updateBookingStatus(bookingId, "confirmed");
      // The status-update endpoint returns the bare booking (no resolved
      // user/hotel/room) - merge just the fields that can change, keep the
      // already-resolved relationships from the original list.
      setBookings((prev) =>
        prev.map((b) =>
          b.id === bookingId
            ? { ...b, status: updated.status, updated_at: updated.updated_at }
            : b
        )
      );
    } catch (err) {
      setMutationError(extractErrorMessage(err));
    } finally {
      setMutatingId(null);
    }
  }

  async function handleCancelBooking(bookingId: number) {
    if (mutatingId !== null) return;

    const ok = await confirm({
      title: `Cancel booking #${bookingId}?`,
      description:
        "The guest will see this booking as cancelled. This cannot be undone.",
      confirmLabel: "Cancel booking",
      destructive: true,
    });
    if (!ok) return;

    setMutatingId(bookingId);
    setMutationError("");
    try {
      const updated = await updateBookingStatus(bookingId, "cancelled");
      setBookings((prev) =>
        prev.map((b) =>
          b.id === bookingId
            ? { ...b, status: updated.status, updated_at: updated.updated_at }
            : b
        )
      );
    } catch (err) {
      setMutationError(extractErrorMessage(err));
    } finally {
      setMutatingId(null);
    }
  }

  return (
    <div className="mx-auto max-w-6xl">
      <AdminPageHeader
        title="Bookings"
        subtitle="Confirm or cancel any guest's reservation."
        actions={
          <button
            onClick={refresh}
            disabled={isLoading}
            title="Refresh"
            aria-label="Refresh bookings"
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
            placeholder="Search by booking ID, guest, hotel, or room..."
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

        {!isLoading && errors.bookings && (
          <AdminErrorState
            message={`Failed to load bookings: ${errors.bookings}`}
            onRetry={refresh}
          />
        )}

        {!isLoading && !errors.bookings && filtered.length === 0 && (
          <AdminEmptyState
            icon={ClipboardList}
            title="No bookings found"
            description={
              bookings.length === 0
                ? "No bookings have been made yet."
                : "Try a different search or filter."
            }
          />
        )}

        {!isLoading && !errors.bookings && filtered.length > 0 && (
          <div className="flex flex-col gap-2.5">
            {filtered.map((booking) => {
              const isBusy = mutatingId === booking.id;
              const canConfirm = booking.status === "pending";
              const canCancel =
                booking.status === "pending" || booking.status === "confirmed";

              return (
                <div
                  key={booking.id}
                  className={`flex flex-col gap-4 rounded-[14px] border bg-white p-4 sm:flex-row sm:items-center ${
                    booking.status === "pending" ? "border-amber-200" : "border-line"
                  }`}
                >
                  <div className="flex h-11 w-11 shrink-0 items-center justify-center rounded-[11px] bg-cream text-teal">
                    <BedDouble size={17} />
                  </div>

                  <div className="min-w-0 flex-1">
                    <Link
                      to={booking.hotel ? `/admin/hotels/${booking.hotel.id}/edit` : "#"}
                      className="truncate text-sm font-semibold text-ink hover:text-teal"
                    >
                      {booking.hotel
                        ? booking.hotel.name
                        : booking.hotel_name_snapshot || `Room #${booking.room?.id ?? "?"}`}
                      {booking.room
                        ? ` · ${booking.room.type} ${booking.room.room_number}`
                        : booking.room_number_snapshot
                          ? ` · ${booking.room_type_snapshot} ${booking.room_number_snapshot} (deleted)`
                          : ""}
                    </Link>
                    <p className="mt-0.5 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted">
                      <span>#{booking.id}</span>
                      <span className="flex items-center gap-1">
                        <UserRound size={12} />
                        {booking.user ? booking.user.name : "Unknown guest"}
                      </span>
                      <span className="flex items-center gap-1">
                        <Calendar size={12} />
                        {new Date(booking.check_in).toLocaleDateString()} –{" "}
                        {new Date(booking.check_out).toLocaleDateString()}
                      </span>
                      <span>{booking.guests} guest{booking.guests === 1 ? "" : "s"}</span>
                      <span className="font-semibold text-ink">
                        ${booking.total_price.toFixed(2)}
                      </span>
                    </p>
                  </div>

                  <div className="flex shrink-0 items-center gap-2">
                    <StatusBadge status={booking.status} />

                    <button
                      onClick={() => handleConfirmBooking(booking.id)}
                      disabled={!canConfirm || isBusy}
                      title="Confirm booking"
                      aria-label="Confirm booking"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-emerald-300 hover:text-emerald-600 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      {isBusy ? (
                        <Loader2 size={15} className="animate-spin" />
                      ) : (
                        <Check size={15} />
                      )}
                    </button>

                    <button
                      onClick={() => handleCancelBooking(booking.id)}
                      disabled={!canCancel || isBusy}
                      title="Cancel booking"
                      aria-label="Cancel booking"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-red-300 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      {isBusy ? (
                        <Loader2 size={15} className="animate-spin" />
                      ) : (
                        <X size={15} />
                      )}
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
