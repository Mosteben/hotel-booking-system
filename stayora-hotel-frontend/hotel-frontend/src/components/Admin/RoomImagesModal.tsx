import { useEffect, useRef, useState } from "react";
import { BedDouble, ImagePlus, Loader2, Star, Trash2, X } from "lucide-react";
import { Modal } from "@/components/Admin/Modal";
import { ConfirmDialog } from "@/components/Admin/ConfirmDialog";
import { useConfirmDialog } from "@/hooks/useConfirmDialog";
import {
  deleteRoomImage,
  setMainRoomImage,
  uploadRoomImages,
} from "@/api/roomImageApi";
import { extractErrorMessage } from "@/api/client";
import type { Room, RoomImage } from "@/types/room";

const MAX_IMAGE_SIZE = 5 * 1024 * 1024;
const ALLOWED_IMAGE_TYPES = ["image/jpeg", "image/png", "image/webp"];

function validateImageFile(file: File): string | null {
  if (!ALLOWED_IMAGE_TYPES.includes(file.type)) {
    return `${file.name} must be a JPEG, PNG, or WebP image.`;
  }
  if (file.size > MAX_IMAGE_SIZE) {
    return `${file.name} is larger than 5MB.`;
  }
  return null;
}

interface PendingImage {
  file: File;
  previewUrl: string;
}

export function RoomImagesModal({
  open,
  onClose,
  room,
  onImagesChanged,
}: {
  open: boolean;
  onClose: () => void;
  room: Room | null;
  // Called after any successful upload/delete/set-main so the parent list
  // can refresh its own copy of this room's images/thumbnail.
  onImagesChanged: () => void;
}) {
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();

  const [images, setImages] = useState<RoomImage[]>([]);
  const [pendingImages, setPendingImages] = useState<PendingImage[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState("");
  const [busyImageId, setBusyImageId] = useState<number | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!open) return;
    setImages(room?.images ?? []);
    setError("");
  }, [open, room]);

  useEffect(() => {
    return () => {
      pendingImages.forEach((p) => URL.revokeObjectURL(p.previewUrl));
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function handleSelectFiles(fileList: FileList | null) {
    if (!fileList || fileList.length === 0) return;
    setError("");

    const next: PendingImage[] = [];
    for (const file of Array.from(fileList)) {
      const validationError = validateImageFile(file);
      if (validationError) {
        setError(validationError);
        continue;
      }
      next.push({ file, previewUrl: URL.createObjectURL(file) });
    }

    setPendingImages((prev) => [...prev, ...next]);
    if (fileInputRef.current) fileInputRef.current.value = "";
  }

  function removePendingImage(previewUrl: string) {
    setPendingImages((prev) => {
      const target = prev.find((p) => p.previewUrl === previewUrl);
      if (target) URL.revokeObjectURL(target.previewUrl);
      return prev.filter((p) => p.previewUrl !== previewUrl);
    });
  }

  async function handleUpload() {
    if (!room || pendingImages.length === 0 || isUploading) return;
    setIsUploading(true);
    setError("");
    try {
      const res = await uploadRoomImages(
        room.id,
        pendingImages.map((p) => p.file)
      );
      if (res.success) {
        pendingImages.forEach((p) => URL.revokeObjectURL(p.previewUrl));
        setPendingImages([]);
        setImages((prev) => [...prev, ...res.data]);
        onImagesChanged();
      } else {
        setError(res.message || "Upload failed.");
      }
    } catch (err) {
      setError(extractErrorMessage(err));
    } finally {
      setIsUploading(false);
    }
  }

  async function handleDeleteImage(imageId: number) {
    if (!room) return;

    const ok = await confirm({
      title: "Delete this photo?",
      description:
        "This removes it from the room's gallery. If it's the main photo, another one is promoted automatically.",
      confirmLabel: "Delete photo",
      destructive: true,
    });
    if (!ok) return;

    setBusyImageId(imageId);
    setError("");
    try {
      const res = await deleteRoomImage(room.id, imageId);
      if (res.success) {
        setImages((prev) => {
          const deleted = prev.find((img) => img.id === imageId);
          const rest = prev.filter((img) => img.id !== imageId);
          if (deleted?.is_main && rest.length > 0) {
            rest[0] = { ...rest[0], is_main: true };
          }
          return rest;
        });
        onImagesChanged();
      } else {
        setError(res.message || "Couldn't delete this image.");
      }
    } catch (err) {
      setError(extractErrorMessage(err));
    } finally {
      setBusyImageId(null);
    }
  }

  async function handleSetMain(imageId: number) {
    if (!room) return;
    setBusyImageId(imageId);
    setError("");
    try {
      const res = await setMainRoomImage(room.id, imageId);
      if (res.success) {
        setImages((prev) => prev.map((img) => ({ ...img, is_main: img.id === imageId })));
        onImagesChanged();
      } else {
        setError(res.message || "Couldn't update the main image.");
      }
    } catch (err) {
      setError(extractErrorMessage(err));
    } finally {
      setBusyImageId(null);
    }
  }

  if (!room) return null;

  return (
    <>
      <Modal
        open={open}
        onClose={onClose}
        title={`Photos · Room ${room.room_number}`}
        maxWidth="max-w-2xl"
      >
        <p className="text-sm text-muted">
          The main photo is what guests see first when browsing this room.
        </p>

        {images.length > 0 && (
          <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-3">
            {images.map((image) => {
              const isBusy = busyImageId === image.id;
              return (
                <div
                  key={image.id}
                  className="group relative aspect-[4/3] overflow-hidden rounded-[13px] bg-cream"
                >
                  <img src={image.url} alt="" className="h-full w-full object-cover" />
                  {image.is_main && (
                    <span className="absolute left-2 top-2 flex items-center gap-1 rounded-pill bg-white/95 px-2 py-1 text-[11px] font-semibold text-teal shadow-sm">
                      <Star size={10} className="fill-teal text-teal" />
                      Main
                    </span>
                  )}
                  <div className="absolute inset-0 flex items-end justify-end gap-1.5 bg-gradient-to-t from-black/50 via-transparent to-transparent p-2 opacity-0 transition-opacity group-hover:opacity-100">
                    {!image.is_main && (
                      <button
                        type="button"
                        onClick={() => handleSetMain(image.id)}
                        disabled={isBusy}
                        title="Set as main photo"
                        aria-label="Set as main photo"
                        className="flex h-8 w-8 items-center justify-center rounded-full bg-white/95 text-ink transition-colors hover:text-teal disabled:opacity-50 cursor-pointer"
                      >
                        {isBusy ? (
                          <Loader2 size={14} className="animate-spin" />
                        ) : (
                          <Star size={14} />
                        )}
                      </button>
                    )}
                    <button
                      type="button"
                      onClick={() => handleDeleteImage(image.id)}
                      disabled={isBusy}
                      title="Delete photo"
                      aria-label="Delete photo"
                      className="flex h-8 w-8 items-center justify-center rounded-full bg-white/95 text-red-500 transition-colors hover:bg-red-50 disabled:opacity-50 cursor-pointer"
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

        {images.length === 0 && pendingImages.length === 0 && (
          <div className="mt-4 flex flex-col items-center gap-2 rounded-[14px] border border-dashed border-line px-6 py-10 text-center">
            <span className="flex h-11 w-11 items-center justify-center rounded-full bg-bg text-teal">
              <BedDouble size={20} />
            </span>
            <p className="text-xs font-medium text-muted">No photos yet</p>
          </div>
        )}

        {pendingImages.length > 0 && (
          <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-3">
            {pendingImages.map((p) => (
              <div
                key={p.previewUrl}
                className="relative aspect-[4/3] overflow-hidden rounded-[13px] border-2 border-dashed border-teal/40 bg-cream"
              >
                <img src={p.previewUrl} alt="" className="h-full w-full object-cover opacity-80" />
                <button
                  type="button"
                  onClick={() => removePendingImage(p.previewUrl)}
                  title="Remove"
                  aria-label="Remove selected image"
                  className="absolute right-2 top-2 flex h-7 w-7 items-center justify-center rounded-full bg-white/95 text-ink transition-colors hover:text-red-500 cursor-pointer"
                >
                  <X size={14} />
                </button>
              </div>
            ))}
          </div>
        )}

        {error && (
          <p role="alert" className="mt-3 text-xs font-medium text-red-500">
            {error}
          </p>
        )}

        <div className="mt-4 flex flex-wrap items-center gap-3">
          <input
            ref={fileInputRef}
            type="file"
            accept="image/jpeg,image/png,image/webp"
            multiple
            className="hidden"
            onChange={(e) => handleSelectFiles(e.target.files)}
          />
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="flex items-center gap-2 rounded-pill border border-line px-4 py-2.5 text-sm font-semibold text-ink transition-colors hover:border-teal hover:text-teal cursor-pointer"
          >
            <ImagePlus size={15} />
            Choose photos
          </button>

          {pendingImages.length > 0 && (
            <button
              type="button"
              onClick={handleUpload}
              disabled={isUploading}
              className="flex items-center gap-2 rounded-pill bg-teal px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark disabled:cursor-not-allowed disabled:opacity-60 cursor-pointer"
            >
              {isUploading && <Loader2 size={14} className="animate-spin" />}
              Upload {pendingImages.length} photo{pendingImages.length === 1 ? "" : "s"}
            </button>
          )}
        </div>
        <p className="mt-2 text-xs text-muted">
          JPEG, PNG, or WebP. Up to 5MB per photo.
        </p>
      </Modal>

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
    </>
  );
}
