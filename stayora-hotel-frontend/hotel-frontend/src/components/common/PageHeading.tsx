// Mirrors the Login/Register page heading pattern exactly (eyebrow label +
// display heading + muted subtitle) so every page in the app reads as part
// of the same visual product.

export function PageHeading({
  eyebrow,
  title,
  subtitle,
}: {
  eyebrow: string;
  title: string;
  subtitle?: string;
}) {
  return (
    <div>
      <p className="text-sm font-semibold tracking-wide text-teal">
        {eyebrow}
      </p>
      <h1 className="mt-1.5 font-display text-2xl font-bold text-ink sm:text-3xl">
        {title}
      </h1>
      {subtitle && (
        <p className="mt-2 text-sm leading-relaxed text-muted">{subtitle}</p>
      )}
    </div>
  );
}
