import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "@/context/AuthContext";
import { ProtectedRoute } from "@/components/auth/ProtectedRoute";
import { Home } from "@/pages/Home/Home";
import { Login } from "@/pages/Login/Login";
import { Register } from "@/pages/Register/Register";
import { HotelDetails } from "@/pages/HotelDetails/HotelDetails";
import { Search } from "@/pages/Search/Search";
import { Booking } from "@/pages/Booking/Booking";
import { Payment } from "@/pages/Payment/Payment";
import { MyBookings } from "@/pages/MyBookings/MyBookings";
import { BookingDetails } from "@/pages/BookingDetails/BookingDetails";
import { Favorites } from "@/pages/Favorites/Favorites";
import { Profile } from "@/pages/Profile/Profile";
import { AdminLayout } from "@/layouts/AdminLayout/AdminLayout";
import { Overview } from "@/pages/Admin/Overview";
import { HotelList } from "@/pages/Admin/HotelList";
import { HotelForm } from "@/pages/Admin/HotelForm";
import { RoomList } from "@/pages/Admin/RoomList";
import { BookingList } from "@/pages/Admin/BookingList";
import { PaymentList } from "@/pages/Admin/PaymentList";
import { UserList } from "@/pages/Admin/UserList";
import { ReviewList } from "@/pages/Admin/ReviewList";

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/search" element={<Search />} />
          <Route path="/hotels/:id" element={<HotelDetails />} />

          <Route
            path="/booking/:roomId"
            element={
              <ProtectedRoute>
                <Booking />
              </ProtectedRoute>
            }
          />
          <Route
            path="/bookings/:id/payment"
            element={
              <ProtectedRoute>
                <Payment />
              </ProtectedRoute>
            }
          />
          <Route
            path="/my-bookings"
            element={
              <ProtectedRoute>
                <MyBookings />
              </ProtectedRoute>
            }
          />
          <Route
            path="/bookings/:id"
            element={
              <ProtectedRoute>
                <BookingDetails />
              </ProtectedRoute>
            }
          />
          <Route
            path="/favorites"
            element={
              <ProtectedRoute>
                <Favorites />
              </ProtectedRoute>
            }
          />
          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <Profile />
              </ProtectedRoute>
            }
          />

          {/* Admin: its own layout (sidebar + topbar), nested under one
              role-gated parent route rather than repeating ProtectedRoute
              on every admin page. */}
          <Route
            path="/admin"
            element={
              <ProtectedRoute roles={["admin", "manager"]}>
                <AdminLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<Overview />} />
            <Route path="bookings" element={<BookingList />} />
            <Route path="payments" element={<PaymentList />} />
            <Route path="hotels" element={<HotelList />} />
            <Route path="hotels/new" element={<HotelForm />} />
            <Route path="hotels/:id/edit" element={<HotelForm />} />
            <Route path="hotels/:hotelId/rooms" element={<RoomList />} />
            <Route path="users" element={<UserList />} />
            <Route path="reviews" element={<ReviewList />} />
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
