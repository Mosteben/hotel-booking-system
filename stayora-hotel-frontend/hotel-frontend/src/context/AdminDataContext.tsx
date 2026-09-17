import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type Dispatch,
  type ReactNode,
  type SetStateAction,
} from "react";
import { getAllBookings } from "@/api/bookingApi";
import { getAllPayments } from "@/api/paymentApi";
import { getAllReviews } from "@/api/reviewApi";
import { getAllUsers } from "@/api/userApi";
import { getHotels } from "@/api/hotelApi";
import { getRoomsByHotelID } from "@/api/roomApi";
import { extractErrorMessage } from "@/api/client";
import type { AdminBookingSummary } from "@/types/booking";
import type { AdminPaymentSummary } from "@/types/payment";
import type { AdminReviewSummary } from "@/types/review";
import type { AdminUserSummary } from "@/types/user";
import type { Hotel } from "@/types/hotel";

export interface AdminHotelSummary extends Hotel {
  roomCount: number;
}

// Per-resource errors, not one blanket flag - a failed /users request must
// never be indistinguishable from "there are genuinely zero users". Each
// resource that fails to load keeps its own message here; the ones that
// succeeded still populate normally even if a sibling request failed.
export interface AdminDataErrors {
  bookings?: string;
  payments?: string;
  reviews?: string;
  users?: string;
  hotels?: string;
}

interface AdminDataContextValue {
  bookings: AdminBookingSummary[];
  payments: AdminPaymentSummary[];
  reviews: AdminReviewSummary[];
  users: AdminUserSummary[];
  hotels: AdminHotelSummary[];
  isLoading: boolean;
  errors: AdminDataErrors;
  refresh: () => Promise<void>;
  setBookings: Dispatch<SetStateAction<AdminBookingSummary[]>>;
  setPayments: Dispatch<SetStateAction<AdminPaymentSummary[]>>;
  setReviews: Dispatch<SetStateAction<AdminReviewSummary[]>>;
}

const AdminDataContext = createContext<AdminDataContextValue | undefined>(undefined);

// Fetches every dataset the Admin area needs exactly once (on entering
// /admin), shared across the sidebar badges, the Overview dashboard, and
// every management page - so opening the Admin area doesn't re-fetch the
// same lists once per page. Individual pages still own their own mutation
// (confirm/cancel/delete) logic; they just read/update this shared state
// instead of holding a separate copy of the same data.
export function AdminDataProvider({ children }: { children: ReactNode }) {
  const [bookings, setBookings] = useState<AdminBookingSummary[]>([]);
  const [payments, setPayments] = useState<AdminPaymentSummary[]>([]);
  const [reviews, setReviews] = useState<AdminReviewSummary[]>([]);
  const [users, setUsers] = useState<AdminUserSummary[]>([]);
  const [hotels, setHotels] = useState<AdminHotelSummary[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [errors, setErrors] = useState<AdminDataErrors>({});

  const load = useCallback(async () => {
    setIsLoading(true);

    const nextErrors: AdminDataErrors = {};

    const [bookingsResult, paymentsResult, reviewsResult, usersResult, hotelsResult] =
      await Promise.allSettled([
        getAllBookings(),
        getAllPayments(),
        getAllReviews(),
        getAllUsers(),
        getHotels(),
      ]);

    if (bookingsResult.status === "fulfilled") {
      setBookings(bookingsResult.value);
    } else {
      setBookings([]);
      nextErrors.bookings = extractErrorMessage(bookingsResult.reason);
    }

    if (paymentsResult.status === "fulfilled" && paymentsResult.value.success) {
      setPayments(paymentsResult.value.data);
    } else {
      setPayments([]);
      nextErrors.payments =
        paymentsResult.status === "fulfilled"
          ? paymentsResult.value.message
          : extractErrorMessage(paymentsResult.reason);
    }

    if (reviewsResult.status === "fulfilled" && reviewsResult.value.success) {
      setReviews(reviewsResult.value.data);
    } else {
      setReviews([]);
      nextErrors.reviews =
        reviewsResult.status === "fulfilled"
          ? reviewsResult.value.message
          : extractErrorMessage(reviewsResult.reason);
    }

    if (usersResult.status === "fulfilled" && usersResult.value.success) {
      setUsers(usersResult.value.data);
    } else {
      setUsers([]);
      nextErrors.users =
        usersResult.status === "fulfilled"
          ? usersResult.value.message
          : extractErrorMessage(usersResult.reason);
    }

    if (hotelsResult.status === "fulfilled" && hotelsResult.value.success) {
      const hotelList = hotelsResult.value.data;

      // Room counts per hotel are needed for "hotels without rooms" (Needs
      // Attention) and for the Hotels list - bounded by hotel count (not
      // booking count), same shape as the room/hotel enrichment already
      // used elsewhere in this app.
      const withRoomCounts = await Promise.all(
        hotelList.map(async (hotel): Promise<AdminHotelSummary> => {
          try {
            const roomsRes = await getRoomsByHotelID(hotel.id);
            return {
              ...hotel,
              roomCount: roomsRes.success ? roomsRes.data.length : 0,
            };
          } catch {
            return { ...hotel, roomCount: 0 };
          }
        })
      );

      setHotels(withRoomCounts);
    } else {
      setHotels([]);
      nextErrors.hotels =
        hotelsResult.status === "fulfilled"
          ? hotelsResult.value.message
          : extractErrorMessage(hotelsResult.reason);
    }

    setErrors(nextErrors);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const value = useMemo<AdminDataContextValue>(
    () => ({
      bookings,
      payments,
      reviews,
      users,
      hotels,
      isLoading,
      errors,
      refresh: load,
      setBookings,
      setPayments,
      setReviews,
    }),
    [bookings, payments, reviews, users, hotels, isLoading, errors, load]
  );

  return (
    <AdminDataContext.Provider value={value}>{children}</AdminDataContext.Provider>
  );
}

export function useAdminData(): AdminDataContextValue {
  const ctx = useContext(AdminDataContext);
  if (!ctx) {
    throw new Error("useAdminData must be used within an AdminDataProvider");
  }
  return ctx;
}
