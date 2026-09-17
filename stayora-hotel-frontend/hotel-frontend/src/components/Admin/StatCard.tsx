import type { LucideIcon } from "lucide-react";

type Tone = "default" | "warning" | "danger" | "success";

const TONE_STYLES: Record<Tone, string> = {
  default: "bg-cream text-teal",
  warning: "bg-amber-50 text-amber-600",
  danger: "bg-red-50 text-red-500",
  success: "bg-emerald-50 text-emerald-600",
};

export function StatCard({
  label,
  value,
  icon: Icon,
  tone = "default",
}: {
  label: string;
  value: string | number;
  icon: LucideIcon;
  tone?: Tone;
}) {
  return (
    <div className="flex items-center gap-4 rounded-[16px] border border-line bg-white p-5">
      <span
        className={`flex h-11 w-11 shrink-0 items-center justify-center rounded-[12px] ${TONE_STYLES[tone]}`}
      >
        <Icon size={18} />
      </span>
      <div className="min-w-0">
        <p className="text-xs font-medium uppercase tracking-wide text-muted">
          {label}
        </p>
        <p className="mt-0.5 font-display text-2xl font-bold text-ink">{value}</p>
      </div>
    </div>
  );
}
