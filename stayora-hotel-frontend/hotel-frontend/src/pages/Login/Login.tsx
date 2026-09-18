import { useState, type FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import { AuthLayout } from "@/components/auth/AuthLayout";
import { FormField } from "@/components/auth/FormField";
import { PasswordField } from "@/components/auth/PasswordField";
import { AuthButton } from "@/components/auth/AuthButton";
import { login as loginRequest } from "@/api/authApi";
import { extractErrorMessage } from "@/api/client";
import { useAuth } from "@/context/AuthContext";

export function Login() {
  const navigate = useNavigate();
  const { loginWithToken } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [formError, setFormError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  function validate(): boolean {
    const errors: Record<string, string> = {};
    if (!email.trim()) errors.email = "Email is required.";
    else if (!/^\S+@\S+\.\S+$/.test(email)) errors.email = "Enter a valid email.";
    if (!password) errors.password = "Password is required.";
    setFieldErrors(errors);
    return Object.keys(errors).length === 0;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setFormError("");
    if (!validate()) return;

    setIsSubmitting(true);
    try {
      const response = await loginRequest({ email, password });
      if (response.success) {
        await loginWithToken(response.data.token);
        navigate("/");
      } else {
        setFormError(response.message || "Invalid email or password.");
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
        Welcome back
      </p>
      <h1 className="mt-1.5 font-display text-3xl font-bold text-ink">
        Sign in to NileStay
      </h1>
      <p className="mt-2 text-sm leading-relaxed text-muted">
        Sign in to continue your journey and discover your perfect stay.
      </p>

      <form onSubmit={handleSubmit} noValidate className="mt-8 flex flex-col gap-5">
        <FormField
          label="Email"
          type="email"
          autoComplete="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          error={fieldErrors.email}
        />

        <PasswordField
          label="Password"
          autoComplete="current-password"
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          error={fieldErrors.password}
        />

        {formError && (
          <p
            role="alert"
            className="rounded-xl border border-red-100 bg-red-50 px-4 py-3 text-sm font-medium text-red-600"
          >
            {formError}
          </p>
        )}

        <AuthButton isSubmitting={isSubmitting}>Sign in</AuthButton>
      </form>

      {/* Social login isn't implemented in the backend yet — shown
          disabled rather than as fake working buttons. Delete this block
          entirely if you'd rather not show them at all. */}
      <div className="mt-7">
        <div className="flex items-center gap-3">
          <span className="h-px flex-1 bg-line" />
          <span className="text-xs font-medium uppercase tracking-wider text-muted">
            or
          </span>
          <span className="h-px flex-1 bg-line" />
        </div>
        <div className="mt-4 flex justify-center gap-3">
          {["Facebook", "Google"].map((name) => (
            <button
              key={name}
              type="button"
              disabled
              title="Coming soon"
              className="flex h-11 w-11 cursor-not-allowed items-center justify-center rounded-[13px] border border-line text-xs font-semibold text-muted opacity-50"
            >
              {name[0]}
            </button>
          ))}
        </div>
      </div>

      <p className="mt-7 text-center text-sm text-muted">
        Don&apos;t have an account?{" "}
        <Link
          to="/register"
          className="font-semibold text-teal transition-colors duration-200 hover:text-teal-dark"
        >
          Create an account
        </Link>
      </p>
    </AuthLayout>
  );
}
