import { useEffect, useState } from "react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Hero } from "@/components/Hero/Hero";
import { HotelSection } from "@/components/HotelSection/HotelSection";
import { PopularDestinations } from "@/components/HotelSection/PopularDestinations";
import { WhyNileStay } from "@/components/WhyNileStay/WhyNileStay";
import { Footer } from "@/components/common/Footer";
import { getHotels } from "@/api/hotelApi";
import { getMyFavorites } from "@/api/favoriteApi";
import { extractErrorMessage } from "@/api/client";
import { useAuth } from "@/context/AuthContext";
import { getHotelImageUrl } from "@/utils/hotelImages";
import type { Hotel } from "@/types/hotel";

export type HotelListStatus = "loading" | "success" | "empty" | "error";

// The hotel list is fetched once here (not per-section) so Hero's featured
// image, the "Stays ready to book" grid, and "Popular destinations" all
// read from the same request instead of each hitting /hotels
// independently - see HotelSection's previous version, which used to own
// this fetch by itself.
export function Home() {
  const { isAuthenticated } = useAuth();
  const [hotels, setHotels] = useState<Hotel[]>([]);
  const [favoriteIds, setFavoriteIds] = useState<Set<number> | null>(null);
  const [status, setStatus] = useState<HotelListStatus>("loading");
  const [errorMessage, setErrorMessage] = useState("");

  async function load() {
    setStatus("loading");
    try {
      const response = await getHotels();
      if (response.success) {
        if (response.data.length === 0) {
          setStatus("empty");
        } else {
          setHotels(response.data);
          setStatus("success");
        }
      } else {
        setErrorMessage(response.message || "Couldn't load hotels right now.");
        setStatus("error");
      }
    } catch (err) {
      setErrorMessage(extractErrorMessage(err));
      setStatus("error");
    }

    if (isAuthenticated) {
      try {
        const favRes = await getMyFavorites();
        if (favRes.success) {
          setFavoriteIds(new Set(favRes.data.map((f) => f.hotel_id)));
        }
      } catch {
        setFavoriteIds(new Set());
      }
    } else {
      setFavoriteIds(new Set());
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const featuredHotel = hotels.find((h) => getHotelImageUrl(h)) ?? null;

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />
      <Hero featuredHotel={featuredHotel} />
      <HotelSection
        hotels={hotels}
        status={status}
        errorMessage={errorMessage}
        favoriteIds={favoriteIds}
        onRetry={load}
      />
      <PopularDestinations hotels={hotels} />
      <WhyNileStay />
      <Footer />
    </div>
  );
}
