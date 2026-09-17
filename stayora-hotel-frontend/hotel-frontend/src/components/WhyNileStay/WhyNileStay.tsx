import { ShieldCheck, BadgeDollarSign, Clock, Heart } from "lucide-react";

// Static marketing copy describing the platform itself, not business/hotel
// data - explicitly allowed as UI content. No hotel names, prices, or
// counts appear here.
const BENEFITS = [
  {
    icon: ShieldCheck,
    title: "Verified listings",
    description: "Every hotel and room comes straight from our live catalog.",
  },
  {
    icon: BadgeDollarSign,
    title: "Transparent pricing",
    description: "The price you see is calculated by our system, not guessed.",
  },
  {
    icon: Clock,
    title: "Instant confirmation",
    description: "Your booking status updates the moment we hear back.",
  },
  {
    icon: Heart,
    title: "Save your favorites",
    description: "Keep track of the stays you love and come back anytime.",
  },
];

export function WhyNileStay() {
  return (
    <section className="bg-cream px-4 py-20 sm:px-6 sm:py-24">
      <div className="mx-auto max-w-6xl">
        <div className="mx-auto max-w-xl text-center">
          <p className="text-sm font-semibold tracking-wide text-teal">
            Why NileStay
          </p>
          <h2 className="mt-1.5 font-display text-2xl font-bold text-ink sm:text-3xl">
            Booking made simple
          </h2>
        </div>

        <div className="mt-12 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {BENEFITS.map(({ icon: Icon, title, description }) => (
            <div
              key={title}
              className="flex flex-col items-start gap-3.5 rounded-card bg-white p-7 shadow-[0_8px_24px_rgba(16,24,40,0.06)]"
            >
              <span className="flex h-11 w-11 items-center justify-center rounded-full bg-teal/10 text-teal">
                <Icon size={20} />
              </span>
              <h3 className="font-display text-base font-semibold text-ink">
                {title}
              </h3>
              <p className="text-sm leading-relaxed text-muted">
                {description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
