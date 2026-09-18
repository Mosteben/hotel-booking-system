import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  BedDouble,
  CalendarDays,
  CalendarX,
} from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { StatusBadge } from "@/components/common/StatusBadge";
import { PageHeading } from "@/components/common/PageHeading";
import { EmptyState } from "@/components/common/EmptyState";
import { ErrorState } from "@/components/common/ErrorState";
import { Skeleton } from "@/components/common/Skeleton";
import { getMyBookings, cancelBooking } from "@/api/bookingApi";
import { getRoomByID } from "@/api/roomApi";
import { getHotelByID } from "@/api/hotelApi";
import { getMyPayments } from "@/api/paymentApi";
import { extractErrorMessage } from "@/api/client";
import { getHotelImageUrl } from "@/utils/hotelImages";
import type { Booking } from "@/types/booking";
import type { Room } from "@/types/room";
import type { Hotel } from "@/types/hotel";
import type { Payment } from "@/types/payment";

type Status = "loading" | "success" | "empty" | "error";
type Group = "upcoming" | "completed" | "cancelled";

interface EnrichedBooking {
  booking: Booking;
  room: Room | null;
  hotel: Hotel | null;
  payment: Payment | null;
}

function groupOf(booking: Booking): Group {
  if (booking.status === "cancelled") return "cancelled";
  const checkOut = new Date(booking.check_out).getTime();
  return checkOut >= Date.now() ? "upcoming" : "completed";
}

function BookingCardSkeleton() {
  return (
    <div className="flex flex-col gap-6 rounded-card border border-line bg-white p-7 sm:p-8 lg:grid lg:grid-cols-[160px_1fr_200px] lg:gap-8">
      <Skeleton className="h-40 w-full rounded-[16px] lg:h-auto" />
      <div className="flex flex-col gap-4 lg:border-l lg:border-line lg:pl-8">
        <div className="flex flex-col gap-2">
          <Skeleton className="h-5 w-48" />
          <Skeleton className="h-3 w-32" />
        </div>
        <div className="flex flex-wrap gap-6">
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-10 w-24" />
          <Skeleton className="h-10 w-20" />
        </div>
      </div>
      <div className="flex flex-col justify-between gap-4 lg:border-l lg:border-line lg:pl-8">
        <Skeleton className="h-8 w-24" />
        <Skeleton className="h-10 w-full rounded-pill" />
      </div>
    </div>
  );
}

export function MyBookings() {
  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [items, setItems] = useState<EnrichedBooking[]>([]);
  const [activeGroup, setActiveGroup] = useState<Group>("upcoming");
  const [cancellingId, setCancellingId] = useState<number | null>(null);

  async function load() {
    setStatus("loading");
    try {
      const bookings = await getMyBookings();

      if (bookings.length === 0) {
        setItems([]);
        setStatus("empty");
        return;
      }

      const roomCache = new Map<number, Room | null>();
      const hotelCache = new Map<number, Hotel | null>();

      let payments: Payment[] = [];
      try {
        const paymentsRes = await getMyPayments();
        if (paymentsRes.success) payments = paymentsRes.data;
      } catch {
        // Payment status is supplementary - don't fail the whole page for it.
      }

      const enriched = await Promise.all(
        bookings.map(async (booking): Promise<EnrichedBooking> => {
          let room = roomCache.get(booking.room_id) ?? null;
          if (!roomCache.has(booking.room_id)) {
            try {
              const roomRes = await getRoomByID(booking.room_id);
              room = roomRes.success ? roomRes.data : null;
            } catch {
              room = null;
            }
            roomCache.set(booking.room_id, room);
          }

          let hotel: Hotel | null = null;
          if (room) {
            hotel = hotelCache.get(room.hotel_id) ?? null;
            if (!hotelCache.has(room.hotel_id)) {
              try {
                const hotelRes = await getHotelByID(room.hotel_id);
                hotel = hotelRes.success ? hotelRes.data : null;
              } catch {
                hotel = null;
              }
              hotelCache.set(room.hotel_id, hotel);
            }
          }

          const payment = payments.find((p) => p.booking_id === booking.id) ?? null;

          return { booking, room, hotel, payment };
        })
      );

      setItems(enriched);
      setStatus("success");
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleCancel(bookingId: number) {
    setCancellingId(bookingId);
    try {
      await cancelBooking(bookingId);
      setItems((prev) =>
        prev.map((item) =>
          item.booking.id === bookingId
            ? { ...item, booking: { ...item.booking, status: "cancelled" } }
            : item
        )
      );
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
    } finally {
      setCancellingId(null);
    }
  }

  const grouped = useMemo(() => {
    const groups: Record<Group, EnrichedBooking[]> = {
      upcoming: [],
      completed: [],
      cancelled: [],
    };
    for (const item of items) {
      groups[groupOf(item.booking)].push(item);
    }
    return groups;
  }, [items]);

  const TABS: { key: Group; label: string }[] = [
    { key: "upcoming", label: `Upcoming (${grouped.upcoming.length})` },
    { key: "completed", label: `Completed (${grouped.completed.length})` },
    { key: "cancelled", label: `Cancelled (${grouped.cancelled.length})` },
  ];

  const visible = grouped[activeGroup];

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-5xl px-4 py-10 sm:px-6 sm:py-14">
        <PageHeading
          eyebrow="Your trips"
          title="My bookings"
          subtitle="Everything you've booked with NileStay, in one place."
        />

        {status === "loading" && (
          <div className="mt-8 flex flex-col gap-6">
            {Array.from({ length: 3 }).map((_, i) => (
              <BookingCardSkeleton key={i} />
            ))}
          </div>
        )}

        {status === "error" && (
          <div className="mt-10">
            <ErrorState
              title="Couldn't load your bookings"
              message={errorMessage}
              onRetry={load}
            />
          </div>
        )}

        {status === "empty" && (
          <div className="mt-10">
            <EmptyState
              icon={CalendarX}
              title="No bookings yet"
              description="Once you book a stay, it'll show up here."
              action={
                <Link
                  to="/"
                  className="mt-2 rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white hover:bg-teal-dark"
                >
                  Browse hotels
                </Link>
              }
            />
          </div>
        )}

        {status === "success" && (
          <>
            <div className="mt-8 flex gap-2 overflow-x-auto">
              {TABS.map(({ key, label }) => (
                <button
                  key={key}
                  onClick={() => setActiveGroup(key)}
                  className={`shrink-0 rounded-pill px-4 py-2 text-sm font-semibold transition-colors cursor-pointer ${
                    activeGroup === key
                      ? "bg-teal text-white"
                      : "bg-bg text-muted hover:text-ink"
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>

            {visible.length === 0 ? (
              <div className="mt-6">
                <EmptyState
                  icon={CalendarX}
                  title={`No ${activeGroup} bookings`}
                />
              </div>
            ) : (
              <div className="mt-8 flex flex-col gap-6">
                {visible.map(({ booking, room, hotel, payment }, i) => {
                  const isCancellable =
                    booking.status === "pending" || booking.status === "confirmed";
                  // Prefer the actual room's own photo - it's what the
                  // guest booked - and only fall back to the hotel's photo
                  // when this specific room has none.
                  const imageUrl =
                    room?.image_url ?? (hotel ? getHotelImageUrl(hotel) : null);
                  return (
                    <article
                      key={booking.id}
                      className="card-enter flex flex-col gap-6 rounded-card border border-line bg-white p-7 shadow-[var(--shadow-card)] transition-shadow hover:shadow-[var(--shadow-hover)] sm:p-8 lg:grid lg:grid-cols-[160px_1fr_200px] lg:items-stretch lg:gap-8"
                      style={{ animationDelay: `${Math.min(i, 8) * 60}ms` }}
                    >
                      {/* Zone 1: image */}
                      <Link
                        to={`/bookings/${booking.id}`}
                        className="block h-40 w-full shrink-0 overflow-hidden rounded-[16px] bg-cream lg:h-auto"
                      >
                        {imageUrl ? (
                          <img src={imageUrl} alt="" className="h-full w-full object-cover" />
                        ) : (
                          <div className="flex h-full w-full items-center justify-center">
                            <BedDouble size={28} className="text-teal" />
                          </div>
                        )}
                      </Link>

                      {/* Zone 2: hotel/room info + dates/guests/status,
                          visually separated from the image and from the
                          price/actions column by a vertical divider on
                          desktop. */}
                      <div className="flex flex-col justify-center gap-5 lg:border-l lg:border-line lg:pl-8">
                        <Link to={`/bookings/${booking.id}`}>
                          <h2 className="font-display text-lg font-semibold text-ink">
                            {hotel ? hotel.name : `Hotel #${room?.hotel_id ?? "?"}`}
                          </h2>
                          <p className="mt-1 text-sm text-muted">
                            {room
                              ? `${room.type} · Room ${room.room_number}`
                              : `Room #${booking.room_id}`}
                          </p>
                        </Link>

                        <div className="flex flex-wrap gap-x-8 gap-y-4">
                          <div>
                            <p className="text-xs font-medium uppercase tracking-wide text-muted">
                              Check in
                            </p>
                            <p className="mt-1 flex items-center gap-1.5 text-sm font-semibold text-ink">
                              <CalendarDays size={14} className="text-muted" />
                              {new Date(booking.check_in).toLocaleDateString()}
                            </p>
                          </div>
                          <div>
                            <p className="text-xs font-medium uppercase tracking-wide text-muted">
                              Check out
                            </p>
                            <p className="mt-1 flex items-center gap-1.5 text-sm font-semibold text-ink">
                              <CalendarDays size={14} className="text-muted" />
                              {new Date(booking.check_out).toLocaleDateString()}
                            </p>
                          </div>
                          <div>
                            <p className="text-xs font-medium uppercase tracking-wide text-muted">
                              Guests
                            </p>
                            <p className="mt-1 text-sm font-semibold text-ink">
                              {booking.guests}
                            </p>
                          </div>
                        </div>

                        <div className="flex flex-wrap items-center gap-2">
                          <StatusBadge status={booking.status} />
                          {payment && <StatusBadge status={payment.status} />}
                        </div>
                      </div>

                      {/* Zone 3: price + actions, separated from the info
                          column by its own divider on desktop. */}
                      <div className="flex flex-row items-center justify-between gap-4 border-t border-line pt-6 lg:flex-col lg:items-end lg:justify-center lg:border-l lg:border-t-0 lg:pl-8 lg:pt-0 lg:text-right">
                        <div>
                          <p className="text-xs font-medium uppercase tracking-wide text-muted">
                            Total price
                          </p>
                          <p className="mt-1 font-display text-2xl font-bold text-gold">
                            ${booking.total_price.toFixed(2)}
                          </p>
                        </div>
                        <div className="flex flex-col gap-2.5 sm:flex-row lg:w-full lg:flex-col">
                          {!payment && booking.status !== "cancelled" && (
                            <Link
                              to={`/bookings/${booking.id}/payment`}
                              className="rounded-pill bg-teal px-5 py-2.5 text-center text-sm font-semibold text-white transition-colors hover:bg-teal-dark"
                            >
                              Pay now
                            </Link>
                          )}
                          {isCancellable && (
                            <button
                              onClick={() => handleCancel(booking.id)}
                              disabled={cancellingId === booking.id}
                              className="rounded-pill border border-line px-5 py-2.5 text-sm font-semibold text-muted transition-colors hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed cursor-pointer"
                            >
                              {cancellingId === booking.id ? "Cancelling..." : "Cancel"}
                            </button>
                          )}
                        </div>
                      </div>
                    </article>
                  );
                })}
              </div>
            )}
          </>
        )}
      </div>

      <Footer />
    </div>
  );
}
