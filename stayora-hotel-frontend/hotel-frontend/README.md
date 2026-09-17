# NileStay — Hotel Booking Frontend (Phase 1)

React + Vite + TypeScript + Tailwind CSS frontend for the Hotel Booking
System. This phase implements exactly three pages and connects to the
**real** Go/Gin backend — there is no mock data anywhere in this codebase.

## Pages in this phase

- `/` — Landing / Home (hero, search bar, categories, live hotel list from `GET /hotels`)
- `/login` — Login (`POST /auth/login`)
- `/register` — Register (`POST /auth/register`)

Rooms, bookings, payments, reviews, favorites, profile, and dashboard/admin
pages are intentionally **not** included — they come in a later phase.

## Getting started

```bash
npm install
cp .env.example .env   # already points at http://localhost:8081
npm run dev
```

Make sure the Go backend is running on `http://localhost:8081` (its CORS
config already allows `http://localhost:5173`, so no backend changes are
needed).

Open http://localhost:5173.

## Environment variables

```env
VITE_API_URL=http://localhost:8081
```

Change this if your backend runs on a different host/port.

## Project structure

```
src/
├── api/
│   ├── client.ts        # Axios instance, JWT storage, error normalization
│   ├── authApi.ts       # login / register / getCurrentUser
│   └── hotelApi.ts      # getHotels
├── components/
│   ├── Navbar/
│   ├── Hero/
│   ├── SearchBar/
│   ├── Categories/
│   ├── HotelCard/        # card + skeleton
│   ├── HotelSection/      # loading / empty / error states, real data only
│   ├── auth/              # AuthLayout, FormField, FormSelect
│   └── common/            # Footer
├── pages/
│   ├── Home/
│   ├── Login/
│   └── Register/
├── context/
│   └── AuthContext.tsx    # JWT-based auth state
├── types/
│   ├── api.ts     # generic ApiResponse envelope
│   ├── auth.ts    # matches real /auth request & response shapes
│   └── hotel.ts   # matches real GET /hotels response shape
└── utils/
    └── hotelImages.ts   # centralized presentation-image mapping
                          # (backend has no image_url field yet — swap
                          # this file's internals when it does)
```

## API contracts implemented (verified against the real backend)

**POST /auth/login**
```json
// request
{ "email": "string", "password": "string" }
// response
{ "success": true, "message": "Login successful", "data": { "token": "..." } }
```

**POST /auth/register**
```json
// request
{
  "first_name": "string", "last_name": "string", "email": "string",
  "password": "string", "confirm_password": "string", "phone": "string",
  "date_of_birth": "YYYY-MM-DD", "gender": "male",
  "nationality": "string", "national_id": "string",
  "passport_number": "string", "state": "string", "postal_code": "string",
  "address": "string", "city": "string", "country": "string"
}
// response — NOTE: no token is returned, so the app redirects to /login
{ "success": true, "message": "User registered successfully" }
```

**GET /auth/me**
```json
{
  "success": true,
  "data": {
    "id": "...", "email": "...", "first_name": "...", "last_name": "...",
    "phone": "...", "is_active": true, "is_verified": false,
    "role": "customer",
    "profile": { "DateOfBirth": "...", "Gender": "...", "...": "..." }
  }
}
```

**GET /hotels**
```json
{
  "success": true,
  "message": "hotels retrieved successfully",
  "data": [
    {
      "id": 2, "name": "Grand Nile Hotel", "description": "...",
      "address": "...", "city": "Cairo", "country": "Egypt",
      "phone": "...", "email": "...", "stars": 5, "rooms": null,
      "created_at": "...", "updated_at": "..."
    }
  ]
}
```

Only these fields are shown on hotel cards. No price or rating is
displayed because the backend doesn't return them.

## Authentication

- JWT is stored in `localStorage` under `stayora_token` and attached to
  every request via an Axios request interceptor (`Authorization: Bearer <token>`).
- On app load, if a token exists, `AuthContext` calls `GET /auth/me` to
  restore the session; if that fails, the stale token is cleared.
- Register does **not** log the user in automatically (the backend
  doesn't return a token on register) — it redirects to `/login`.

## Design system

| Token | Value |
|---|---|
| Primary (teal) | `#0B3F42` |
| Background | `#F7F8F6` |
| Main text | `#123234` |
| Secondary text | `#7C8A8B` |
| Border | `#E7ECEB` |
| Typography | Manrope (display) / Inter (body) |

Tokens live in `src/index.css` under Tailwind v4's `@theme` block
(`bg-teal`, `text-ink`, `rounded-panel`, `rounded-pill`, etc.).

## Verified

- `npx tsc -b --noEmit` — 0 errors
- `npm run build` — production build succeeds
- `npm run dev` — all three routes (`/`, `/login`, `/register`) return `200`
- No mock/fake data anywhere — hotel list and auth flows call the real backend
