// NileStay brand mark: an abstract river-bend motif (three flowing
// vertical currents, evoking the Nile) rather than a literal
// flag/pyramid/emoji. Renders in `currentColor` so it inherits whatever
// foreground color its container sets.
export function LogoMark({ size = 17 }: { size?: number }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={2}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <path d="M4.5 20c2.2-5.3 2.2-10.7 0-16" />
      <path d="M12 20c2.2-5.3 2.2-10.7 0-16" />
      <path d="M19.5 20c2.2-5.3 2.2-10.7 0-16" />
    </svg>
  );
}

const BADGE_SIZES = {
  compact: { badge: "h-8 w-8", icon: 14, text: "text-base" },
  default: { badge: "h-9 w-9", icon: 16, text: "text-lg" },
} as const;

export function Logo({
  size = "default",
  variant = "solid",
  wordmark = true,
  className = "",
}: {
  size?: "compact" | "default";
  variant?: "solid" | "outline";
  wordmark?: boolean;
  className?: string;
}) {
  const { badge, icon, text } = BADGE_SIZES[size];

  return (
    <span className={`flex items-center gap-2.5 ${className}`}>
      <span
        className={`flex ${badge} shrink-0 items-center justify-center rounded-full ${
          variant === "solid"
            ? "bg-teal text-white"
            : "border border-teal/15 bg-white/80 text-teal shadow-sm backdrop-blur-sm"
        }`}
      >
        <LogoMark size={icon} />
      </span>
      {wordmark && (
        <span className={`font-display ${text} font-semibold tracking-tight text-ink`}>
          NileStay
        </span>
      )}
    </span>
  );
}
