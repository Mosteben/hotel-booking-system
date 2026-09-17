import { Logo } from "@/components/common/Logo";

export function Footer() {
  return (
    <footer className="border-t border-line px-4 py-10 sm:px-6">
      <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 sm:flex-row">
        <Logo size="compact" />
        <p className="text-sm text-muted">
          © {new Date().getFullYear()} NileStay. All rights reserved.
        </p>
      </div>
    </footer>
  );
}
