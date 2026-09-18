import type { ReactNode } from "react";
import { Link } from "react-router-dom";
import { ArrowLeft } from "lucide-react";
import { Logo } from "@/components/common/Logo";

interface AuthLayoutProps {
  imageSide?: "left" | "right";
  children: ReactNode;
}

function BrandImagePanel({ order }: { order: "lg:order-1" | "lg:order-2" }) {
  return (
    <div className={`relative hidden h-full overflow-hidden lg:block ${order}`}>
      <img
        src="/images/hotel-illustration.png"
        alt=""
        aria-hidden="true"
        className="auth-image-in absolute inset-0 h-full w-full object-cover"
      />
      {/* Light, airy overlay (was a dark navy gradient) so the panel reads
          as bright and premium instead of moody. Text below switches to
          dark ink to stay readable against it. */}
      <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-white/92 via-white/55 to-white/15" />
      <div className="pointer-events-none absolute inset-0 bg-gradient-to-tr from-teal/10 via-transparent to-transparent" />

      <div className="relative flex h-full flex-col justify-between p-10 xl:p-14">
        <Link to="/" className="inline-flex w-fit">
          <Logo variant="outline" />
        </Link>

        <div className="max-w-sm">
          <h2 className="font-display text-3xl font-bold leading-[1.15] text-ink xl:text-[2.25rem]">
            Stay somewhere unforgettable.
          </h2>
          <p className="mt-3 text-sm leading-relaxed text-muted">
            Discover exceptional stays, thoughtfully selected for your next
            journey.
          </p>
        </div>
      </div>
    </div>
  );
}

export function AuthLayout({ imageSide = "left", children }: AuthLayoutProps) {
  const imageOrder = imageSide === "right" ? "lg:order-2" : "lg:order-1";
  const formOrder = imageSide === "right" ? "lg:order-1" : "lg:order-2";

  return (
    // Fixed viewport height + overflow-hidden: the page itself never grows,
    // no matter how long the form inside gets.
    <div className="relative flex h-screen w-full items-center justify-center overflow-hidden bg-white p-4 sm:p-6 lg:p-8">
      {/* Soft floating blue blobs behind the card — purely decorative,
          keeps the pure-white background from feeling flat. */}
      <div className="auth-blob pointer-events-none absolute -left-24 -top-24 h-96 w-96 rounded-full bg-teal/10 blur-3xl" />
      <div className="auth-blob auth-blob-delay pointer-events-none absolute -bottom-24 -right-16 h-[28rem] w-[28rem] rounded-full bg-teal-light/15 blur-3xl" />

      {/* The card is a fixed 85vh block, sized up and with a hover lift.
          The image panel lives inside it, so it's bound by this fixed
          height and can never stretch or get dragged along by a tall
          form anymore. */}
      <div className="auth-card-in relative z-10 grid h-[85vh] w-full max-w-[1360px] grid-cols-1 overflow-hidden rounded-panel border border-line bg-white shadow-[var(--shadow-modal)] lg:grid-cols-2">
        <BrandImagePanel order={imageOrder} />

        <div className={`flex h-full min-h-0 flex-col ${formOrder}`}>
          {/* Minimal header: logo (mobile only, since the image panel already
              carries it on desktop) + a quiet back-to-home link. */}
          <div className="flex shrink-0 items-center justify-between px-6 py-5 sm:px-10 lg:justify-end lg:px-12">
            <Link to="/" className="lg:hidden">
              <Logo size="compact" />
            </Link>

            <Link
              to="/"
              className="group inline-flex items-center gap-1.5 text-sm font-medium text-muted transition-colors duration-200 hover:text-ink"
            >
              <ArrowLeft
                size={15}
                className="transition-transform duration-200 group-hover:-translate-x-0.5"
              />
              Back to home
            </Link>
          </div>

          {/* min-h-0 is what lets this area actually scroll within the
              fixed-height card instead of pushing the card taller. */}
          <div className="auth-scroll min-h-0 flex-1 overflow-y-auto px-6 pb-10 sm:px-10">
            <div className="auth-fade-in mx-auto w-full max-w-[420px] py-4">
              {children}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}