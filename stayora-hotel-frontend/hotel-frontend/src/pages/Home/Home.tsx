import { Navbar } from "@/components/Navbar/Navbar";
import { Hero } from "@/components/Hero/Hero";
import { HotelSection } from "@/components/HotelSection/HotelSection";
import { WhyNileStay } from "@/components/WhyNileStay/WhyNileStay";
import { Footer } from "@/components/common/Footer";

export function Home() {
  return (
    <div className="min-h-screen bg-bg">
      <Navbar />
      <Hero />
      <HotelSection />
      <WhyNileStay />
      <Footer />
    </div>
  );
}
