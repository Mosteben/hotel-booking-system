import { useState } from "react";
import { BedDouble, Users } from "lucide-react";
import type { Room } from "@/types/room";

const STATUS_STYLES: Record<string, string> = {
  available: "bg-emerald-50 text-emerald-600",
  occupied: "bg-red-50 text-red-500",
  maintenance: "bg-amber-50 text-amber-600",
};

export function RoomCard({
  room,
  onBook,
}: {
  room: Room;
  onBook: (room: Room) => void;
}) {
  const isBookable = room.status === "available";
  // Cloudinary URLs are used directly, never rewritten - if one ever 404s
  // (deleted asset, network hiccup), fall back to the icon panel instead
  // of a broken-image glyph.
  const [imageFailed, setImageFailed] = useState(false);
  const imageUrl = !imageFailed ? room.image_url : null;

  return (
    <article className="flex flex-col overflow-hidden rounded-card bg-white shadow-[0_8px_24px_rgba(16,24,40,0.06)] transition-shadow hover:shadow-[0_16px_36px_rgba(16,24,40,0.12)] sm:flex-row">
      <div className="h-32 w-full shrink-0 overflow-hidden bg-cream sm:h-auto sm:w-40">
        {imageUrl ? (
          <img
            src={imageUrl}
            alt={`${room.type} · Room ${room.room_number}`}
            className="h-full w-full object-cover"
            onError={() => setImageFailed(true)}
          />
        ) : (
          <div className="flex h-full w-full items-center justify-center">
            <BedDouble size={32} className="text-teal" strokeWidth={1.5} />
          </div>
        )}
      </div>

      <div className="flex flex-1 flex-col gap-3 p-6 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="flex flex-wrap items-center gap-2">
            <h3 className="font-display text-base font-semibold text-ink">
              {room.type}
            </h3>
            <span
              className={`inline-flex items-center rounded-pill px-2.5 py-0.5 text-xs font-semibold capitalize ${
                STATUS_STYLES[room.status] ?? "bg-bg text-muted"
              }`}
            >
              {room.status}
            </span>
          </div>
          <p className="mt-1 text-sm text-muted">Room {room.room_number}</p>
          {room.description && (
            <p className="mt-1 max-w-md text-sm leading-relaxed text-muted">
              {room.description}
            </p>
          )}
          <p className="mt-2 flex items-center gap-1.5 text-sm text-muted">
            <Users size={14} />
            Up to {room.capacity} guests
          </p>
        </div>

        <div className="flex shrink-0 items-center justify-between gap-4 sm:flex-col sm:items-end sm:justify-center sm:gap-4">
          <p className="font-display text-lg font-bold text-ink">
            ${room.price_per_night.toFixed(2)}
            <span className="text-sm font-normal text-muted"> / night</span>
          </p>
          <button
            onClick={() => onBook(room)}
            disabled={!isBookable}
            className="rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
          >
            {isBookable ? "Select room" : "Unavailable"}
          </button>
        </div>
      </div>
    </article>
  );
}
