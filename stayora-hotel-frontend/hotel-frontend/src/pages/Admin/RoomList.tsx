import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  BedDouble,
  ImagePlus,
  Loader2,
  Pencil,
  Plus,
  RefreshCw,
  Trash2,
  Users,
} from "lucide-react";
import { getHotelByID } from "@/api/hotelApi";
import { getRoomsByHotelID, deleteRoom } from "@/api/roomApi";
import { useAdminData } from "@/context/AdminDataContext";
import { extractErrorMessage } from "@/api/client";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import { AdminBreadcrumb } from "@/components/Admin/AdminBreadcrumb";
import {
  AdminEmptyState,
  AdminErrorState,
  AdminLoadingRows,
} from "@/components/Admin/AdminStates";
import { ConfirmDialog } from "@/components/Admin/ConfirmDialog";
import { RoomFormModal } from "@/components/Admin/RoomFormModal";
import { RoomImagesModal } from "@/components/Admin/RoomImagesModal";
import { useConfirmDialog } from "@/hooks/useConfirmDialog";
import type { Hotel } from "@/types/hotel";
import type { Room, RoomStatus } from "@/types/room";

const STATUS_STYLES: Record<RoomStatus, string> = {
  available: "bg-emerald-50 text-emerald-600",
  occupied: "bg-red-50 text-red-500",
  maintenance: "bg-amber-50 text-amber-600",
};

type Status = "loading" | "success" | "error";

export function RoomList() {
  const { hotelId } = useParams<{ hotelId: string }>();
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();
  const { refresh: refreshAdminData } = useAdminData();

  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [hotel, setHotel] = useState<Hotel | null>(null);
  const [rooms, setRooms] = useState<Room[]>([]);

  const [formOpen, setFormOpen] = useState(false);
  const [editingRoom, setEditingRoom] = useState<Room | null>(null);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [imagesModalRoom, setImagesModalRoom] = useState<Room | null>(null);

  async function load() {
    if (!hotelId) return;
    setStatus("loading");
    setErrorMessage("");
    try {
      const [hotelRes, roomsRes] = await Promise.all([
        getHotelByID(hotelId),
        getRoomsByHotelID(hotelId),
      ]);

      if (!hotelRes.success) {
        setErrorMessage(hotelRes.message || "Couldn't load this hotel.");
        setStatus("error");
        return;
      }
      setHotel(hotelRes.data);
      setRooms(roomsRes.success ? roomsRes.data : []);
      setStatus("success");
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [hotelId]);

  function openCreate() {
    setEditingRoom(null);
    setFormOpen(true);
  }

  function openEdit(room: Room) {
    setEditingRoom(room);
    setFormOpen(true);
  }

  async function handleDelete(room: Room) {
    const ok = await confirm({
      title: `Delete room ${room.room_number}?`,
      description:
        "This permanently removes the room. Existing bookings for it are not affected, but it can no longer be booked.",
      confirmLabel: "Delete room",
      destructive: true,
    });
    if (!ok) return;

    setDeletingId(room.id);
    try {
      const res = await deleteRoom(room.id);
      if (res.success) {
        setRooms((prev) => prev.filter((r) => r.id !== room.id));
        refreshAdminData();
      } else {
        setErrorMessage(res.message || "Couldn't delete this room.");
      }
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
    } finally {
      setDeletingId(null);
    }
  }

  return (
    <div className="mx-auto max-w-5xl">
      <AdminBreadcrumb
        items={[
          { label: "Hotels", to: "/admin/hotels" },
          {
            label: hotel?.name ?? "Hotel",
            to: hotelId ? `/admin/hotels/${hotelId}/edit` : undefined,
          },
          { label: "Rooms" },
        ]}
      />

      <div className="mt-3">
        <AdminPageHeader
          title="Rooms"
          subtitle="Manage the rooms guests can book at this hotel."
          actions={
            <>
              <button
                onClick={load}
                disabled={status === "loading"}
                title="Refresh"
                aria-label="Refresh rooms"
                className="flex h-10 w-10 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
              >
                <RefreshCw size={16} className={status === "loading" ? "animate-spin" : ""} />
              </button>
              <button
                onClick={openCreate}
                className="flex items-center gap-2 rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark cursor-pointer"
              >
                <Plus size={16} />
                New room
              </button>
            </>
          }
        />
      </div>

      {errorMessage && status !== "error" && (
        <p
          role="alert"
          className="mt-4 rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
        >
          {errorMessage}
        </p>
      )}

      <div className="mt-6">
        {status === "loading" && <AdminLoadingRows count={4} />}

        {status === "error" && (
          <AdminErrorState message={errorMessage} onRetry={load} />
        )}

        {status === "success" && rooms.length === 0 && (
          <AdminEmptyState
            icon={BedDouble}
            title="No rooms yet"
            description="This hotel can't be booked until at least one room is added."
          />
        )}

        {status === "success" && rooms.length > 0 && (
          <div className="flex flex-col gap-2.5">
            {rooms.map((room) => {
              const isBusy = deletingId === room.id;
              return (
                <div
                  key={room.id}
                  className="flex flex-col gap-4 rounded-[14px] border border-line bg-white p-4 transition-shadow duration-200 hover:shadow-[var(--shadow-card)] sm:flex-row sm:items-center"
                >
                  <div className="h-14 w-14 shrink-0 overflow-hidden rounded-[11px] bg-cream">
                    {room.image_url ? (
                      <img
                        src={room.image_url}
                        alt=""
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center text-teal">
                        <BedDouble size={18} />
                      </div>
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <p className="flex flex-wrap items-center gap-2 text-sm font-semibold text-ink">
                      {room.type}
                      <span
                        className={`inline-flex items-center rounded-pill px-2 py-0.5 text-[11px] font-semibold capitalize ${STATUS_STYLES[room.status]}`}
                      >
                        {room.status}
                      </span>
                    </p>
                    <p className="mt-0.5 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted">
                      <span>Room {room.room_number}</span>
                      <span className="flex items-center gap-1">
                        <Users size={12} />
                        Up to {room.capacity}
                      </span>
                      <span className="font-semibold text-ink">
                        ${room.price_per_night.toFixed(2)} / night
                      </span>
                      <span>
                        {room.images?.length ?? 0} photo{(room.images?.length ?? 0) === 1 ? "" : "s"}
                      </span>
                    </p>
                  </div>

                  <div className="flex shrink-0 items-center gap-2">
                    <button
                      onClick={() => setImagesModalRoom(room)}
                      title="Manage photos"
                      aria-label="Manage photos"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal cursor-pointer"
                    >
                      <ImagePlus size={14} />
                    </button>
                    <button
                      onClick={() => openEdit(room)}
                      title="Edit room"
                      aria-label="Edit room"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal cursor-pointer"
                    >
                      <Pencil size={14} />
                    </button>
                    <button
                      onClick={() => handleDelete(room)}
                      disabled={isBusy}
                      title="Delete room"
                      aria-label="Delete room"
                      className="flex h-9 w-9 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-red-300 hover:text-red-500 disabled:cursor-not-allowed disabled:opacity-40 cursor-pointer"
                    >
                      {isBusy ? (
                        <Loader2 size={14} className="animate-spin" />
                      ) : (
                        <Trash2 size={14} />
                      )}
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>

      {hotelId && (
        <RoomFormModal
          open={formOpen}
          onClose={() => setFormOpen(false)}
          hotelId={Number(hotelId)}
          room={editingRoom}
          onSaved={() => {
            load();
            refreshAdminData();
          }}
        />
      )}

      <RoomImagesModal
        open={!!imagesModalRoom}
        onClose={() => setImagesModalRoom(null)}
        room={imagesModalRoom}
        onImagesChanged={() => {
          load();
          refreshAdminData();
        }}
      />

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
