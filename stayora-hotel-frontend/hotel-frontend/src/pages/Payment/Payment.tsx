import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  BedDouble,
  Loader2,
  TriangleAlert,
  Wallet,
  CreditCard,
  CheckCircle2,
  Clock,
  XCircle,
  Undo2,
} from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { StatusBadge } from "@/components/common/StatusBadge";
import { BookingProgress } from "@/components/common/BookingProgress";
import { AuthButton } from "@/components/auth/AuthButton";
import { getBookingByID } from "@/api/bookingApi";
import { getRoomByID } from "@/api/roomApi";
import { getHotelByID } from "@/api/hotelApi";
import { createPayment, getMyPayments } from "@/api/paymentApi";
import { extractErrorMessage } from "@/api/client";
import type { Booking } from "@/types/booking";
import type { Room } from "@/types/room";
import type { Hotel } from "@/types/hotel";
import type { Payment as PaymentType, PaymentMethod } from "@/types/payment";

type Status = "loading" | "success" | "error";

const CONFIRMATION_CONTENT: Record<
  string,
  { icon: typeof Clock; tone: string; title: string; body: string }
> = {
  pending: {
    icon: Clock,
    tone: "bg-amber-50 text-amber-600",
    title: "Payment pending",
    body: "We've recorded your payment and it's awaiting confirmation.",
  },
  paid: {
    icon: CheckCircle2,
    tone: "bg-emerald-50 text-emerald-600",
    title: "Booking confirmed",
    body: "Your payment was confirmed. We're looking forward to hosting you.",
  },
  failed: {
    icon: XCircle,
    tone: "bg-red-50 text-red-500",
    title: "Payment failed",
    body: "This payment didn't go through. Please contact support from My Bookings.",
  },
  refunded: {
    icon: Undo2,
    tone: "bg-bg text-muted",
    title: "Payment refunded",
    body: "This payment has been refunded.",
  },
};

export function Payment() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [booking, setBooking] = useState<Booking | null>(null);
  const [room, setRoom] = useState<Room | null>(null);
  const [hotel, setHotel] = useState<Hotel | null>(null);
  const [existingPayment, setExistingPayment] = useState<PaymentType | null>(null);

  const [method, setMethod] = useState<PaymentMethod>("card");
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [createdPayment, setCreatedPayment] = useState<PaymentType | null>(null);

  useEffect(() => {
    async function load() {
      if (!id) return;
      setStatus("loading");
      try {
        const bookingData = await getBookingByID(id);
        setBooking(bookingData);

        const roomRes = await getRoomByID(bookingData.room_id);
        if (roomRes.success) {
          setRoom(roomRes.data);
          const hotelRes = await getHotelByID(roomRes.data.hotel_id);
          if (hotelRes.success) setHotel(hotelRes.data);
        }

        const paymentsRes = await getMyPayments();
        if (paymentsRes.success) {
          const match = paymentsRes.data.find((p) => p.booking_id === bookingData.id);
          if (match) setExistingPayment(match);
        }

        setStatus("success");
      } catch (err) {
        setErrorMessage(extractErrorMessage(err));
        setStatus("error");
      }
    }
    load();
  }, [id]);

  async function handlePay() {
    if (!booking || isSubmitting) return;
    setFormError("");
    setIsSubmitting(true);
    try {
      const response = await createPayment(booking.id, method);
      if (response.success) {
        setCreatedPayment(response.data);
      } else {
        setFormError(response.message || "Couldn't process your payment.");
      }
    } catch (err) {
      setFormError(extractErrorMessage(err));
    } finally {
      setIsSubmitting(false);
    }
  }

  if (status === "loading") {
    return (
      <div className="min-h-screen bg-bg">
        <Navbar />
        <div className="flex min-h-[60vh] items-center justify-center">
          <Loader2 size={28} className="animate-spin text-teal" />
        </div>
      </div>
    );
  }

  if (status === "error" || !booking) {
    return (
      <div className="min-h-screen bg-bg">
        <Navbar />
        <div className="mx-auto flex max-w-2xl flex-col items-center gap-3 px-4 py-24 text-center">
          <span className="flex h-12 w-12 items-center justify-center rounded-full bg-red-50 text-red-500">
            <TriangleAlert size={22} />
          </span>
          <p className="font-display text-lg font-semibold text-ink">
            Couldn't load this booking
          </p>
          <p className="max-w-sm text-sm text-muted">{errorMessage}</p>
        </div>
      </div>
    );
  }

  const paymentToShow = createdPayment ?? existingPayment;

  // ============ Confirmation view ============
  if (paymentToShow) {
    const content = CONFIRMATION_CONTENT[paymentToShow.status] ?? CONFIRMATION_CONTENT.pending;
    const Icon = content.icon;

    return (
      <div className="min-h-screen bg-bg">
        <Navbar />
        <div className="page-fade-in mx-auto max-w-xl px-4 py-14 sm:px-6">
          <div className="rounded-card border border-line bg-white p-8 text-center shadow-[0_16px_40px_rgba(16,24,40,0.08)]">
            <span className={`mx-auto flex h-14 w-14 items-center justify-center rounded-full ${content.tone}`}>
              <Icon size={26} />
            </span>
            <h1 className="mt-4 font-display text-2xl font-bold text-ink">
              {content.title}
            </h1>
            <p className="mt-2 text-sm leading-relaxed text-muted">{content.body}</p>

            <div className="mt-6 flex items-center justify-center gap-2">
              <StatusBadge status={booking.status} />
              <StatusBadge status={paymentToShow.status} />
            </div>

            <div className="mt-6 flex flex-col gap-2 rounded-[13px] bg-bg px-5 py-4 text-left text-sm">
              <div className="flex justify-between">
                <span className="text-muted">Booking ID</span>
                <span className="font-semibold text-ink">#{booking.id}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted">Hotel</span>
                <span className="font-semibold text-ink">
                  {hotel ? hotel.name : `#${room?.hotel_id ?? "?"}`}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted">Room</span>
                <span className="font-semibold text-ink">
                  {room ? `${room.type} · ${room.room_number}` : `#${booking.room_id}`}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted">Dates</span>
                <span className="font-semibold text-ink">
                  {new Date(booking.check_in).toLocaleDateString()} –{" "}
                  {new Date(booking.check_out).toLocaleDateString()}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted">Guests</span>
                <span className="font-semibold text-ink">{booking.guests}</span>
              </div>
              <div className="flex justify-between border-t border-line pt-2">
                <span className="text-muted">Total paid</span>
                <span className="font-display text-base font-bold text-ink">
                  ${booking.total_price.toFixed(2)}
                </span>
              </div>
            </div>

            <div className="mt-7 flex flex-col gap-2.5 sm:flex-row sm:justify-center">
              <button
                onClick={() => navigate(`/bookings/${booking.id}`)}
                className="rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark cursor-pointer"
              >
                View booking
              </button>
              <button
                onClick={() => navigate("/my-bookings")}
                className="rounded-pill border border-line px-5 py-2.5 text-sm font-semibold text-ink transition-colors hover:border-teal/40 cursor-pointer"
              >
                My bookings
              </button>
              <button
                onClick={() => navigate("/")}
                className="rounded-pill px-5 py-2.5 text-sm font-semibold text-muted transition-colors hover:text-ink cursor-pointer"
              >
                Back home
              </button>
            </div>
          </div>
        </div>
        <Footer />
      </div>
    );
  }

  // ============ Payment method selection ============
  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-5xl px-4 py-10 sm:px-6">
        <p className="text-sm font-semibold tracking-wide text-teal">
          {hotel ? hotel.name : `Booking #${booking.id}`}
        </p>
        <h1 className="mt-1.5 font-display text-2xl font-bold text-ink sm:text-3xl">
          Payment
        </h1>

        <div className="mt-6 rounded-card border border-line bg-white p-4 sm:p-5">
          <BookingProgress current={3} />
        </div>

        <div className="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-5">
          <div className="flex flex-col gap-6 rounded-card border border-line bg-white p-6 lg:col-span-3">
            <div>
              <h2 className="font-display text-lg font-semibold text-ink">
                Choose a payment method
              </h2>
              <p className="mt-1 text-sm text-muted">
                Your booking is reserved as pending until payment is recorded.
              </p>
            </div>

            <div className="grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() => setMethod("card")}
                className={`flex flex-col items-center justify-center gap-2 rounded-[13px] border px-4 py-6 text-sm font-semibold transition-colors cursor-pointer ${
                  method === "card"
                    ? "border-teal bg-teal/5 text-teal"
                    : "border-line text-muted hover:text-ink"
                }`}
              >
                <CreditCard size={22} />
                Card
              </button>
              <button
                type="button"
                onClick={() => setMethod("cash")}
                className={`flex flex-col items-center justify-center gap-2 rounded-[13px] border px-4 py-6 text-sm font-semibold transition-colors cursor-pointer ${
                  method === "cash"
                    ? "border-teal bg-teal/5 text-teal"
                    : "border-line text-muted hover:text-ink"
                }`}
              >
                <Wallet size={22} />
                Cash
              </button>
            </div>

            {formError && (
              <p role="alert" className="rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600">
                {formError}
              </p>
            )}

            <AuthButton
              isSubmitting={isSubmitting}
              disabled={isSubmitting}
              onClick={handlePay}
            >
              Pay ${booking.total_price.toFixed(2)}
            </AuthButton>
          </div>

          <div className="lg:col-span-2 lg:sticky lg:top-24 lg:self-start">
            <div className="overflow-hidden rounded-card border border-line bg-white shadow-[0_16px_40px_rgba(16,24,40,0.08)]">
              <div className="h-24 overflow-hidden bg-cream">
                {room?.image_url ? (
                  <img
                    src={room.image_url}
                    alt={`${room.type} · Room ${room.room_number}`}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center">
                    <BedDouble size={26} className="text-teal" strokeWidth={1.5} />
                  </div>
                )}
              </div>
              <div className="flex flex-col gap-4 p-6">
                <div>
                  <h3 className="font-display text-base font-semibold text-ink">
                    {hotel ? hotel.name : `Booking #${booking.id}`}
                  </h3>
                  <p className="mt-1 text-sm text-muted">
                    {room ? `${room.type} · Room ${room.room_number}` : `Room #${booking.room_id}`}
                  </p>
                </div>
                <div className="flex flex-col gap-2 border-t border-line pt-4 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted">Dates</span>
                    <span className="font-medium text-ink">
                      {new Date(booking.check_in).toLocaleDateString()} –{" "}
                      {new Date(booking.check_out).toLocaleDateString()}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted">Guests</span>
                    <span className="font-medium text-ink">{booking.guests}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted">Status</span>
                    <StatusBadge status={booking.status} />
                  </div>
                </div>
                <div className="flex items-center justify-between border-t border-line pt-4">
                  <span className="font-display text-base font-bold text-ink">Total</span>
                  <span className="font-display text-xl font-bold text-ink">
                    ${booking.total_price.toFixed(2)}
                  </span>
                </div>
              </div>
            </div>

            <p className="mt-3 text-center text-xs text-muted">
              <Link to="/my-bookings" className="font-semibold text-teal">
                View all my bookings
              </Link>
            </p>
          </div>
        </div>
      </div>

      <Footer />
    </div>
  );
}
