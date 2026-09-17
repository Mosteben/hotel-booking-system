import { useEffect, useState, type FormEvent, type ReactNode } from "react";
import { Link } from "react-router-dom";
import {
  UserRound,
  KeyRound,
  MapPin,
  CalendarCheck,
  Heart,
  ChevronRight,
} from "lucide-react";
import { Navbar } from "@/components/Navbar/Navbar";
import { Footer } from "@/components/common/Footer";
import { PageHeading } from "@/components/common/PageHeading";
import { FormField, FormSelect } from "@/components/auth/FormField";
import { PasswordField } from "@/components/auth/PasswordField";
import { AuthButton } from "@/components/auth/AuthButton";
import { updateProfile, changePassword } from "@/api/authApi";
import { extractErrorMessage } from "@/api/client";
import { useAuth } from "@/context/AuthContext";
import type { ChangePasswordRequest, UpdateProfileRequest } from "@/types/auth";

const EMPTY_PROFILE_FORM: UpdateProfileRequest = {
  first_name: "",
  last_name: "",
  phone: "",
  gender: "male",
  nationality: "",
  national_id: "",
  passport_number: "",
  address: "",
  city: "",
  state: "",
  country: "",
  postal_code: "",
};

const EMPTY_PASSWORD_FORM: ChangePasswordRequest = {
  current_password: "",
  new_password: "",
  confirm_password: "",
};

const QUICK_LINKS = [
  { label: "My bookings", to: "/my-bookings", icon: CalendarCheck },
  { label: "Favorites", to: "/favorites", icon: Heart },
];

function SectionCard({
  icon: Icon,
  title,
  description,
  children,
}: {
  icon: typeof UserRound;
  title: string;
  description: string;
  children: ReactNode;
}) {
  return (
    <div className="rounded-card border border-line bg-white p-8 sm:p-10">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
          <Icon size={19} />
        </span>
        <div>
          <h2 className="font-display text-lg font-semibold text-ink">{title}</h2>
          <p className="mt-0.5 text-sm text-muted">{description}</p>
        </div>
      </div>
      <div className="mt-8 flex flex-col gap-6">{children}</div>
    </div>
  );
}

export function Profile() {
  const { user, refreshUser } = useAuth();

  const [form, setForm] = useState<UpdateProfileRequest>(EMPTY_PROFILE_FORM);
  const [profileMessage, setProfileMessage] = useState("");
  const [profileError, setProfileError] = useState("");
  const [isSavingProfile, setIsSavingProfile] = useState(false);

  const [passwordForm, setPasswordForm] =
    useState<ChangePasswordRequest>(EMPTY_PASSWORD_FORM);
  const [passwordMessage, setPasswordMessage] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [isSavingPassword, setIsSavingPassword] = useState(false);

  useEffect(() => {
    if (!user) return;
    setForm({
      first_name: user.first_name,
      last_name: user.last_name,
      phone: user.phone,
      gender: (user.profile?.Gender as "male" | "female") || "male",
      nationality: user.profile?.Nationality || "",
      national_id: user.profile?.NationalID || "",
      passport_number: user.profile?.PassportNumber || "",
      address: user.profile?.Address || "",
      city: user.profile?.City || "",
      state: user.profile?.State || "",
      country: user.profile?.Country || "",
      postal_code: user.profile?.PostalCode || "",
    });
  }, [user]);

  function update<K extends keyof UpdateProfileRequest>(
    key: K,
    value: UpdateProfileRequest[K]
  ) {
    setForm((prev) => ({ ...prev, [key]: value }));
  }

  async function handleProfileSubmit(e: FormEvent) {
    e.preventDefault();
    setProfileMessage("");
    setProfileError("");
    setIsSavingProfile(true);
    try {
      const response = await updateProfile(form);
      if (response.success) {
        setProfileMessage("Profile updated successfully.");
        await refreshUser();
      } else {
        setProfileError(response.message || "Couldn't update your profile.");
      }
    } catch (err) {
      setProfileError(extractErrorMessage(err));
    } finally {
      setIsSavingProfile(false);
    }
  }

  async function handlePasswordSubmit(e: FormEvent) {
    e.preventDefault();
    setPasswordMessage("");
    setPasswordError("");

    if (passwordForm.new_password !== passwordForm.confirm_password) {
      setPasswordError("New passwords do not match.");
      return;
    }
    if (passwordForm.new_password.length < 8) {
      setPasswordError("New password must be at least 8 characters.");
      return;
    }

    setIsSavingPassword(true);
    try {
      const response = await changePassword(passwordForm);
      if (response.success) {
        setPasswordMessage("Password changed successfully.");
        setPasswordForm(EMPTY_PASSWORD_FORM);
      } else {
        setPasswordError(response.message || "Couldn't change your password.");
      }
    } catch (err) {
      setPasswordError(extractErrorMessage(err));
    } finally {
      setIsSavingPassword(false);
    }
  }

  if (!user) return null;

  return (
    <div className="min-h-screen bg-bg">
      <Navbar />

      <div className="page-fade-in mx-auto max-w-4xl px-4 py-10 sm:px-6 sm:py-14">
        <PageHeading eyebrow="Account" title="Your profile" />

        {/* Identity header */}
        <div className="mt-8 flex flex-col items-start gap-6 rounded-card border border-line bg-white p-8 sm:flex-row sm:items-center sm:p-10">
          <span className="flex h-20 w-20 shrink-0 items-center justify-center rounded-full bg-teal font-display text-2xl font-bold text-white">
            {user.first_name[0]}
            {user.last_name[0]}
          </span>
          <div>
            <h1 className="font-display text-2xl font-bold text-ink">
              {user.first_name} {user.last_name}
            </h1>
            <p className="mt-1.5 text-sm text-muted">{user.email}</p>
            <span className="mt-3 inline-flex items-center rounded-pill bg-teal/10 px-3 py-1 text-xs font-semibold capitalize text-teal">
              {user.role}
            </span>
          </div>
        </div>

        {/* Quick links */}
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2">
          {QUICK_LINKS.map(({ label, to, icon: Icon }) => (
            <Link
              key={to}
              to={to}
              className="flex items-center justify-between rounded-card border border-line bg-white px-6 py-5 transition-colors hover:border-teal/30"
            >
              <span className="flex items-center gap-3 text-sm font-semibold text-ink">
                <span className="flex h-10 w-10 items-center justify-center rounded-full bg-bg text-teal">
                  <Icon size={16} />
                </span>
                {label}
              </span>
              <ChevronRight size={16} className="text-muted" />
            </Link>
          ))}
        </div>

        {/* Personal information + Contact & address: one submit action,
            but rendered as two clearly separate cards rather than one long
            list of fields. */}
        <form onSubmit={handleProfileSubmit} className="mt-6 flex flex-col gap-6">
          <SectionCard
            icon={UserRound}
            title="Personal information"
            description="Your identity on NileStay."
          >
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <FormField
                label="First name"
                value={form.first_name}
                onChange={(e) => update("first_name", e.target.value)}
              />
              <FormField
                label="Last name"
                value={form.last_name}
                onChange={(e) => update("last_name", e.target.value)}
              />
            </div>
            <FormField label="Email" value={user.email} disabled className="opacity-60" />
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <FormField
                label="Phone number"
                type="tel"
                value={form.phone}
                onChange={(e) => update("phone", e.target.value)}
              />
              <FormSelect
                label="Gender"
                value={form.gender}
                onChange={(e) => update("gender", e.target.value as "male" | "female")}
              >
                <option value="male">Male</option>
                <option value="female">Female</option>
              </FormSelect>
            </div>
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <FormField
                label="Nationality"
                value={form.nationality}
                onChange={(e) => update("nationality", e.target.value)}
              />
              <FormField
                label="National ID"
                value={form.national_id}
                onChange={(e) => update("national_id", e.target.value)}
              />
            </div>
            <FormField
              label="Passport number"
              optional
              value={form.passport_number}
              onChange={(e) => update("passport_number", e.target.value)}
            />
          </SectionCard>

          <SectionCard
            icon={MapPin}
            title="Contact & address"
            description="Where we can reach you."
          >
            <FormField
              label="Address"
              value={form.address}
              onChange={(e) => update("address", e.target.value)}
            />
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <FormField
                label="City"
                value={form.city}
                onChange={(e) => update("city", e.target.value)}
              />
              <FormField
                label="State"
                value={form.state}
                onChange={(e) => update("state", e.target.value)}
              />
            </div>
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <FormField
                label="Postal code"
                value={form.postal_code}
                onChange={(e) => update("postal_code", e.target.value)}
              />
              <FormField
                label="Country"
                value={form.country}
                onChange={(e) => update("country", e.target.value)}
              />
            </div>
          </SectionCard>

          <div className="rounded-card border border-line bg-white p-8 sm:p-10">
            {profileMessage && (
              <p className="mb-6 rounded-xl border border-emerald-100 bg-emerald-50 px-4 py-3 text-sm font-medium text-emerald-600">
                {profileMessage}
              </p>
            )}
            {profileError && (
              <p role="alert" className="mb-6 rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600">
                {profileError}
              </p>
            )}
            <AuthButton isSubmitting={isSavingProfile} className="w-auto px-8 sm:w-auto">
              Save changes
            </AuthButton>
          </div>
        </form>

        {/* Security */}
        <form
          onSubmit={handlePasswordSubmit}
          className="mt-6 rounded-card border border-line bg-white p-8 sm:p-10"
        >
          <div className="flex items-center gap-4">
            <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-bg text-teal">
              <KeyRound size={19} />
            </span>
            <div>
              <h2 className="font-display text-lg font-semibold text-ink">
                Account settings
              </h2>
              <p className="mt-0.5 text-sm text-muted">Change your password.</p>
            </div>
          </div>

          <div className="mt-8 flex flex-col gap-6">
            <PasswordField
              label="Current password"
              autoComplete="current-password"
              value={passwordForm.current_password}
              onChange={(e) =>
                setPasswordForm((prev) => ({ ...prev, current_password: e.target.value }))
              }
            />
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
              <PasswordField
                label="New password"
                autoComplete="new-password"
                value={passwordForm.new_password}
                onChange={(e) =>
                  setPasswordForm((prev) => ({ ...prev, new_password: e.target.value }))
                }
              />
              <PasswordField
                label="Confirm new password"
                autoComplete="new-password"
                value={passwordForm.confirm_password}
                onChange={(e) =>
                  setPasswordForm((prev) => ({ ...prev, confirm_password: e.target.value }))
                }
              />
            </div>

            {passwordMessage && (
              <p className="rounded-xl border border-emerald-100 bg-emerald-50 px-4 py-3 text-sm font-medium text-emerald-600">
                {passwordMessage}
              </p>
            )}
            {passwordError && (
              <p role="alert" className="rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600">
                {passwordError}
              </p>
            )}

            <AuthButton isSubmitting={isSavingPassword} className="w-auto self-start px-8">
              Update password
            </AuthButton>
          </div>
        </form>
      </div>

      <Footer />
    </div>
  );
}
