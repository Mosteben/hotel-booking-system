import { useEffect, type ReactNode } from "react";
import { X } from "lucide-react";
import { useMountTransition } from "@/hooks/useMountTransition";

export function Modal({
  open,
  onClose,
  title,
  children,
  maxWidth = "max-w-md",
}: {
  open: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
  maxWidth?: string;
}) {
  const shouldMount = useMountTransition(open);

  useEffect(() => {
    if (!open) return;
    function handleKey(e: KeyboardEvent) {
      if (e.key === "Escape") onClose();
    }
    window.addEventListener("keydown", handleKey);
    return () => window.removeEventListener("keydown", handleKey);
  }, [open, onClose]);

  if (!shouldMount) return null;

  return (
    <div
      className={`fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-ink/40 px-4 py-8 backdrop-blur-sm sm:items-center ${
        open ? "modal-backdrop-in" : "modal-backdrop-out"
      }`}
      onClick={onClose}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="admin-modal-title"
        onClick={(e) => e.stopPropagation()}
        className={`flex max-h-[85vh] w-full ${maxWidth} flex-col rounded-card border border-line bg-white shadow-[var(--shadow-modal)] ${
          open ? "modal-panel-in" : "modal-panel-out"
        }`}
      >
        <div className="flex shrink-0 items-center justify-between gap-3 p-6 pb-0">
          <h2
            id="admin-modal-title"
            className="font-display text-base font-semibold text-ink"
          >
            {title}
          </h2>
          <button
            onClick={onClose}
            aria-label="Close"
            className="rounded-full p-1 text-muted transition-colors hover:text-ink cursor-pointer"
          >
            <X size={16} />
          </button>
        </div>
        {/* Content scrolls internally once it exceeds the panel's own
            max-height, so a tall form/photo grid never pushes its Save/
            Cancel buttons past what's reachable on a short viewport -
            the header above stays put instead of scrolling away too. */}
        <div className="min-h-0 flex-1 overflow-y-auto p-6 pt-4">{children}</div>
      </div>
    </div>
  );
}
