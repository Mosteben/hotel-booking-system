import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthLayout } from "@/components/auth/AuthLayout";
import { FormField, FormSelect } from "@/components/auth/FormField";
import { PasswordField } from "@/components/auth/PasswordField";
import { AuthButton } from "@/components/auth/AuthButton";
import { register as registerRequest } from "@/api/authApi";
import { extractErrorMessage } from "@/api/client";
import type { RegisterRequest } from "@/types/auth";

type FormState = RegisterRequest;

const INITIAL_STATE: FormState = {
  first_name: "",
  last_name: "",
  email: "",
  password: "",
  confirm_password: "",
  phone: "",
  date_of_birth: "",
  gender: "male",
  nationality: "",
  national_id: "",
  passport_number: "",
  state: "",
  postal_code: "",
  address: "",
  city: "",
  country: "",
};

const REQUIRED_FIELDS: (keyof FormState)[] = [
  "first_name",
  "last_name",
  "email",
  "password",
  "confirm_password",
  "phone",
  "date_of_birth",
  "gender",
  "nationality",
  "national_id",
  "state",
  "postal_code",
  "address",
  "city",
  "country",
];

const LABELS: Record<keyof FormState, string> = {
  first_name: "First name",
  last_name: "Last name",
  email: "Email",
  password: "Password",
  confirm_password: "Confirm password",
  phone: "Phone number",
  date_of_birth: "Date of birth",
  gender: "Gender",
  nationality: "Nationality",
  national_id: "National ID",
  passport_number: "Passport number",
  state: "State",
  postal_code: "Postal code",
  address: "Address",
  city: "City",
  country: "Country",
};

function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <legend className="mb-1 text-xs font-semibold uppercase tracking-wider text-teal-light">
      {children}
    </legend>
  );
}

export function Register() {
  const navigate = useNavigate();
  const [form, setForm] = useState<FormState>(INITIAL_STATE);
  const [errors, setErrors] = useState<Partial<Record<keyof FormState, string>>>({});
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  function update<K extends keyof FormState>(key: K, value: FormState[K]) {
    setForm((prev) => ({ ...prev, [key]: value }));
  }

  function validate(): boolean {
    const next: Partial<Record<keyof FormState, string>> = {};

    for (const field of REQUIRED_FIELDS) {
      if (!form[field] || !String(form[field]).trim()) {
        next[field] = `${LABELS[field]} is required.`;
      }
    }

    if (form.email && !/^\S+@\S+\.\S+$/.test(form.email)) {
      next.email = "Enter a valid email.";
    }
    if (form.password && form.password.length < 8) {
      next.password = "Password must be at least 8 characters.";
    }
    if (form.confirm_password && form.password !== form.confirm_password) {
      next.confirm_password = "Passwords do not match.";
    }

    setErrors(next);
    return Object.keys(next).length === 0;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setFormError("");
    if (!validate()) return;

    setIsSubmitting(true);
    try {
      const payload: RegisterRequest = {
        ...form,
        passport_number: form.passport_number?.trim() || "",
      };
      const response = await registerRequest(payload);
      if (response.success) {
        navigate("/login", {
          state: { justRegistered: true, email: form.email },
        });
      } else {
        setFormError(response.message || "Registration failed. Please try again.");
      }
    } catch (err) {
      setFormError(extractErrorMessage(err));
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <AuthLayout imageSide="left">
      <p className="text-sm font-semibold tracking-wide text-teal">
        Join NileStay
      </p>
      <h1 className="mt-1.5 font-display text-3xl font-bold text-ink">
        Create your account
      </h1>
      <p className="mt-2 text-sm leading-relaxed text-muted">
        Join us and start discovering your next stay.
      </p>

      <form onSubmit={handleSubmit} noValidate className="mt-8 flex flex-col gap-7">
        <fieldset className="flex flex-col gap-4">
          <SectionLabel>Personal information</SectionLabel>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={LABELS.first_name}
              autoComplete="given-name"
              value={form.first_name}
              onChange={(e) => update("first_name", e.target.value)}
              error={errors.first_name}
            />
            <FormField
              label={LABELS.last_name}
              autoComplete="family-name"
              value={form.last_name}
              onChange={(e) => update("last_name", e.target.value)}
              error={errors.last_name}
            />
          </div>
          <FormField
            label={LABELS.email}
            type="email"
            autoComplete="email"
            value={form.email}
            onChange={(e) => update("email", e.target.value)}
            error={errors.email}
          />
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={LABELS.phone}
              type="tel"
              autoComplete="tel"
              value={form.phone}
              onChange={(e) => update("phone", e.target.value)}
              error={errors.phone}
            />
            <FormField
              label={LABELS.date_of_birth}
              type="date"
              autoComplete="bday"
              value={form.date_of_birth}
              onChange={(e) => update("date_of_birth", e.target.value)}
              error={errors.date_of_birth}
            />
          </div>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormSelect
              label={LABELS.gender}
              value={form.gender}
              onChange={(e) => update("gender", e.target.value as FormState["gender"])}
              error={errors.gender}
            >
              <option value="male">Male</option>
              <option value="female">Female</option>
            </FormSelect>
            <FormField
              label={LABELS.nationality}
              value={form.nationality}
              onChange={(e) => update("nationality", e.target.value)}
              error={errors.nationality}
            />
          </div>
        </fieldset>

        <fieldset className="flex flex-col gap-4 border-t border-line pt-6">
          <SectionLabel>Address</SectionLabel>
          <FormField
            label={LABELS.address}
            autoComplete="street-address"
            value={form.address}
            onChange={(e) => update("address", e.target.value)}
            error={errors.address}
          />
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={LABELS.city}
              autoComplete="address-level2"
              value={form.city}
              onChange={(e) => update("city", e.target.value)}
              error={errors.city}
            />
            <FormField
              label={LABELS.state}
              autoComplete="address-level1"
              value={form.state}
              onChange={(e) => update("state", e.target.value)}
              error={errors.state}
            />
          </div>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={LABELS.postal_code}
              autoComplete="postal-code"
              value={form.postal_code}
              onChange={(e) => update("postal_code", e.target.value)}
              error={errors.postal_code}
            />
            <FormField
              label={LABELS.country}
              autoComplete="country-name"
              value={form.country}
              onChange={(e) => update("country", e.target.value)}
              error={errors.country}
            />
          </div>
        </fieldset>

        <fieldset className="flex flex-col gap-4 border-t border-line pt-6">
          <SectionLabel>Identity</SectionLabel>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField
              label={LABELS.national_id}
              value={form.national_id}
              onChange={(e) => update("national_id", e.target.value)}
              error={errors.national_id}
            />
            <FormField
              label={LABELS.passport_number}
              optional
              value={form.passport_number}
              onChange={(e) => update("passport_number", e.target.value)}
              error={errors.passport_number}
            />
          </div>
        </fieldset>

        <fieldset className="flex flex-col gap-4 border-t border-line pt-6">
          <SectionLabel>Security</SectionLabel>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <PasswordField
              label={LABELS.password}
              autoComplete="new-password"
              value={form.password}
              onChange={(e) => update("password", e.target.value)}
              error={errors.password}
            />
            <PasswordField
              label={LABELS.confirm_password}
              autoComplete="new-password"
              value={form.confirm_password}
              onChange={(e) => update("confirm_password", e.target.value)}
              error={errors.confirm_password}
            />
          </div>
        </fieldset>

        {formError && (
          <p
            role="alert"
            className="rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
          >
            {formError}
          </p>
        )}

        <AuthButton isSubmitting={isSubmitting}>Create account</AuthButton>
      </form>

      <p className="mt-8 text-center text-sm text-muted">
        Already have an account?{" "}
        <Link
          to="/login"
          className="font-semibold text-teal transition-colors duration-200 hover:text-teal-dark"
        >
          Log in
        </Link>
      </p>
    </AuthLayout>
  );
}
