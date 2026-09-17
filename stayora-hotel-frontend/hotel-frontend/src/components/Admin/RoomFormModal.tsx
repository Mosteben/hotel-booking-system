import { useEffect, useState, type FormEvent } from "react";
import { Modal } from "@/components/Admin/Modal";
import { FormField, FormSelect } from "@/components/auth/FormField";
import { AuthButton } from "@/components/auth/AuthButton";
import { createRoom, updateRoom } from "@/api/roomApi";
import { extractErrorMessage } from "@/api/client";
import type { Room, RoomRequest, RoomStatus } from "@/types/room";

const EMPTY_FORM: RoomRequest = {
  room_number: "",
  type: "",
  description: "",
  price_per_night: 0,
  capacity: 1,
  status: "available",
};

export function RoomFormModal({
  open,
  onClose,
  hotelId,
  room,
  onSaved,
}: {
  open: boolean;
  onClose: () => void;
  hotelId: number;
  // null = create mode, otherwise edit this room.
  room: Room | null;
  onSaved: () => void;
}) {
  const [form, setForm] = useState<RoomRequest>(EMPTY_FORM);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (!open) return;
    setError("");
    setForm(
      room
        ? {
            room_number: room.room_number,
            type: room.type,
            description: room.description ?? "",
            price_per_night: room.price_per_night,
            capacity: room.capacity,
            status: room.status,
          }
        : EMPTY_FORM
    );
  }, [open, room]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");

    if (!form.room_number.trim() || !form.type.trim()) {
      setError("Room number and type are required.");
      return;
    }
    if (form.price_per_night <= 0) {
      setError("Price per night must be greater than zero.");
      return;
    }
    if (form.capacity <= 0) {
      setError("Capacity must be greater than zero.");
      return;
    }

    setIsSaving(true);
    try {
      if (room) {
        const res = await updateRoom(room.id, form);
        if (!res.success) {
          setError(res.message || "Couldn't save this room.");
          return;
        }
      } else {
        const res = await createRoom(hotelId, form);
        if (!res.success) {
          setError(res.message || "Couldn't create this room.");
          return;
        }
      }
      onSaved();
      onClose();
    } catch (err) {
      setError(extractErrorMessage(err));
    } finally {
      setIsSaving(false);
    }
  }

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={room ? "Edit room" : "New room"}
    >
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div className="grid grid-cols-2 gap-4">
          <FormField
            label="Room number"
            value={form.room_number}
            onChange={(e) => setForm((f) => ({ ...f, room_number: e.target.value }))}
            required
          />
          <FormField
            label="Type"
            placeholder="e.g. Deluxe"
            value={form.type}
            onChange={(e) => setForm((f) => ({ ...f, type: e.target.value }))}
            required
          />
        </div>

        <div className="flex flex-col gap-1.5">
          <label className="text-sm font-medium text-ink">
            Description <span className="text-xs font-normal text-muted">(optional)</span>
          </label>
          <textarea
            value={form.description}
            onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            rows={2}
            className="w-full resize-none rounded-[13px] border border-line bg-white px-4 py-3 text-sm text-ink outline-none transition-colors placeholder:text-muted focus:border-teal focus:ring-4 focus:ring-teal/10"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <FormField
            label="Price / night"
            type="number"
            min={0}
            step="0.01"
            value={form.price_per_night}
            onChange={(e) =>
              setForm((f) => ({ ...f, price_per_night: Number(e.target.value) }))
            }
            required
          />
          <FormField
            label="Capacity"
            type="number"
            min={1}
            value={form.capacity}
            onChange={(e) => setForm((f) => ({ ...f, capacity: Number(e.target.value) }))}
            required
          />
        </div>

        <FormSelect
          label="Status"
          value={form.status}
          onChange={(e) =>
            setForm((f) => ({ ...f, status: e.target.value as RoomStatus }))
          }
        >
          <option value="available">Available</option>
          <option value="occupied">Occupied</option>
          <option value="maintenance">Maintenance</option>
        </FormSelect>

        {error && (
          <p role="alert" className="text-xs font-medium text-red-500">
            {error}
          </p>
        )}

        <div className="mt-1 flex justify-end gap-2.5">
          <button
            type="button"
            onClick={onClose}
            className="rounded-pill border border-line px-4 py-2.5 text-sm font-semibold text-ink transition-colors hover:border-teal/40 cursor-pointer"
          >
            Cancel
          </button>
          <AuthButton isSubmitting={isSaving} className="mt-0 w-auto px-6">
            {room ? "Save changes" : "Create room"}
          </AuthButton>
        </div>
      </form>
    </Modal>
  );
}
