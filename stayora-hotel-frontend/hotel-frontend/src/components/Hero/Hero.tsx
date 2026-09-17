import { SearchBar } from "@/components/SearchBar/SearchBar";

export function Hero() {
  return (
    <section className="relative overflow-hidden px-4 pb-4 pt-10 sm:px-6 sm:pt-14 lg:pt-16">
      <div className="auth-blob pointer-events-none absolute -left-32 -top-32 h-96 w-96 rounded-full bg-teal/10 blur-3xl" />
      <div className="auth-blob auth-blob-delay pointer-events-none absolute -right-24 top-10 h-80 w-80 rounded-full bg-teal-light/10 blur-3xl" />

      <div className="relative mx-auto max-w-6xl">
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-sm font-semibold tracking-wide text-teal">
            Real stays, verified by NileStay
          </p>
          <h1 className="mt-2 font-display text-4xl font-bold leading-[1.1] tracking-tight text-ink sm:text-5xl lg:text-[3.4rem]">
            Stay somewhere
            <br />
            worth remembering
          </h1>
          <p className="mx-auto mt-5 max-w-md text-base leading-relaxed text-muted sm:text-lg">
            Search, compare and book hotels with transparent pricing - every
            listing pulled straight from our live catalog.
          </p>
        </div>

        <div className="mt-10 sm:mt-12">
          <SearchBar />
        </div>
      </div>
    </section>
  );
}
