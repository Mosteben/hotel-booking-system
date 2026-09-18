import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import {
  BedDouble,
  CalendarDays,
  Users,
  Loader2,
  MapPin,
  ArrowLeft,
} from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { StatusBadge } from "@/components/common/StatusBadge";
import { ErrorState } from "@/components/common/ErrorState";
import { getBookingByID, cancelBooking } from "@/api/bookingApi";
import { getRoomByID } from "@/api/roomApi";
import { getHotelByID } from "@/api/hotelApi";
import { getMyPayments } from "@/api/paymentApi";
import { extractErrorMessage } from "@/api/client";
import { getHotelImageUrl } from "@/utils/hotelImages";
import type { Booking } from "@/types/booking";
import type { Room } from "@/types/room";
import type { Hotel } from "@/types/hotel";
import type { Payment } from "@/types/payment";

type Status = "loading" | "success" | "error";

export function BookingDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();

  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [booking, setBooking] = useState<Booking | null>(null);
  const [room, setRoom] = useState<Room | null>(null);
  const [hotel, setHotel] = useState<Hotel | null>(null);
  const [payment, setPayment] = useState<Payment | null>(null);
  const [isCancelling, setIsCancelling] = useState(false);
  const [cancelError, setCancelError] = useState("");

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

      try {
        const paymentsRes = await getMyPayments();
        if (paymentsRes.success) {
          setPayment(
            paymentsRes.data.find((p) => p.booking_id === bookingData.id) ?? null
          );
        }
      } catch {
        // Payment status is supplementary here.
      }

      setStatus("success");
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }
  }

  useEffect(() => {
    load();
  }, [id]);

  async function handleCancel() {
    if (!booking) return;
    setCancelError("");
    setIsCancelling(true);
    try {
      await cancelBooking(booking.id);
      setBooking({ ...booking, status: "cancelled" });
    } catch (err) {
      setCancelError(extractErrorMessage(err));
    } finally {
      setIsCancelling(false);
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
        <div className="mx-auto max-w-2xl px-4 py-24 sm:px-6">
          <ErrorState title="Couldn't load this booking" message={errorMessage} />
        </div>
      </div>
    );
  }

  // The room's own photo is what the guest actually booked - only fall
  // back to the hotel's photo when this specific room has none.
  const imageUrl = room?.image_url ?? (hotel ? getHotelImageUrl(hotel) : null);
  // The backend only supports cancelling a pending/confirmed booking via
  // DELETE /bookings/:id - no modify endpoint exists, so no "Modify" action
  // is shown.
  const isCancellable = booking.status === "pending" || booking.status === "confirmed";

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-3xl px-4 py-10 sm:px-6">
        <button
          onClick={() => navigate("/my-bookings")}
          className="flex items-center gap-1.5 text-sm font-medium text-muted transition-colors hover:text-ink cursor-pointer"
        >
          <ArrowLeft size={15} />
          Back to my bookings
        </button>

        <div className="mt-4 flex flex-wrap items-center justify-between gap-3">
          <div>
            <p className="text-sm font-semibold tracking-wide text-teal">
              Booking #{booking.id}
            </p>
            <h1 className="mt-1 font-display text-2xl font-bold text-ink sm:text-3xl">
              {hotel ? hotel.name : `Hotel #${room?.hotel_id ?? "?"}`}
            </h1>
          </div>
          <div className="flex items-center gap-2">
            <StatusBadge status={booking.status} />
            {payment && <StatusBadge status={payment.status} />}
          </div>
        </div>

        <div className="mt-6 overflow-hidden rounded-card border border-line bg-white shadow-[var(--shadow-card)]">
          <div className="flex h-40 items-center justify-center bg-cream">
            {imageUrl ? (
              <img src={imageUrl} alt={hotel?.name} className="h-full w-full object-cover" />
            ) : (
              <BedDouble size={32} className="text-teal" strokeWidth={1.5} />
            )}
          </div>

          <div className="flex flex-col gap-6 p-6">
            {hotel && (
              <p className="flex items-center gap-1.5 text-sm text-muted">
                <MapPin size={14} />
                {hotel.city}
                {hotel.country ? `, ${hotel.country}` : ""}
              </p>
            )}

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="flex items-start gap-3">
                <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
                  <BedDouble size={16} />
                </span>
                <div>
                  <p className="text-xs font-medium text-muted">Room</p>
                  <p className="text-sm font-semibold text-ink">
                    {room ? `${room.type} · ${room.room_number}` : `#${booking.room_id}`}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
                  <Users size={16} />
                </span>
                <div>
                  <p className="text-xs font-medium text-muted">Guests</p>
                  <p className="text-sm font-semibold text-ink">{booking.guests}</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
                  <CalendarDays size={16} />
                </span>
                <div>
                  <p className="text-xs font-medium text-muted">Check in</p>
                  <p className="text-sm font-semibold text-ink">
                    {new Date(booking.check_in).toLocaleDateString()}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
                  <CalendarDays size={16} />
                </span>
                <div>
                  <p className="text-xs font-medium text-muted">Check out</p>
                  <p className="text-sm font-semibold text-ink">
                    {new Date(booking.check_out).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between rounded-[13px] bg-bg px-4 py-3">
              <span className="text-sm font-medium text-muted">Total price</span>
              <span className="font-display text-xl font-bold text-gold">
                ${booking.total_price.toFixed(2)}
              </span>
            </div>

            {cancelError && (
              <p role="alert" className="rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600">
                {cancelError}
              </p>
            )}

            <div className="flex flex-wrap gap-3">
              {!payment && booking.status !== "cancelled" && (
                <Link
                  to={`/bookings/${booking.id}/payment`}
                  className="rounded-pill bg-teal px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark"
                >
                  Pay now
                </Link>
              )}
              {isCancellable && (
                <button
                  onClick={handleCancel}
                  disabled={isCancelling}
                  className="rounded-pill border border-line px-5 py-2.5 text-sm font-semibold text-muted transition-colors hover:border-red-200 hover:text-red-500 disabled:cursor-not-allowed cursor-pointer"
                >
                  {isCancelling ? "Cancelling..." : "Cancel booking"}
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      <Footer />
    </div>
  );
}
