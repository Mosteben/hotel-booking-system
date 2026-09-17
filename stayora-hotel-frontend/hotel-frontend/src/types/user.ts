// Mirrors GET /users exactly (internal/user/service AdminUserSummary) -
// deliberately a narrower, safe view of the user, never the raw model.
// There is no password field here at all, by design.

import type { UserRole } from "@/types/auth";

export interface AdminUserSummary {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  role: UserRole;
  is_active: boolean;
  is_verified: boolean;
  created_at: string;
}
