import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  ClipboardList,
  Wallet,
  Building2,
  Users,
  MessageSquareText,
  X,
  type LucideIcon,
} from "lucide-react";
import { LogoMark } from "@/components/common/Logo";
import { useAdminData } from "@/context/AdminDataContext";

interface NavItem {
  label: string;
  to: string;
  icon: LucideIcon;
  badge?: number;
}

interface NavSection {
  label?: string;
  items: NavItem[];
}

function SidebarContent({ onNavigate }: { onNavigate: () => void }) {
  const { bookings, payments } = useAdminData();

  // Only real, unambiguous "needs a human" counts - not invented metrics.
  const pendingBookings = bookings.filter((b) => b.status === "pending").length;
  const attentionPayments = payments.filter(
    (p) => p.status === "pending" || p.status === "failed"
  ).length;

  const sections: NavSection[] = [
    { items: [{ label: "Overview", to: "/admin", icon: LayoutDashboard }] },
    {
      label: "Reservations",
      items: [
        {
          label: "Bookings",
          to: "/admin/bookings",
          icon: ClipboardList,
          badge: pendingBookings,
        },
        {
          label: "Payments",
          to: "/admin/payments",
          icon: Wallet,
          badge: attentionPayments,
        },
      ],
    },
    {
      label: "Property",
      items: [{ label: "Hotels", to: "/admin/hotels", icon: Building2 }],
    },
    {
      label: "Customers",
      items: [{ label: "Guests", to: "/admin/users", icon: Users }],
    },
    {
      label: "Content",
      items: [
        { label: "Reviews", to: "/admin/reviews", icon: MessageSquareText },
      ],
    },
  ];

  return (
    <nav className="flex h-full flex-col overflow-y-auto pb-6">
      <div className="flex items-center gap-2.5 px-5 py-5">
        <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-teal text-white">
          <LogoMark size={16} />
        </span>
        <div>
          <p className="font-display text-sm font-semibold text-admin-sidebar-text">
            NileStay
          </p>
          <p className="text-[11px] font-medium text-admin-sidebar-text-secondary">
            Admin
          </p>
        </div>
      </div>

      <div className="mt-2 flex flex-1 flex-col gap-5 px-3">
        {sections.map((section, i) => (
          <div key={i}>
            {section.label && (
              <p className="px-3 pb-1.5 text-[11px] font-semibold uppercase tracking-wider text-admin-sidebar-text-secondary/70">
                {section.label}
              </p>
            )}
            <ul className="flex flex-col gap-0.5">
              {section.items.map(({ label, to, icon: Icon, badge }) => (
                <li key={to}>
                  <NavLink
                    to={to}
                    end={to === "/admin"}
                    onClick={onNavigate}
                    className={({ isActive }) =>
                      `flex items-center justify-between gap-2 rounded-[10px] px-3 py-2.5 text-sm font-medium transition-colors duration-150 ${
                        isActive
                          ? "bg-admin-sidebar-active text-admin-sidebar-text hover:bg-admin-sidebar-active-hover"
                          : "text-admin-sidebar-text-secondary hover:bg-admin-sidebar-hover hover:text-admin-sidebar-text"
                      }`
                    }
                  >
                    <span className="flex items-center gap-2.5">
                      <Icon size={16} />
                      {label}
                    </span>
                    {!!badge && badge > 0 && (
                      <span className="flex h-5 min-w-5 items-center justify-center rounded-pill bg-white/25 px-1.5 text-[11px] font-semibold text-admin-sidebar-text">
                        {badge}
                      </span>
                    )}
                  </NavLink>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </div>
    </nav>
  );
}

export function AdminSidebar({
  mobileOpen,
  onClose,
}: {
  mobileOpen: boolean;
  onClose: () => void;
}) {
  return (
    <>
      {/* Desktop: permanent sidebar */}
      <aside className="hidden md:block md:w-64 md:shrink-0 md:bg-admin-sidebar-bg">
        <div className="md:fixed md:h-screen md:w-64">
          <SidebarContent onNavigate={() => {}} />
        </div>
      </aside>

      {/* Mobile: slide-in drawer */}
      {mobileOpen && (
        <div className="fixed inset-0 z-40 md:hidden">
          <div
            className="absolute inset-0 bg-ink/50 backdrop-blur-sm"
            onClick={onClose}
          />
          <div className="absolute inset-y-0 left-0 w-64 bg-admin-sidebar-bg shadow-2xl">
            <button
              onClick={onClose}
              aria-label="Close menu"
              className="absolute right-3 top-5 rounded-full p-1.5 text-admin-sidebar-text-secondary transition-colors hover:text-admin-sidebar-text cursor-pointer"
            >
              <X size={18} />
            </button>
            <SidebarContent onNavigate={onClose} />
          </div>
        </div>
      )}
    </>
  );
}
