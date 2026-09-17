import { useEffect, useState, type FormEvent } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import {
  BedDouble,
  CalendarDays,
  Users,
  Loader2,
  TriangleAlert,
  CheckCircle2,
  XCircle,
  MapPin,
} from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { AuthButton } from "@/components/auth/AuthButton";
import { BookingProgress } from "@/components/common/BookingProgress";
import { getRoomByID, checkRoomAvailability } from "@/api/roomApi";
import { getHotelByID } from "@/api/hotelApi";
import { createBooking } from "@/api/bookingApi";
import { extractErrorMessage } from "@/api/client";
import type { Room } from "@/types/room";
import type { Hotel } from "@/types/hotel";

type Status = "loading" | "success" | "error";
type AvailabilityStatus = "idle" | "checking" | "available" | "unavailable";

function nightsBetween(checkIn: string, checkOut: string): number {
  if (!checkIn || !checkOut) return 0;
  const inDate = new Date(checkIn);
  const outDate = new Date(checkOut);
  const diff = outDate.getTime() - inDate.getTime();
  return Math.max(0, Math.round(diff / (1000 * 60 * 60 * 24)));
}

export function Booking() {
  const { roomId } = useParams<{ roomId: string }>();
  const navigate = useNavigate();

  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [room, setRoom] = useState<Room | null>(null);
  const [hotel, setHotel] = useState<Hotel | null>(null);

  const [checkIn, setCheckIn] = useState("");
  const [checkOut, setCheckOut] = useState("");
  const [guests, setGuests] = useState(1);
  const [formError, setFormError] = useState("");
  const [availability, setAvailability] = useState<AvailabilityStatus>("idle");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    async function load() {
      if (!roomId) return;
      setStatus("loading");
      try {
        const roomRes = await getRoomByID(roomId);
        if (!roomRes.success) {
          setErrorMessage(roomRes.message || "Couldn't load this room.");
          setStatus("error");
          return;
        }
        setRoom(roomRes.data);

        const hotelRes = await getHotelByID(roomRes.data.hotel_id);
        if (hotelRes.success) setHotel(hotelRes.data);

        setStatus("success");
      } catch (err) {
        setErrorMessage(extractErrorMessage(err));
        setStatus("error");
      }
    }
    load();
  }, [roomId]);

  const nights = nightsBetween(checkIn, checkOut);
  const estimatedTotal = room ? nights * room.price_per_night : 0;

  async function handleCheckAvailability() {
    if (!room || !checkIn || !checkOut) return;
    setAvailability("checking");
    try {
      const result = await checkRoomAvailability(room.id, checkIn, checkOut);
      setAvailability(result.available ? "available" : "unavailable");
    } catch {
      setAvailability("idle");
    }
  }

  function validate(): boolean {
    if (!checkIn || !checkOut) {
      setFormError("Please select both check-in and check-out dates.");
      return false;
    }
    if (nights <= 0) {
      setFormError("Check-out date must be after check-in date.");
      return false;
    }
    if (room && guests > room.capacity) {
      setFormError(`This room fits up to ${room.capacity} guests.`);
      return false;
    }
    if (guests < 1) {
      setFormError("At least 1 guest is required.");
      return false;
    }
    setFormError("");
    return true;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!room || !validate() || isSubmitting) return;

    setIsSubmitting(true);
    setFormError("");
    try {
      // The backend re-validates availability and computes the total price
      // itself - what we show below is a preview, never the source of truth.
      const booking = await createBooking({
        room_id: room.id,
        check_in: checkIn,
        check_out: checkOut,
        guests,
      });
      navigate(`/bookings/${booking.id}/payment`);
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

  if (status === "error" || !room) {
    return (
      <div className="min-h-screen bg-bg">
        <Navbar />
        <div className="mx-auto flex max-w-2xl flex-col items-center gap-3 px-4 py-24 text-center">
          <span className="flex h-12 w-12 items-center justify-center rounded-full bg-red-50 text-red-500">
            <TriangleAlert size={22} />
          </span>
          <p className="font-display text-lg font-semibold text-ink">
            Couldn't load this room
          </p>
          <p className="max-w-sm text-sm text-muted">{errorMessage}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-5xl px-4 py-10 sm:px-6">
        <p className="text-sm font-semibold tracking-wide text-teal">
          {hotel ? hotel.name : `Hotel #${room.hotel_id}`}
        </p>
        <h1 className="mt-1.5 font-display text-2xl font-bold text-ink sm:text-3xl">
          Complete your booking
        </h1>

        <div className="mt-6 rounded-card border border-line bg-white p-4 sm:p-5">
          <BookingProgress current={nights > 0 ? 2 : 1} />
        </div>

        <div className="mt-8 grid grid-cols-1 gap-8 lg:grid-cols-5">
          {/* Form */}
          <form
            onSubmit={handleSubmit}
            className="flex flex-col gap-6 rounded-card border border-line bg-white p-6 lg:col-span-3"
          >
            <div>
              <h2 className="font-display text-lg font-semibold text-ink">
                Dates &amp; guests
              </h2>
              <p className="mt-1 text-sm text-muted">
                Choose when you'd like to stay and how many guests are joining.
              </p>
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="flex flex-col gap-1.5">
                <label className="flex items-center gap-1.5 text-sm font-medium text-ink">
                  <CalendarDays size={15} />
                  Check in
                </label>
                <input
                  type="date"
                  value={checkIn}
                  min={new Date().toISOString().slice(0, 10)}
                  onChange={(e) => {
                    setCheckIn(e.target.value);
                    setAvailability("idle");
                  }}
                  className="w-full rounded-[13px] border border-line bg-white px-4 py-3 text-sm text-ink outline-none focus:border-teal focus:ring-4 focus:ring-teal/10 [color-scheme:light]"
                />
              </div>
              <div className="flex flex-col gap-1.5">
                <label className="flex items-center gap-1.5 text-sm font-medium text-ink">
                  <CalendarDays size={15} />
                  Check out
                </label>
                <input
                  type="date"
                  value={checkOut}
                  min={checkIn || new Date().toISOString().slice(0, 10)}
                  onChange={(e) => {
                    setCheckOut(e.target.value);
                    setAvailability("idle");
                  }}
                  className="w-full rounded-[13px] border border-line bg-white px-4 py-3 text-sm text-ink outline-none focus:border-teal focus:ring-4 focus:ring-teal/10 [color-scheme:light]"
                />
              </div>
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="flex items-center gap-1.5 text-sm font-medium text-ink">
                <Users size={15} />
                Guests
              </label>
              <input
                type="number"
                min={1}
                max={room.capacity}
                value={guests}
                onChange={(e) => setGuests(Number(e.target.value))}
                className="w-full rounded-[13px] border border-line bg-white px-4 py-3 text-sm text-ink outline-none focus:border-teal focus:ring-4 focus:ring-teal/10 sm:w-40"
              />
              <p className="text-xs text-muted">Up to {room.capacity} guests for this room.</p>
            </div>

            {checkIn && checkOut && nights > 0 && (
              <button
                type="button"
                onClick={handleCheckAvailability}
                disabled={availability === "checking"}
                className="self-start rounded-pill border border-line px-4 py-2 text-xs font-semibold text-ink transition-colors hover:border-teal/40 disabled:cursor-not-allowed cursor-pointer"
              >
                {availability === "checking"
                  ? "Checking availability..."
                  : "Check availability"}
              </button>
            )}

            {availability === "available" && (
              <p className="flex items-center gap-1.5 text-sm font-medium text-emerald-600">
                <CheckCircle2 size={15} />
                This room is available for your dates.
              </p>
            )}
            {availability === "unavailable" && (
              <p className="flex items-center gap-1.5 text-sm font-medium text-red-500">
                <XCircle size={15} />
                This room is already booked for those dates.
              </p>
            )}

            {formError && (
              <p
                role="alert"
                className="rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
              >
                {formError}
              </p>
            )}

            <AuthButton
              isSubmitting={isSubmitting}
              disabled={isSubmitting || availability === "unavailable"}
              className="mt-2"
            >
              Confirm booking
            </AuthButton>

            <p className="text-center text-xs text-muted">
              Availability and pricing are re-verified by the server when you
              confirm.
            </p>
          </form>

          {/* Sticky summary */}
          <div className="lg:col-span-2 lg:sticky lg:top-24 lg:self-start">
            <div className="overflow-hidden rounded-card border border-line bg-white shadow-[0_16px_40px_rgba(16,24,40,0.08)]">
              <div className="h-28 overflow-hidden bg-cream">
                {room.image_url ? (
                  <img
                    src={room.image_url}
                    alt={`${room.type} · Room ${room.room_number}`}
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="flex h-full w-full items-center justify-center">
                    <BedDouble size={30} className="text-teal" strokeWidth={1.5} />
                  </div>
                )}
              </div>
              <div className="flex flex-col gap-4 p-6">
                <div>
                  <h3 className="font-display text-base font-semibold text-ink">
                    {hotel ? hotel.name : `Hotel #${room.hotel_id}`}
                  </h3>
                  {hotel && (
                    <p className="mt-1 flex items-center gap-1.5 text-xs text-muted">
                      <MapPin size={12} />
                      {hotel.city}
                      {hotel.country ? `, ${hotel.country}` : ""}
                    </p>
                  )}
                  <p className="mt-2 text-sm text-muted">
                    {room.type} &middot; Room {room.room_number}
                  </p>
                </div>

                <div className="flex flex-col gap-2 border-t border-line pt-4 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted">Check in</span>
                    <span className="font-medium text-ink">
                      {checkIn ? new Date(checkIn).toLocaleDateString() : "—"}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted">Check out</span>
                    <span className="font-medium text-ink">
                      {checkOut ? new Date(checkOut).toLocaleDateString() : "—"}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted">Guests</span>
                    <span className="font-medium text-ink">{guests}</span>
                  </div>
                </div>

                <div className="flex flex-col gap-2 border-t border-line pt-4 text-sm">
                  <div className="flex justify-between text-muted">
                    <span>
                      ${room.price_per_night.toFixed(2)} &times; {nights || 0}{" "}
                      {nights === 1 ? "night" : "nights"}
                    </span>
                    <span>${estimatedTotal.toFixed(2)}</span>
                  </div>
                  <div className="flex items-center justify-between border-t border-line pt-3">
                    <span className="font-display text-base font-bold text-ink">
                      Total
                    </span>
                    <span className="font-display text-xl font-bold text-ink">
                      ${estimatedTotal.toFixed(2)}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <p className="mt-3 text-center text-xs text-muted">
              <Link to={`/hotels/${room.hotel_id}`} className="font-semibold text-teal">
                Back to hotel
              </Link>
            </p>
          </div>
        </div>
      </div>

      <Footer />
    </div>
  );
}
