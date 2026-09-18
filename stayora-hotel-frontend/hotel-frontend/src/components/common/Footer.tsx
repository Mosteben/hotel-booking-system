import { Link } from "react-router-dom";
import { Logo } from "@/components/common/Logo";
import { useAuth } from "@/context/AuthContext";

// Every link below points to a route that actually exists in App.tsx - no
// invented pages (no "About"/"Careers"/"Blog" placeholders). Account
// column adapts to auth state the same way Navbar does.
const EXPLORE_LINKS = [
  { label: "Home", to: "/" },
  { label: "Search hotels", to: "/search" },
];

export function Footer() {
  const { isAuthenticated } = useAuth();

  const accountLinks = isAuthenticated
    ? [
        { label: "My bookings", to: "/my-bookings" },
        { label: "Favorites", to: "/favorites" },
        { label: "Profile", to: "/profile" },
      ]
    : [
        { label: "Login", to: "/login" },
        { label: "Create account", to: "/register" },
      ];

  return (
    <footer className="border-t border-line bg-white">
      <div className="mx-auto max-w-6xl px-4 py-14 sm:px-6">
        <div className="grid grid-cols-2 gap-10 sm:grid-cols-3">
          <div className="col-span-2 sm:col-span-1">
            <Logo />
            <p className="mt-4 max-w-[220px] text-sm leading-relaxed text-muted">
              Real hotels, transparent pricing, and instant confirmation
              along the Nile and beyond.
            </p>
          </div>

          <div>
            <p className="text-xs font-semibold uppercase tracking-wider text-muted">
              Explore
            </p>
            <ul className="mt-4 flex flex-col gap-3">
              {EXPLORE_LINKS.map((link) => (
                <li key={link.to}>
                  <Link
                    to={link.to}
                    className="text-sm text-ink transition-colors hover:text-teal"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <p className="text-xs font-semibold uppercase tracking-wider text-muted">
              Account
            </p>
            <ul className="mt-4 flex flex-col gap-3">
              {accountLinks.map((link) => (
                <li key={link.to}>
                  <Link
                    to={link.to}
                    className="text-sm text-ink transition-colors hover:text-teal"
                  >
                    {link.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="mt-12 border-t border-line pt-6 text-center sm:text-left">
          <p className="text-sm text-muted">
            © {new Date().getFullYear()} NileStay. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
