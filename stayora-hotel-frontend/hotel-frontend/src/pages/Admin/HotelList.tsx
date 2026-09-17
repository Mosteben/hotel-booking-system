import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  BedDouble,
  Building2,
  ImageOff,
  Pencil,
  Plus,
  RefreshCw,
  Search,
  Star,
  Trash2,
} from "lucide-react";
import { useAdminData } from "@/context/AdminDataContext";
import { deleteHotel } from "@/api/hotelApi";
import { extractErrorMessage } from "@/api/client";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import {
  AdminEmptyState,
  AdminErrorState,
  AdminLoadingRows,
} from "@/components/Admin/AdminStates";
import { ConfirmDialog } from "@/components/Admin/ConfirmDialog";
import { useConfirmDialog } from "@/hooks/useConfirmDialog";

export function HotelList() {
  const { hotels, isLoading, errors, refresh } = useAdminData();
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();

  const [search, setSearch] = useState("");
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [actionError, setActionError] = useState("");

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) return hotels;
    return hotels.filter((h) =>
      [h.name, h.city, h.country].join(" ").toLowerCase().includes(query)
    );
  }, [hotels, search]);

  async function handleDelete(hotelId: number, name: string) {
    const ok = await confirm({
      title: `Delete ${name}?`,
      description:
        "This permanently removes the hotel, its rooms, and its photos. This cannot be undone.",
      confirmLabel: "Delete hotel",
      destructive: true,
    });
    if (!ok) return;

    setDeletingId(hotelId);
    setActionError("");
    try {
      const res = await deleteHotel(hotelId);
      if (res.success) {
        await refresh();
      } else {
        setActionError(res.message || "Couldn't delete this hotel.");
      }
    } catch (err) {
      setActionError(extractErrorMessage(err));
    } finally {
      setDeletingId(null);
    }
  }

  return (
    <div className="mx-auto max-w-6xl">
      <AdminPageHeader
        title="Hotels"
        subtitle="Create hotels, manage their rooms, and manage their photos."
        actions={
          <>
            <button
              onClick={refresh}
              disabled={isLoading}
              title="Refresh"
              aria-label="Refresh hotels"
              className="flex h-10 w-10 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
            >
              <RefreshCw size={16} className={isLoading ? "animate-spin" : ""} />
            </button>
            <Link
              to="/admin/hotels/new"
              className="flex items-center gap-2 rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark"
            >
              <Plus size={16} />
              New hotel
            </Link>
          </>
        }
      />

      <div className="mt-6 relative max-w-xs">
        <Search
          size={15}
          className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-muted"
        />
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search hotels..."
          className="w-full rounded-pill border border-line bg-white py-2.5 pl-10 pr-4 text-sm text-ink outline-none focus:border-teal focus:ring-4 focus:ring-teal/10"
        />
      </div>

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

        {!isLoading && errors.hotels && (
          <AdminErrorState message={`Failed to load hotels: ${errors.hotels}`} onRetry={refresh} />
        )}

        {!isLoading && !errors.hotels && filtered.length === 0 && (
          <AdminEmptyState
            icon={Building2}
            title="No hotels found"
            description={
              hotels.length === 0
                ? "Create your first hotel to get started."
                : "Try a different search."
            }
          />
        )}

        {!isLoading && !errors.hotels && filtered.length > 0 && (
          <div className="flex flex-col gap-2.5">
            {filtered.map((hotel) => {
              const isBusy = deletingId === hotel.id;
              return (
                <div
                  key={hotel.id}
                  className="flex flex-col gap-4 rounded-[14px] border border-line bg-white p-4 sm:flex-row sm:items-center"
                >
                  <div className="h-14 w-14 shrink-0 overflow-hidden rounded-[12px] bg-cream">
                    {hotel.image_url ? (
                      <img
                        src={hotel.image_url}
                        alt={hotel.name}
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center text-teal">
                        <Building2 size={20} />
                      </div>
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <p className="truncate text-sm font-semibold text-ink">{hotel.name}</p>
                    <p className="truncate text-xs text-muted">
                      {hotel.city}
                      {hotel.country ? `, ${hotel.country}` : ""}
                    </p>
                    <div className="mt-1.5 flex flex-wrap items-center gap-1.5">
                      {hotel.stars > 0 && (
                        <span className="inline-flex items-center gap-1 rounded-pill bg-bg px-2 py-0.5 text-[11px] font-semibold text-ink">
                          <Star size={10} className="fill-teal text-teal" />
                          {hotel.stars}
                        </span>
                      )}
                      <span className="inline-flex items-center gap-1 rounded-pill bg-bg px-2 py-0.5 text-[11px] font-semibold text-muted">
                        <BedDouble size={10} />
                        {hotel.roomCount} room{hotel.roomCount === 1 ? "" : "s"}
                      </span>
                      {!hotel.image_url && (
                        <span className="inline-flex items-center gap-1 rounded-pill bg-amber-50 px-2 py-0.5 text-[11px] font-semibold text-amber-600">
                          <ImageOff size={10} />
                          No photo
                        </span>
                      )}
                    </div>
                  </div>

                  <div className="flex shrink-0 items-center gap-2">
                    <Link
                      to={`/admin/hotels/${hotel.id}/rooms`}
                      className="rounded-pill border border-line px-3.5 py-2 text-xs font-semibold text-ink transition-colors hover:border-teal/40 hover:text-teal"
                    >
                      Manage rooms
                    </Link>
                    <Link
                      to={`/admin/hotels/${hotel.id}/edit`}
                      title="Edit hotel"
                      aria-label="Edit hotel"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal"
                    >
                      <Pencil size={14} />
                    </Link>
                    <button
                      onClick={() => handleDelete(hotel.id, hotel.name)}
                      disabled={isBusy}
                      title="Delete hotel"
                      aria-label="Delete hotel"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-red-300 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      <Trash2 size={14} />
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
