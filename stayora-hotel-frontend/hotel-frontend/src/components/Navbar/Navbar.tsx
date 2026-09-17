import { useState } from "react";
import { Link, NavLink, useNavigate } from "react-router-dom";
import {
  Menu,
  X,
  Heart,
  CalendarCheck,
  UserRound,
  LogOut,
  Hotel,
  LayoutDashboard,
} from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { Logo } from "@/components/common/Logo";

const PRIMARY_LINKS = [
  { label: "Hotels", to: "/search", icon: Hotel },
  { label: "Bookings", to: "/my-bookings", icon: CalendarCheck },
  { label: "Favorites", to: "/favorites", icon: Heart },
];

function navLinkClass({ isActive }: { isActive: boolean }) {
  return `rounded-pill px-4 py-2 text-sm font-medium transition-colors ${
    isActive ? "bg-teal/10 text-teal" : "text-muted hover:text-ink"
  }`;
}

export function Navbar() {
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();
  const { isAuthenticated, user, logout } = useAuth();

  // Admin/manager gets exactly one entry point into the separate Admin
  // Dashboard (its own sidebar/layout) rather than admin links mixed into
  // the customer nav - this stays customer-only.
  const isHotelManager = user?.role === "admin" || user?.role === "manager";
  const links = PRIMARY_LINKS;

  return (
    <header className="sticky top-0 z-30 px-4 pt-4 sm:px-6 sm:pt-6">
      <nav className="mx-auto flex max-w-6xl items-center justify-between rounded-pill border border-line bg-white/90 px-4 py-2.5 shadow-[0_8px_30px_rgba(11,63,66,0.06)] backdrop-blur sm:px-6 sm:py-3">
        <Link to="/">
          <Logo />
        </Link>

        {isAuthenticated && (
          <ul className="hidden items-center gap-1 md:flex">
            {links.map(({ label, to, icon: Icon }) => (
              <li key={to}>
                <NavLink to={to} className={navLinkClass}>
                  <span className="flex items-center gap-1.5">
                    <Icon size={15} />
                    {label}
                  </span>
                </NavLink>
              </li>
            ))}
          </ul>
        )}

        <div className="hidden items-center gap-2 md:flex">
          {isAuthenticated ? (
            <>
              {isHotelManager && (
                <Link
                  to="/admin"
                  className="flex items-center gap-1.5 rounded-pill border border-line px-3.5 py-2 text-sm font-medium text-ink transition-colors hover:border-teal/40 hover:text-teal"
                >
                  <LayoutDashboard size={15} />
                  Admin Dashboard
                </Link>
              )}
              <NavLink
                to="/profile"
                className={({ isActive }) =>
                  `flex items-center gap-2 rounded-pill py-1.5 pl-1.5 pr-3.5 text-sm font-medium transition-colors ${
                    isActive ? "bg-teal/10 text-teal" : "text-ink hover:bg-bg"
                  }`
                }
              >
                <span className="flex h-7 w-7 items-center justify-center rounded-full bg-teal/10 text-teal">
                  <UserRound size={14} />
                </span>
                {user?.first_name}
              </NavLink>
              <button
                onClick={logout}
                title="Log out"
                aria-label="Log out"
                className="flex h-9 w-9 items-center justify-center rounded-full text-muted transition-colors hover:bg-bg hover:text-red-500 cursor-pointer"
              >
                <LogOut size={16} />
              </button>
            </>
          ) : (
            <>
              <button
                onClick={() => navigate("/register")}
                className="px-4 py-2.5 text-sm font-semibold text-ink transition-colors hover:text-teal cursor-pointer"
              >
                Sign up
              </button>
              <button
                onClick={() => navigate("/login")}
                className="rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark cursor-pointer"
              >
                Login
              </button>
            </>
          )}
        </div>

        <button
          className="rounded-full p-2 text-ink md:hidden"
          onClick={() => setOpen((v) => !v)}
          aria-label="Toggle menu"
        >
          {open ? <X size={22} /> : <Menu size={22} />}
        </button>
      </nav>

      {open && (
        <div className="mx-auto mt-2 max-w-6xl rounded-card border border-line bg-white p-5 shadow-lg md:hidden">
          {isAuthenticated ? (
            <>
              <div className="flex items-center gap-3 border-b border-line pb-4">
                <span className="flex h-10 w-10 items-center justify-center rounded-full bg-teal/10 text-teal">
                  <UserRound size={18} />
                </span>
                <div>
                  <p className="text-sm font-semibold text-ink">
                    {user?.first_name} {user?.last_name}
                  </p>
                  <p className="text-xs text-muted">{user?.email}</p>
                </div>
              </div>
              <ul className="mt-4 flex flex-col gap-1">
                {[...links, { label: "Profile", to: "/profile", icon: UserRound }].map(
                  ({ label, to, icon: Icon }) => (
                    <li key={to}>
                      <NavLink
                        to={to}
                        onClick={() => setOpen(false)}
                        className={({ isActive }) =>
                          `flex items-center gap-2.5 rounded-[13px] px-3.5 py-3 text-sm font-semibold transition-colors ${
                            isActive
                              ? "bg-teal/10 text-teal"
                              : "text-ink hover:bg-bg"
                          }`
                        }
                      >
                        <Icon size={17} />
                        {label}
                      </NavLink>
                    </li>
                  )
                )}
              </ul>
              {isHotelManager && (
                <Link
                  to="/admin"
                  onClick={() => setOpen(false)}
                  className="mt-3 flex items-center justify-center gap-2 rounded-pill bg-ink py-2.5 text-sm font-semibold text-white"
                >
                  <LayoutDashboard size={15} />
                  Admin Dashboard
                </Link>
              )}
              <button
                onClick={() => {
                  logout();
                  setOpen(false);
                }}
                className="mt-3 flex w-full items-center justify-center gap-2 rounded-pill border border-line py-2.5 text-sm font-semibold text-muted"
              >
                <LogOut size={15} />
                Log out
              </button>
            </>
          ) : (
            <div className="flex flex-col gap-2.5">
              <button
                onClick={() => {
                  navigate("/register");
                  setOpen(false);
                }}
                className="rounded-pill border border-line px-5 py-2.5 text-center text-sm font-semibold text-ink"
              >
                Sign up
              </button>
              <button
                onClick={() => {
                  navigate("/login");
                  setOpen(false);
                }}
                className="rounded-pill bg-teal px-5 py-2.5 text-center text-sm font-semibold text-white"
              >
                Login
              </button>
            </div>
          )}
        </div>
      )}
    </header>
  );
}
