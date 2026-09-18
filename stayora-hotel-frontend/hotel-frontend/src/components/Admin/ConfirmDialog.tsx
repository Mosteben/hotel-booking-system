import { useEffect, useRef } from "react";
import { TriangleAlert, X } from "lucide-react";
import type { ConfirmOptions } from "@/hooks/useConfirmDialog";
import { useMountTransition } from "@/hooks/useMountTransition";

interface ConfirmDialogProps extends ConfirmOptions {
  open: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export function ConfirmDialog({
  open,
  title,
  description,
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  destructive = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  const confirmRef = useRef<HTMLButtonElement>(null);
  const shouldMount = useMountTransition(open);

  useEffect(() => {
    if (!open) return;
    confirmRef.current?.focus();

    function handleKey(e: KeyboardEvent) {
      if (e.key === "Escape") onCancel();
    }
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [open, onCancel]);

  if (!shouldMount) return null;

  return (
    <div
      className={`fixed inset-0 z-50 flex items-center justify-center bg-ink/40 px-4 backdrop-blur-sm ${
        open ? "modal-backdrop-in" : "modal-backdrop-out"
      }`}
      onClick={onCancel}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-dialog-title"
        onClick={(e) => e.stopPropagation()}
        className={`w-full max-w-sm rounded-card border border-line bg-white p-6 shadow-[var(--shadow-modal)] ${
          open ? "modal-panel-in" : "modal-panel-out"
        }`}
      >
        <div className="flex items-start gap-3">
          <span
            className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-full ${
              destructive ? "bg-error-soft text-error" : "bg-teal/10 text-teal"
            }`}
          >
            <TriangleAlert size={18} />
          </span>
          <div className="min-w-0 flex-1 pt-1">
            <h2
              id="confirm-dialog-title"
              className="font-display text-base font-semibold text-ink"
            >
              {title}
            </h2>
            <p className="mt-1.5 text-sm leading-relaxed text-muted">
              {description}
            </p>
          </div>
          <button
            onClick={onCancel}
            aria-label="Close"
            className="shrink-0 rounded-full p-1 text-muted transition-colors hover:text-ink cursor-pointer"
          >
            <X size={16} />
          </button>
        </div>

        <div className="mt-6 flex justify-end gap-2.5">
          <button
            onClick={onCancel}
            className="rounded-pill border border-line px-4 py-2.5 text-sm font-semibold text-ink transition-colors hover:border-teal/40 cursor-pointer"
          >
            {cancelLabel}
          </button>
          <button
            ref={confirmRef}
            onClick={onConfirm}
            className={`rounded-pill px-4 py-2.5 text-sm font-semibold text-white transition-all duration-150 active:scale-[0.97] cursor-pointer ${
              destructive ? "bg-error hover:bg-red-700" : "bg-teal hover:bg-teal-dark"
            }`}
          >
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
