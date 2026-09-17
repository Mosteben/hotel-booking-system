import { Link, useLocation } from "react-router-dom";
import { Bell, ExternalLink, LogOut, Menu } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { useAdminData } from "@/context/AdminDataContext";

const ROUTE_TITLES: { pattern: RegExp; label: string }[] = [
  { pattern: /^\/admin\/?$/, label: "Overview" },
  { pattern: /^\/admin\/bookings/, label: "Bookings" },
  { pattern: /^\/admin\/payments/, label: "Payments" },
  { pattern: /^\/admin\/hotels\/new/, label: "New Hotel" },
  { pattern: /^\/admin\/hotels\/[^/]+\/rooms/, label: "Rooms" },
  { pattern: /^\/admin\/hotels\/[^/]+\/edit/, label: "Edit Hotel" },
  { pattern: /^\/admin\/hotels/, label: "Hotels" },
  { pattern: /^\/admin\/users/, label: "Guests" },
  { pattern: /^\/admin\/reviews/, label: "Reviews" },
];

function useAdminPageTitle(): string {
  const { pathname } = useLocation();
  return ROUTE_TITLES.find((r) => r.pattern.test(pathname))?.label ?? "Admin";
}

export function AdminTopbar({ onMenuClick }: { onMenuClick: () => void }) {
  const { user, logout } = useAuth();
  const { bookings, payments } = useAdminData();
  const pageTitle = useAdminPageTitle();

  // Real counts, reused from the same source as the sidebar badges - not a
  // separate invented notification system.
  const needsAttentionCount =
    bookings.filter((b) => b.status === "pending").length +
    payments.filter((p) => p.status === "pending" || p.status === "failed").length;

  return (
    <header className="sticky top-0 z-20 flex items-center justify-between gap-4 border-b border-line bg-white/95 px-4 py-3.5 backdrop-blur sm:px-6">
      <div className="flex min-w-0 items-center gap-3">
        <button
          onClick={onMenuClick}
          aria-label="Open menu"
          className="shrink-0 rounded-full p-2 text-muted transition-colors hover:bg-bg hover:text-ink md:hidden cursor-pointer"
        >
          <Menu size={20} />
        </button>
        <h1 className="truncate font-display text-base font-semibold text-ink">
          {pageTitle}
        </h1>
        <Link
          to="/"
          className="hidden shrink-0 items-center gap-1.5 rounded-pill border border-line px-2.5 py-1 text-xs font-medium text-muted transition-colors hover:border-teal/40 hover:text-teal lg:flex"
        >
          <ExternalLink size={12} />
          View site
        </Link>
      </div>

      <div className="flex items-center gap-1.5">
        <Link
          to="/admin"
          title="Needs attention"
          aria-label={`${needsAttentionCount} items need attention`}
          className="relative flex h-9 w-9 items-center justify-center rounded-full text-muted transition-colors hover:bg-bg hover:text-ink"
        >
          <Bell size={17} />
          {needsAttentionCount > 0 && (
            <span className="absolute right-1 top-1 flex h-4 min-w-4 items-center justify-center rounded-pill bg-red-500 px-1 text-[10px] font-semibold text-white">
              {needsAttentionCount}
            </span>
          )}
        </Link>

        <div className="mx-1.5 hidden h-6 w-px bg-line sm:block" />

        <div className="hidden items-center gap-2 sm:flex">
          <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-teal/10 text-xs font-semibold text-teal">
            {user?.first_name?.[0]}
            {user?.last_name?.[0]}
          </span>
          <div className="leading-tight">
            <p className="text-sm font-semibold text-ink">
              {user?.first_name} {user?.last_name}
            </p>
            <p className="text-xs capitalize text-muted">{user?.role}</p>
          </div>
        </div>

        <button
          onClick={logout}
          title="Log out"
          aria-label="Log out"
          className="flex h-9 w-9 items-center justify-center rounded-full text-muted transition-colors hover:bg-bg hover:text-red-500 cursor-pointer"
        >
          <LogOut size={16} />
        </button>
      </div>
    </header>
  );
}
