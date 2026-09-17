import { Check } from "lucide-react";

const STEPS = ["Dates & guests", "Summary", "Payment"];

export function BookingProgress({ current }: { current: 1 | 2 | 3 }) {
  return (
    <ol className="flex items-center gap-2 sm:gap-3">
      {STEPS.map((label, index) => {
        const step = (index + 1) as 1 | 2 | 3;
        const isDone = step < current;
        const isActive = step === current;
        return (
          <li key={label} className="flex flex-1 items-center gap-2 sm:gap-3">
            <div className="flex items-center gap-2">
              <span
                className={`flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold transition-colors ${
                  isDone
                    ? "bg-teal text-white"
                    : isActive
                      ? "bg-teal text-white"
                      : "bg-line text-muted"
                }`}
              >
                {isDone ? <Check size={14} /> : step}
              </span>
              <span
                className={`hidden text-xs font-semibold sm:inline ${
                  isActive || isDone ? "text-ink" : "text-muted"
                }`}
              >
                {label}
              </span>
            </div>
            {step < 3 && (
              <span
                className={`h-px flex-1 ${isDone ? "bg-teal" : "bg-line"}`}
              />
            )}
          </li>
        );
      })}
    </ol>
  );
}
