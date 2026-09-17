// These shapes mirror the real backend request/response bodies exactly
// (verified against Swagger + live /auth responses). Do not rename fields.

export type Gender = "male" | "female";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  confirm_password: string;
  phone: string;
  date_of_birth: string; // "YYYY-MM-DD"
  gender: Gender;
  nationality: string;
  national_id: string;
  passport_number?: string;
  state: string;
  postal_code: string;
  address: string;
  city: string;
  country: string;
}

export interface LoginResponseData {
  token: string;
}

// Registration currently only returns { success, message } — no token, no user.
export type RegisterResponseData = undefined;

// The nested profile object comes back with Go's default (PascalCase) field
// names because the backend struct has no custom json tags on it.
export interface UserProfile {
  ID: number;
  UserID: string;
  DateOfBirth: string;
  Gender: string;
  Nationality: string;
  NationalID: string;
  PassportNumber: string;
  Address: string;
  City: string;
  State: string;
  Country: string;
  PostalCode: string;
  Avatar: string;
  Bio: string;
  CreatedAt: string;
  UpdatedAt: string;
}

export type UserRole = "customer" | "manager" | "admin";

export interface AuthenticatedUser {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
  is_active: boolean;
  is_verified: boolean;
  role: UserRole;
  profile: UserProfile;
}

// PUT /auth/profile - every field is optional (validate:"omitempty,...").
export interface UpdateProfileRequest {
  first_name?: string;
  last_name?: string;
  phone?: string;
  gender?: Gender;
  nationality?: string;
  national_id?: string;
  passport_number?: string;
  address?: string;
  city?: string;
  state?: string;
  country?: string;
  postal_code?: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
  confirm_password: string;
}
