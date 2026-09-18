import { useEffect, useMemo, useState, type FormEvent } from "react";
import { useNavigate, useParams } from "react-router-dom";
import {
  Building2,
  MapPin,
  Phone,
  Mail,
  Star,
  MessageSquareText,
  BedDouble,
  ArrowDown,
} from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { FavoriteButton } from "@/components/common/FavoriteButton";
import { EmptyState } from "@/components/common/EmptyState";
import { ErrorState } from "@/components/common/ErrorState";
import { RoomCard } from "@/components/RoomCard/RoomCard";
import { AuthButton } from "@/components/auth/AuthButton";
import { getHotelDetails } from "@/api/hotelApi";
import { getHotelAverageRating, getHotelReviews, createReview } from "@/api/reviewApi";
import { getMyFavorites } from "@/api/favoriteApi";
import { extractErrorMessage } from "@/api/client";
import { getHotelImageUrl } from "@/utils/hotelImages";
import { Skeleton } from "@/components/common/Skeleton";
import { useAuth } from "@/context/AuthContext";
import type { HotelDetails as HotelDetailsType } from "@/types/hotel";
import type { Room } from "@/types/room";
import type { Review } from "@/types/review";

type Status = "loading" | "success" | "error";

export function HotelDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const [status, setStatus] = useState<Status>("loading");
  const [errorMessage, setErrorMessage] = useState("");
  const [hotel, setHotel] = useState<HotelDetailsType | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [averageRating, setAverageRating] = useState<number | null>(null);
  const [isFavorite, setIsFavorite] = useState(false);

  const [rating, setRating] = useState(5);
  const [comment, setComment] = useState("");
  const [reviewError, setReviewError] = useState("");
  const [isSubmittingReview, setIsSubmittingReview] = useState(false);

  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  async function load() {
    if (!id) return;
    setStatus("loading");
    try {
      const [hotelRes, reviewsRes, ratingRes] = await Promise.all([
        getHotelDetails(id),
        getHotelReviews(id),
        getHotelAverageRating(id),
      ]);

      if (!hotelRes.success) {
        setErrorMessage(hotelRes.message || "Couldn't load this hotel.");
        setStatus("error");
        return;
      }

      setHotel(hotelRes.data);
      setReviews(reviewsRes.success ? reviewsRes.data : []);
      setAverageRating(ratingRes.success ? ratingRes.data.average_rating : null);

      if (isAuthenticated) {
        try {
          const favRes = await getMyFavorites();
          if (favRes.success) {
            setIsFavorite(favRes.data.some((f) => f.hotel_id === Number(id)));
          }
        } catch {
          // Favorite state is a nice-to-have here - don't fail the page for it.
        }
      }

      setStatus("success");
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [id]);

  function handleBook(room: Room) {
    navigate(`/booking/${room.id}`);
  }

  function scrollToRooms() {
    document.getElementById("rooms")?.scrollIntoView({ behavior: "smooth" });
  }

  async function handleReviewSubmit(e: FormEvent) {
    e.preventDefault();
    if (!id) return;
    setReviewError("");
    setIsSubmittingReview(true);
    try {
      const response = await createReview(id, { rating, comment });
      if (response.success) {
        setComment("");
        setRating(5);
        const reviewsRes = await getHotelReviews(id);
        if (reviewsRes.success) setReviews(reviewsRes.data);
        const ratingRes = await getHotelAverageRating(id);
        if (ratingRes.success) setAverageRating(ratingRes.data.average_rating);
      } else {
        setReviewError(response.message || "Couldn't submit your review.");
      }
    } catch (err) {
      setReviewError(extractErrorMessage(err));
    } finally {
      setIsSubmittingReview(false);
    }
  }

  const images = useMemo(() => hotel?.images ?? [], [hotel]);

  useEffect(() => {
    if (images.length === 0) return;
    const mainIndex = images.findIndex((img) => img.is_main);
    setSelectedImageIndex(mainIndex >= 0 ? mainIndex : 0);
  }, [images]);

  const rooms = useMemo(() => hotel?.rooms ?? [], [hotel]);
  const startingPrice = useMemo(() => {
    const priced = rooms.filter((r) => r.price_per_night > 0);
    if (priced.length === 0) return null;
    return Math.min(...priced.map((r) => r.price_per_night));
  }, [rooms]);
  const availableRoomCount = rooms.filter((r) => r.status === "available").length;

  if (status === "loading") {
    return (
      <div className="min-h-screen bg-bg">
        <Navbar />
        <div className="mx-auto max-w-6xl px-4 pt-8 sm:px-6">
          <Skeleton className="h-72 w-full rounded-panel sm:h-96" />
          <div className="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-3">
            <div className="flex flex-col gap-3 lg:col-span-2">
              <Skeleton className="h-7 w-72" />
              <Skeleton className="h-4 w-48" />
              {Array.from({ length: 2 }).map((_, i) => (
                <Skeleton key={i} className="mt-4 h-28 w-full rounded-card" />
              ))}
            </div>
            <Skeleton className="h-56 w-full rounded-card" />
          </div>
        </div>
      </div>
    );
  }

  if (status === "error" || !hotel) {
    return (
      <div className="min-h-screen bg-bg">
        <Navbar />
        <div className="mx-auto max-w-6xl px-4 py-24 sm:px-6">
          <ErrorState
            title="Couldn't load this hotel"
            message={errorMessage}
            onRetry={load}
          />
        </div>
      </div>
    );
  }

  const imageUrl = getHotelImageUrl(hotel);
  const bannerImageUrl = images[selectedImageIndex]?.url ?? imageUrl;

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-6xl px-4 pt-8 sm:px-6">
        {/* Banner */}
        <div className="relative h-72 w-full overflow-hidden rounded-panel bg-cream sm:h-96">
          {bannerImageUrl ? (
            <img
              src={bannerImageUrl}
              alt={hotel.name}
              className="h-full w-full object-cover"
            />
          ) : (
            <div className="flex h-full w-full flex-col items-center justify-center gap-2">
              <span className="flex h-14 w-14 items-center justify-center rounded-full bg-white text-teal shadow-sm">
                <Building2 size={26} />
              </span>
              <p className="text-sm font-medium text-muted">No photo yet</p>
            </div>
          )}
          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-gradient-to-t from-black/40 to-transparent" />
          <div className="absolute bottom-4 left-4 right-4 flex items-end justify-between gap-3 sm:bottom-6 sm:left-6 sm:right-6">
            <div>
              <h1 className="font-display text-2xl font-bold text-white drop-shadow sm:text-3xl">
                {hotel.name}
              </h1>
              <p className="mt-1 flex items-center gap-1.5 text-sm text-white/90">
                <MapPin size={14} />
                {hotel.city}
                {hotel.country ? `, ${hotel.country}` : ""}
              </p>
            </div>
            <FavoriteButton
              hotelId={hotel.id}
              isFavorite={isFavorite}
              onChange={setIsFavorite}
            />
          </div>
        </div>

        {/* Gallery thumbnails - only worth showing once there's a choice to make */}
        {images.length > 1 && (
          <div className="mt-3 flex gap-2 overflow-x-auto pb-1">
            {images.map((image, index) => (
              <button
                key={image.id}
                type="button"
                onClick={() => setSelectedImageIndex(index)}
                aria-label={`Show photo ${index + 1}`}
                className={`h-16 w-20 shrink-0 overflow-hidden rounded-[10px] transition-all ${
                  index === selectedImageIndex
                    ? "ring-2 ring-teal ring-offset-2 ring-offset-bg"
                    : "opacity-70 hover:opacity-100"
                }`}
              >
                <img
                  src={image.url}
                  alt=""
                  className="h-full w-full object-cover"
                />
              </button>
            ))}
          </div>
        )}

        {/* Main content + sticky sidebar */}
        <div className="mt-8 grid grid-cols-1 gap-8 pb-16 lg:grid-cols-3">
          <div className="flex flex-col gap-14 lg:col-span-2">
            {/* Overview */}
            <section>
              <div className="flex flex-wrap items-center gap-2.5">
                {hotel.stars > 0 && (
                  <span className="inline-flex items-center gap-1 rounded-pill bg-bg px-2.5 py-1 text-xs font-semibold text-ink">
                    <Star size={12} className="fill-gold text-gold" />
                    {hotel.stars}-star hotel
                  </span>
                )}
                {averageRating !== null && averageRating > 0 && (
                  <span className="inline-flex items-center gap-1 rounded-pill bg-gold-light px-2.5 py-1 text-xs font-semibold text-gold">
                    <Star size={12} className="fill-gold text-gold" />
                    {averageRating.toFixed(1)} ({reviews.length}{" "}
                    {reviews.length === 1 ? "review" : "reviews"})
                  </span>
                )}
                <span className="inline-flex items-center gap-1 rounded-pill bg-bg px-2.5 py-1 text-xs font-semibold text-muted">
                  {availableRoomCount > 0
                    ? `${availableRoomCount} room${availableRoomCount === 1 ? "" : "s"} available`
                    : "No rooms currently available"}
                </span>
              </div>

              {hotel.description && (
                <p className="mt-4 max-w-2xl text-sm leading-relaxed text-muted">
                  {hotel.description}
                </p>
              )}

              {(hotel.address || hotel.phone || hotel.email) && (
                <div className="mt-6 rounded-card border border-line bg-white p-5">
                  <p className="text-xs font-semibold uppercase tracking-wider text-muted">
                    Hotel information
                  </p>
                  <div className="mt-3 flex flex-wrap gap-x-6 gap-y-2.5 text-sm text-ink">
                    {hotel.address && (
                      <span className="flex items-center gap-1.5">
                        <MapPin size={14} className="shrink-0 text-teal" />
                        {hotel.address}
                      </span>
                    )}
                    {hotel.phone && (
                      <span className="flex items-center gap-1.5">
                        <Phone size={14} className="shrink-0 text-teal" />
                        {hotel.phone}
                      </span>
                    )}
                    {hotel.email && (
                      <span className="flex items-center gap-1.5">
                        <Mail size={14} className="shrink-0 text-teal" />
                        {hotel.email}
                      </span>
                    )}
                  </div>
                </div>
              )}
            </section>

            {/* Rooms */}
            <section id="rooms" className="scroll-mt-24">
              <h2 className="font-display text-xl font-bold text-ink">
                Available rooms
              </h2>
              {rooms.length === 0 ? (
                <div className="mt-4">
                  <EmptyState
                    icon={BedDouble}
                    title="No rooms yet"
                    description="Rooms will appear here as soon as they're added to this hotel."
                  />
                </div>
              ) : (
                <div className="mt-6 flex flex-col gap-6">
                  {rooms.map((room, i) => (
                    <div
                      key={room.id}
                      className="card-enter"
                      style={{ animationDelay: `${Math.min(i, 8) * 60}ms` }}
                    >
                      <RoomCard room={room} onBook={handleBook} />
                    </div>
                  ))}
                </div>
              )}
            </section>

            {/* Reviews */}
            <section>
              <h2 className="font-display text-xl font-bold text-ink">Reviews</h2>

              {isAuthenticated ? (
                <form
                  onSubmit={handleReviewSubmit}
                  className="mt-4 rounded-card border border-line bg-white p-5"
                >
                  <div className="flex items-center gap-1">
                    {[1, 2, 3, 4, 5].map((value) => (
                      <button
                        key={value}
                        type="button"
                        onClick={() => setRating(value)}
                        aria-label={`Rate ${value} stars`}
                        className="cursor-pointer p-0.5"
                      >
                        <Star
                          size={20}
                          className={
                            value <= rating ? "fill-gold text-gold" : "text-line"
                          }
                        />
                      </button>
                    ))}
                  </div>
                  <textarea
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                    placeholder="Share your experience..."
                    rows={3}
                    className="mt-3 w-full resize-none rounded-[13px] border border-line bg-white px-4 py-3 text-sm text-ink outline-none placeholder:text-muted focus:border-teal focus:ring-4 focus:ring-teal/10"
                  />
                  {reviewError && (
                    <p role="alert" className="mt-2 text-xs font-medium text-red-500">
                      {reviewError}
                    </p>
                  )}
                  <AuthButton isSubmitting={isSubmittingReview} className="mt-3 w-auto px-6">
                    Submit review
                  </AuthButton>
                </form>
              ) : (
                <p className="mt-3 text-sm text-muted">
                  <button
                    onClick={() => navigate("/login")}
                    className="font-semibold text-teal hover:text-teal-dark cursor-pointer"
                  >
                    Log in
                  </button>{" "}
                  to leave a review.
                </p>
              )}

              {reviews.length === 0 ? (
                <div className="mt-6">
                  <EmptyState icon={MessageSquareText} title="No reviews yet" />
                </div>
              ) : (
                <div className="mt-8 flex flex-col gap-6">
                  {reviews.map((review) => (
                    <div
                      key={review.id}
                      className="rounded-card border border-line bg-white p-5"
                    >
                      <div className="flex items-center gap-1">
                        {Array.from({ length: 5 }).map((_, i) => (
                          <Star
                            key={i}
                            size={14}
                            className={
                              i < review.rating ? "fill-gold text-gold" : "text-line"
                            }
                          />
                        ))}
                      </div>
                      {review.comment && (
                        <p className="mt-2 text-sm leading-relaxed text-ink">
                          {review.comment}
                        </p>
                      )}
                      <p className="mt-2 text-xs text-muted">
                        {new Date(review.created_at).toLocaleDateString()}
                      </p>
                    </div>
                  ))}
                </div>
              )}
            </section>
          </div>

          {/* Sticky booking sidebar */}
          <aside className="lg:sticky lg:top-24 lg:self-start">
            <div className="rounded-card border border-line bg-white p-6 shadow-[var(--shadow-hover)]">
              <p className="text-xs font-semibold uppercase tracking-wider text-muted">
                Starting from
              </p>
              {startingPrice !== null ? (
                <p className="mt-1 font-display text-3xl font-bold text-gold">
                  ${startingPrice.toFixed(0)}
                  <span className="text-sm font-normal text-muted"> / night</span>
                </p>
              ) : (
                <p className="mt-1 text-sm font-medium text-muted">
                  Price on request
                </p>
              )}

              <button
                onClick={scrollToRooms}
                disabled={rooms.length === 0}
                className="mt-5 flex w-full items-center justify-center gap-2 rounded-pill bg-teal py-3.5 text-sm font-semibold text-white transition-colors hover:bg-teal-dark disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
              >
                {rooms.length === 0 ? "No rooms yet" : "Choose a room"}
                {rooms.length > 0 && <ArrowDown size={15} />}
              </button>

              <div className="mt-5 flex flex-col gap-2 border-t border-line pt-5 text-sm text-muted">
                <p>
                  {availableRoomCount > 0
                    ? `${availableRoomCount} available now`
                    : "Check back for availability"}
                </p>
                {hotel.stars > 0 && <p>{hotel.stars}-star property</p>}
              </div>
            </div>
          </aside>
        </div>
      </div>

      <Footer />
    </div>
  );
}
