import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { MapPin, CalendarDays, Users, Search } from "lucide-react";

interface FieldProps {
  icon: React.ReactNode;
  label: string;
  children: React.ReactNode;
  className?: string;
}

function Field({ icon, label, children, className = "" }: FieldProps) {
  return (
    <div className={`flex min-w-0 items-center gap-3 px-5 py-3.5 ${className}`}>
      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
        {icon}
      </span>
      <div className="min-w-0 flex-1">
        <p className="text-xs font-medium text-muted">{label}</p>
        {children}
      </div>
    </div>
  );
}

export function SearchBar() {
  const navigate = useNavigate();
  const [destination, setDestination] = useState("Cairo, Egypt");
  const [checkIn, setCheckIn] = useState("");
  const [checkOut, setCheckOut] = useState("");
  const [guests, setGuests] = useState(2);

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();

    // The backend's /hotels/search only filters by city/name/stars/price/
    // room_type/min_capacity - there's no date-based availability filter,
    // so dates are carried along for display only (see Search page).
    const city = destination.split(",")[0]?.trim();

    const params = new URLSearchParams();
    if (city) params.set("city", city);
    if (guests > 0) params.set("min_capacity", String(guests));
    if (checkIn) params.set("check_in", checkIn);
    if (checkOut) params.set("check_out", checkOut);

    navigate(`/search?${params.toString()}`);
  }

  return (
    <form
      onSubmit={handleSearch}
      className="w-full rounded-panel border border-line bg-white shadow-[var(--shadow-hover)]"
    >
      <div className="flex flex-col divide-y divide-line md:flex-row md:items-center md:divide-x md:divide-y-0">
        <Field icon={<MapPin size={17} />} label="Destination" className="md:w-[26%]">
          <input
            value={destination}
            onChange={(e) => setDestination(e.target.value)}
            placeholder="Where are you going?"
            className="w-full truncate rounded-[6px] bg-transparent text-sm font-semibold text-ink outline-none placeholder:text-muted placeholder:font-normal focus-visible:ring-2 focus-visible:ring-teal/40"
          />
        </Field>

        <Field icon={<CalendarDays size={17} />} label="Check in" className="md:w-[20%]">
          <input
            type="date"
            value={checkIn}
            onChange={(e) => setCheckIn(e.target.value)}
            className="w-full rounded-[6px] bg-transparent text-sm font-semibold text-ink outline-none [color-scheme:light] focus-visible:ring-2 focus-visible:ring-teal/40"
          />
        </Field>

        <Field icon={<CalendarDays size={17} />} label="Check out" className="md:w-[20%]">
          <input
            type="date"
            value={checkOut}
            onChange={(e) => setCheckOut(e.target.value)}
            className="w-full rounded-[6px] bg-transparent text-sm font-semibold text-ink outline-none [color-scheme:light] focus-visible:ring-2 focus-visible:ring-teal/40"
          />
        </Field>

        <Field icon={<Users size={17} />} label="Guests" className="md:w-[24%]">
          <input
            type="number"
            min={1}
            value={guests}
            onChange={(e) => setGuests(Number(e.target.value))}
            className="w-full truncate rounded-[6px] bg-transparent text-sm font-semibold text-ink outline-none focus-visible:ring-2 focus-visible:ring-teal/40"
          />
        </Field>

        <div className="p-3 md:pl-2">
          <button
            type="submit"
            className="flex w-full items-center justify-center gap-2 rounded-pill bg-teal px-6 py-3.5 text-sm font-semibold text-white transition-all duration-150 hover:bg-teal-dark active:scale-[0.97] md:w-14 md:p-0 md:py-3.5 cursor-pointer"
            aria-label="Search"
          >
            <Search size={18} />
            <span className="md:hidden">Search</span>
          </button>
        </div>
      </div>
    </form>
  );
}
