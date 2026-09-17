import { useEffect, useRef, useState, type FormEvent } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  BedDouble,
  Building2,
  CheckCircle2,
  ImagePlus,
  Loader2,
  Star,
  Trash2,
  TriangleAlert,
  X,
} from "lucide-react";
import { FormField, FormSelect } from "@/components/auth/FormField";
import { AuthButton } from "@/components/auth/AuthButton";
import { AdminBreadcrumb } from "@/components/Admin/AdminBreadcrumb";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import { ConfirmDialog } from "@/components/Admin/ConfirmDialog";
import { useConfirmDialog } from "@/hooks/useConfirmDialog";
import { useAdminData } from "@/context/AdminDataContext";
import { createHotel, getHotelByID, updateHotel } from "@/api/hotelApi";
import {
  deleteHotelImage,
  setMainHotelImage,
  uploadHotelImages,
} from "@/api/hotelImageApi";
import { extractErrorMessage } from "@/api/client";
import type { Hotel, HotelImage, HotelRequest } from "@/types/hotel";

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

const emptyForm: HotelRequest = {
  name: "",
  description: "",
  address: "",
  city: "",
  country: "",
  phone: "",
  email: "",
  stars: 3,
};

export function HotelForm() {
  const { id } = useParams<{ id: string }>();
  const isEditing = Boolean(id);
  const navigate = useNavigate();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { refresh } = useAdminData();
  const { confirm, dialogState, handleConfirm, handleCancel } = useConfirmDialog();

  const [loadingHotel, setLoadingHotel] = useState(isEditing);
  const [loadError, setLoadError] = useState("");

  const [form, setForm] = useState<HotelRequest>(emptyForm);
  const [isSaving, setIsSaving] = useState(false);
  const [saveError, setSaveError] = useState("");
  const [saved, setSaved] = useState(false);

  const [images, setImages] = useState<HotelImage[]>([]);
  const [pendingImages, setPendingImages] = useState<PendingImage[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [imageError, setImageError] = useState("");
  const [busyImageId, setBusyImageId] = useState<number | null>(null);

  useEffect(() => {
    if (!id) return;

    let cancelled = false;
    setLoadingHotel(true);
    getHotelByID(id)
      .then((res) => {
        if (cancelled) return;
        if (res.success) {
          const hotel: Hotel = res.data;
          setForm({
            name: hotel.name,
            description: hotel.description ?? "",
            address: hotel.address,
            city: hotel.city,
            country: hotel.country,
            phone: hotel.phone ?? "",
            email: hotel.email ?? "",
            stars: hotel.stars,
          });
          setImages(hotel.images ?? []);
        } else {
          setLoadError(res.message || "Couldn't load this hotel.");
        }
      })
      .catch((err) => setLoadError(extractErrorMessage(err)))
      .finally(() => {
        if (!cancelled) setLoadingHotel(false);
      });

    return () => {
      cancelled = true;
    };
  }, [id]);

  useEffect(() => {
    return () => {
      pendingImages.forEach((p) => URL.revokeObjectURL(p.previewUrl));
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function reloadImages(hotelId: string) {
    const res = await getHotelByID(hotelId);
    if (res.success) {
      setImages(res.data.images ?? []);
    }
    // Keep the shared Admin dataset (sidebar/Overview room+image counts) in
    // sync without forcing every page to individually refetch.
    refresh();
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setSaveError("");
    setSaved(false);

    if (!form.name.trim() || !form.address.trim() || !form.city.trim() || !form.country.trim()) {
      setSaveError("Name, address, city, and country are required.");
      return;
    }

    if (form.stars < 1 || form.stars > 5) {
      setSaveError("Stars must be between 1 and 5.");
      return;
    }

    setIsSaving(true);
    try {
      if (isEditing && id) {
        const res = await updateHotel(id, form);
        if (res.success) {
          setSaved(true);
          refresh();
        } else {
          setSaveError(res.message || "Couldn't save this hotel.");
        }
      } else {
        const res = await createHotel(form);
        if (res.success) {
          refresh();
          navigate(`/admin/hotels/${res.data.id}/edit`, { replace: true });
        } else {
          setSaveError(res.message || "Couldn't create this hotel.");
        }
      }
    } catch (err) {
      setSaveError(extractErrorMessage(err));
    } finally {
      setIsSaving(false);
    }
  }

  function handleSelectFiles(fileList: FileList | null) {
    if (!fileList || fileList.length === 0) return;
    setImageError("");

    const next: PendingImage[] = [];
    for (const file of Array.from(fileList)) {
      const error = validateImageFile(file);
      if (error) {
        setImageError(error);
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
    if (!id || pendingImages.length === 0) return;
    setIsUploading(true);
    setImageError("");
    try {
      const res = await uploadHotelImages(
        id,
        pendingImages.map((p) => p.file)
      );
      if (res.success) {
        pendingImages.forEach((p) => URL.revokeObjectURL(p.previewUrl));
        setPendingImages([]);
        await reloadImages(id);
      } else {
        setImageError(res.message || "Upload failed.");
      }
    } catch (err) {
      setImageError(extractErrorMessage(err));
    } finally {
      setIsUploading(false);
    }
  }

  async function handleDeleteImage(imageId: number) {
    if (!id) return;

    const ok = await confirm({
      title: "Delete this photo?",
      description:
        "This removes it from the hotel's gallery. If it's the main photo, another one is promoted automatically.",
      confirmLabel: "Delete photo",
      destructive: true,
    });
    if (!ok) return;

    setBusyImageId(imageId);
    setImageError("");
    try {
      const res = await deleteHotelImage(id, imageId);
      if (res.success) {
        await reloadImages(id);
      } else {
        setImageError(res.message || "Couldn't delete this image.");
      }
    } catch (err) {
      setImageError(extractErrorMessage(err));
    } finally {
      setBusyImageId(null);
    }
  }

  async function handleSetMain(imageId: number) {
    if (!id) return;
    setBusyImageId(imageId);
    setImageError("");
    try {
      const res = await setMainHotelImage(id, imageId);
      if (res.success) {
        await reloadImages(id);
      } else {
        setImageError(res.message || "Couldn't update the main image.");
      }
    } catch (err) {
      setImageError(extractErrorMessage(err));
    } finally {
      setBusyImageId(null);
    }
  }

  if (loadingHotel) {
    return (
      <div className="mx-auto max-w-3xl">
        <div className="h-64 w-full animate-pulse rounded-[16px] bg-line/60" />
      </div>
    );
  }

  if (loadError) {
    return (
      <div className="mx-auto flex max-w-3xl flex-col items-center gap-3 py-16 text-center">
        <span className="flex h-12 w-12 items-center justify-center rounded-full bg-red-50 text-red-500">
          <TriangleAlert size={22} />
        </span>
        <p className="text-sm text-muted">{loadError}</p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl">
      <AdminBreadcrumb
        items={[
          { label: "Hotels", to: "/admin/hotels" },
          { label: isEditing ? form.name || "Edit hotel" : "New hotel" },
        ]}
      />

      <div className="mt-3">
        <AdminPageHeader
          title={isEditing ? "Edit hotel" : "New hotel"}
          subtitle={
            isEditing
              ? "Update this hotel's details and photos."
              : "Save the hotel details first, then add photos."
          }
          actions={
            isEditing && id ? (
              <Link
                to={`/admin/hotels/${id}/rooms`}
                className="flex items-center gap-2 rounded-pill bg-ink px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-ink/85"
              >
                <BedDouble size={15} />
                Manage rooms
              </Link>
            ) : undefined
          }
        />
      </div>

      <form
        onSubmit={handleSubmit}
        className="mt-6 flex flex-col gap-4 rounded-[16px] border border-line bg-white p-6"
      >
        <FormField
          label="Hotel name"
          value={form.name}
          onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
          required
        />

        <div className="flex flex-col gap-1.5">
          <label className="text-sm font-medium text-ink">
            Description <span className="text-xs font-normal text-muted">(optional)</span>
          </label>
          <textarea
            value={form.description}
            onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
            rows={3}
            className="w-full resize-none rounded-[13px] border border-line bg-white px-4 py-3.5 text-sm text-ink outline-none transition-colors placeholder:text-muted focus:border-teal focus:ring-4 focus:ring-teal/10"
          />
        </div>

        <FormField
          label="Address"
          value={form.address}
          onChange={(e) => setForm((f) => ({ ...f, address: e.target.value }))}
          required
        />

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <FormField
            label="City"
            value={form.city}
            onChange={(e) => setForm((f) => ({ ...f, city: e.target.value }))}
            required
          />
          <FormField
            label="Country"
            value={form.country}
            onChange={(e) => setForm((f) => ({ ...f, country: e.target.value }))}
            required
          />
        </div>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <FormField
            label="Phone"
            optional
            value={form.phone}
            onChange={(e) => setForm((f) => ({ ...f, phone: e.target.value }))}
          />
          <FormField
            label="Email"
            type="email"
            optional
            value={form.email}
            onChange={(e) => setForm((f) => ({ ...f, email: e.target.value }))}
          />
        </div>

        <FormSelect
          label="Stars"
          value={form.stars}
          onChange={(e) => setForm((f) => ({ ...f, stars: Number(e.target.value) }))}
        >
          {[1, 2, 3, 4, 5].map((n) => (
            <option key={n} value={n}>
              {n} star{n === 1 ? "" : "s"}
            </option>
          ))}
        </FormSelect>

        {saveError && (
          <p role="alert" className="text-xs font-medium text-red-500">
            {saveError}
          </p>
        )}

        {saved && (
          <p className="flex items-center gap-1.5 text-xs font-medium text-teal">
            <CheckCircle2 size={14} />
            Saved.
          </p>
        )}

        <AuthButton isSubmitting={isSaving} className="mt-1 w-auto px-6">
          {isEditing ? "Save changes" : "Create hotel"}
        </AuthButton>
      </form>

      {isEditing && id && (
        <div className="mt-6 rounded-[16px] border border-line bg-white p-6">
          <h2 className="font-display text-base font-semibold text-ink">Photos</h2>
          <p className="mt-1 text-sm text-muted">
            The main photo is what customers see on hotel cards and search results.
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
                    <img
                      src={image.url}
                      alt=""
                      className="h-full w-full object-cover"
                    />
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
                <Building2 size={20} />
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
                  <img
                    src={p.previewUrl}
                    alt=""
                    className="h-full w-full object-cover opacity-80"
                  />
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

          {imageError && (
            <p role="alert" className="mt-3 text-xs font-medium text-red-500">
              {imageError}
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
        </div>
      )}

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
