import { useMemo } from "react";
import { Link } from "react-router-dom";
import {
  Banknote,
  Building2,
  CalendarClock,
  CheckCircle2,
  ClipboardList,
  ImageOff,
  ShieldCheck,
  Users,
  Wallet,
  XCircle,
} from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { useAdminData } from "@/context/AdminDataContext";
import { StatCard } from "@/components/Admin/StatCard";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import { AdminErrorState, AdminLoadingRows } from "@/components/Admin/AdminStates";
import { StatusBadge } from "@/components/common/StatusBadge";

type Severity = "critical" | "warning" | "info";

interface AttentionItem {
  key: string;
  icon: typeof ClipboardList;
  title: string;
  description: string;
  actionLabel: string;
  to: string;
  severity: Severity;
}

const SEVERITY_STYLES: Record<Severity, { card: string; icon: string }> = {
  critical: { card: "border-red-100 bg-red-50/60", icon: "bg-red-100 text-red-600" },
  warning: { card: "border-amber-100 bg-amber-50/60", icon: "bg-amber-100 text-amber-700" },
  info: { card: "border-line bg-bg", icon: "bg-cream text-teal" },
};

const RESOURCE_LABELS: Record<string, string> = {
  bookings: "bookings",
  payments: "payments",
  reviews: "reviews",
  users: "guests",
  hotels: "hotels",
};

export function Overview() {
  const { user } = useAuth();
  const { bookings, payments, hotels, users, isLoading, errors, refresh } =
    useAdminData();

  const failedResources = Object.entries(errors).filter(([, msg]) => Boolean(msg));

  const stats = useMemo(() => {
    const pendingBookings = bookings.filter((b) => b.status === "pending").length;
    const confirmedBookings = bookings.filter((b) => b.status === "confirmed").length;
    const cancelledBookings = bookings.filter((b) => b.status === "cancelled").length;

    const paidPayments = payments.filter((p) => p.status === "paid").length;
    const pendingPayments = payments.filter((p) => p.status === "pending").length;
    const failedPayments = payments.filter((p) => p.status === "failed").length;
    const refundedPayments = payments.filter((p) => p.status === "refunded").length;

    const hotelsWithoutRooms = hotels.filter((h) => h.roomCount === 0);
    const hotelsWithoutImage = hotels.filter((h) => !h.image_url);

    return {
      totalBookings: bookings.length,
      pendingBookings,
      confirmedBookings,
      cancelledBookings,
      totalPayments: payments.length,
      paidPayments,
      pendingPayments,
      failedPayments,
      refundedPayments,
      totalHotels: hotels.length,
      totalUsers: users.length,
      hotelsWithoutRooms,
      hotelsWithoutImage,
    };
  }, [bookings, payments, hotels, users]);

  const recentBookings = useMemo(
    () =>
      [...bookings]
        .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
        .slice(0, 5),
    [bookings]
  );

  const recentPayments = useMemo(
    () =>
      [...payments]
        .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
        .slice(0, 5),
    [payments]
  );

  const attentionItems: AttentionItem[] = [];

  // Failed payments are the only truly urgent condition here - real money
  // didn't go through. Everything else is routine/time-sensitive, not an
  // emergency, so it gets a calmer treatment.
  if (stats.failedPayments > 0) {
    attentionItems.push({
      key: "failed-payments",
      icon: XCircle,
      title: `${stats.failedPayments} failed payment${stats.failedPayments === 1 ? "" : "s"}`,
      description: "These payments didn't go through and may need follow-up.",
      actionLabel: "Review payments",
      to: "/admin/payments",
      severity: "critical",
    });
  }

  if (stats.pendingBookings > 0) {
    attentionItems.push({
      key: "pending-bookings",
      icon: CalendarClock,
      title: `${stats.pendingBookings} booking${stats.pendingBookings === 1 ? "" : "s"} awaiting confirmation`,
      description: "New bookings are pending until an admin confirms or cancels them.",
      actionLabel: "Review bookings",
      to: "/admin/bookings",
      severity: "warning",
    });
  }

  if (stats.pendingPayments > 0) {
    attentionItems.push({
      key: "pending-payments",
      icon: Wallet,
      title: `${stats.pendingPayments} payment${stats.pendingPayments === 1 ? "" : "s"} pending`,
      description: "Recorded but not yet confirmed as paid.",
      actionLabel: "Review payments",
      to: "/admin/payments",
      severity: "warning",
    });
  }

  if (stats.hotelsWithoutRooms.length > 0) {
    attentionItems.push({
      key: "hotels-no-rooms",
      icon: Building2,
      title: `${stats.hotelsWithoutRooms.length} hotel${stats.hotelsWithoutRooms.length === 1 ? "" : "s"} with no rooms`,
      description: `${stats.hotelsWithoutRooms
        .slice(0, 3)
        .map((h) => h.name)
        .join(", ")}${stats.hotelsWithoutRooms.length > 3 ? ", ..." : ""} can't be booked until rooms are added.`,
      actionLabel: "Manage hotels",
      to: "/admin/hotels",
      severity: "info",
    });
  }

  if (stats.hotelsWithoutImage.length > 0) {
    attentionItems.push({
      key: "hotels-no-image",
      icon: ImageOff,
      title: `${stats.hotelsWithoutImage.length} hotel${stats.hotelsWithoutImage.length === 1 ? "" : "s"} with no photo`,
      description: `${stats.hotelsWithoutImage
        .slice(0, 3)
        .map((h) => h.name)
        .join(", ")}${stats.hotelsWithoutImage.length > 3 ? ", ..." : ""} show a placeholder to guests.`,
      actionLabel: "Manage hotels",
      to: "/admin/hotels",
      severity: "info",
    });
  }

  if (isLoading) {
    return (
      <div className="mx-auto max-w-6xl">
        <AdminPageHeader title="Overview" />
        <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="h-24 animate-pulse rounded-[16px] bg-line/60" />
          ))}
        </div>
        <div className="mt-8">
          <AdminLoadingRows count={3} />
        </div>
      </div>
    );
  }

  // Only block the whole dashboard if every single resource failed to
  // load - a genuine "nothing to show" state. If some loaded and some
  // didn't, show what's real and flag the rest below instead of hiding
  // working data behind one failed request.
  if (failedResources.length === Object.keys(RESOURCE_LABELS).length) {
    return (
      <div className="mx-auto max-w-6xl">
        <AdminPageHeader title="Overview" />
        <div className="mt-6">
          <AdminErrorState
            message={`Failed to load the dashboard: ${errors.bookings}`}
            onRetry={refresh}
          />
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-6xl">
      <AdminPageHeader
        title={user?.first_name ? `Welcome back, ${user.first_name}` : "Overview"}
        subtitle="A real-time snapshot of NileStay's operations."
      />

      {failedResources.length > 0 && (
        <p
          role="alert"
          className="mt-4 rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
        >
          Couldn't load{" "}
          {failedResources.map(([key]) => RESOURCE_LABELS[key] ?? key).join(", ")}. Some
          numbers below may be incomplete.{" "}
          <button onClick={refresh} className="underline underline-offset-2 cursor-pointer">
            Retry
          </button>
        </p>
      )}

      <div className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        <StatCard label="Total bookings" value={stats.totalBookings} icon={ClipboardList} />
        <StatCard
          label="Pending bookings"
          value={stats.pendingBookings}
          icon={CalendarClock}
          tone={stats.pendingBookings > 0 ? "warning" : "default"}
        />
        <StatCard label="Confirmed" value={stats.confirmedBookings} icon={CheckCircle2} tone="success" />
        <StatCard label="Cancelled" value={stats.cancelledBookings} icon={XCircle} />

        <StatCard label="Paid payments" value={stats.paidPayments} icon={Banknote} tone="success" />
        <StatCard
          label="Failed payments"
          value={stats.failedPayments}
          icon={XCircle}
          tone={stats.failedPayments > 0 ? "danger" : "default"}
        />
        <StatCard label="Refunded" value={stats.refundedPayments} icon={Wallet} />
        <StatCard label="Total hotels" value={stats.totalHotels} icon={Building2} />

        <StatCard label="Total users" value={stats.totalUsers} icon={Users} />
      </div>

      <div className="mt-10">
        <h2 className="font-display text-lg font-semibold text-ink">Needs attention</h2>

        {attentionItems.length === 0 ? (
          <div className="mt-4 flex flex-col items-center gap-2 rounded-[16px] border border-dashed border-line bg-white px-6 py-14 text-center">
            <span className="flex h-11 w-11 items-center justify-center rounded-full bg-emerald-50 text-emerald-600">
              <ShieldCheck size={20} />
            </span>
            <p className="text-sm font-semibold text-ink">All caught up</p>
            <p className="text-sm text-muted">
              Nothing needs your attention right now.
            </p>
          </div>
        ) : (
          <div className="mt-4 flex flex-col gap-2.5">
            {attentionItems.map(({ key, icon: Icon, title, description, actionLabel, to, severity }) => {
              const style = SEVERITY_STYLES[severity];
              return (
                <div
                  key={key}
                  className={`flex flex-col gap-3 rounded-[14px] border p-4 sm:flex-row sm:items-center sm:justify-between ${style.card}`}
                >
                  <div className="flex items-start gap-3">
                    <span className={`flex h-9 w-9 shrink-0 items-center justify-center rounded-full ${style.icon}`}>
                      <Icon size={16} />
                    </span>
                    <div className="min-w-0">
                      <p className="text-sm font-semibold text-ink">{title}</p>
                      <p className="mt-0.5 text-sm text-muted">{description}</p>
                    </div>
                  </div>
                  <Link
                    to={to}
                    className="shrink-0 self-start rounded-pill bg-white px-4 py-2 text-xs font-semibold text-ink shadow-sm ring-1 ring-line transition-colors hover:bg-ink hover:text-white hover:ring-ink sm:self-auto"
                  >
                    {actionLabel}
                  </Link>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <div className="mt-10 grid grid-cols-1 gap-6 lg:grid-cols-2">
        <div>
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold text-ink">Recent bookings</h2>
            <Link to="/admin/bookings" className="text-xs font-semibold text-teal hover:text-teal-dark">
              View all
            </Link>
          </div>
          {recentBookings.length === 0 ? (
            <p className="mt-4 rounded-[14px] border border-dashed border-line bg-white px-4 py-8 text-center text-sm text-muted">
              No bookings yet.
            </p>
          ) : (
            <div className="mt-4 flex flex-col gap-2">
              {recentBookings.map((booking) => (
                <div
                  key={booking.id}
                  className="flex items-center justify-between gap-3 rounded-[12px] border border-line bg-white px-4 py-3"
                >
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-ink">Booking #{booking.id}</p>
                    <p className="text-xs text-muted">
                      {new Date(booking.created_at).toLocaleDateString()} · $
                      {booking.total_price.toFixed(2)}
                    </p>
                  </div>
                  <StatusBadge status={booking.status} />
                </div>
              ))}
            </div>
          )}
        </div>

        <div>
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold text-ink">Recent payments</h2>
            <Link to="/admin/payments" className="text-xs font-semibold text-teal hover:text-teal-dark">
              View all
            </Link>
          </div>
          {recentPayments.length === 0 ? (
            <p className="mt-4 rounded-[14px] border border-dashed border-line bg-white px-4 py-8 text-center text-sm text-muted">
              No payments yet.
            </p>
          ) : (
            <div className="mt-4 flex flex-col gap-2">
              {recentPayments.map((payment) => (
                <div
                  key={payment.id}
                  className="flex items-center justify-between gap-3 rounded-[12px] border border-line bg-white px-4 py-3"
                >
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-ink">${payment.amount.toFixed(2)}</p>
                    <p className="text-xs capitalize text-muted">
                      {payment.payment_method} · {new Date(payment.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <StatusBadge status={payment.status} />
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
