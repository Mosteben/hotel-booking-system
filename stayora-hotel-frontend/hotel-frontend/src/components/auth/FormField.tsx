import { forwardRef, type InputHTMLAttributes, type SelectHTMLAttributes } from "react";

interface FormFieldProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
  optional?: boolean;
}

const baseFieldClass =
  "w-full rounded-[13px] border bg-white px-4 py-3.5 text-sm text-ink outline-none transition-colors duration-200 placeholder:text-muted";

const borderClass = (hasError?: string) =>
  hasError
    ? "border-red-300 focus:border-red-400 focus:ring-4 focus:ring-red-50"
    : "border-line focus:border-teal focus:ring-4 focus:ring-teal/10";

export const FormField = forwardRef<HTMLInputElement, FormFieldProps>(
  function FormField({ label, error, optional, className = "", ...props }, ref) {
    return (
      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium text-ink">
          {label}
          {optional && (
            <span className="ml-1.5 text-xs font-normal text-muted">
              (optional)
            </span>
          )}
        </label>
        <input
          ref={ref}
          className={`${baseFieldClass} ${borderClass(error)} ${className}`}
          {...props}
        />
        {error && (
          <p className="text-xs font-medium text-red-500" role="alert">
            {error}
          </p>
        )}
      </div>
    );
  }
);

interface FormSelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
  label: string;
  error?: string;
  optional?: boolean;
  children: React.ReactNode;
}

export const FormSelect = forwardRef<HTMLSelectElement, FormSelectProps>(
  function FormSelect(
    { label, error, optional, className = "", children, ...props },
    ref
  ) {
    return (
      <div className="flex flex-col gap-1.5">
        <label className="text-sm font-medium text-ink">
          {label}
          {optional && (
            <span className="ml-1.5 text-xs font-normal text-muted">
              (optional)
            </span>
          )}
        </label>
        <select
          ref={ref}
          className={`${baseFieldClass} ${borderClass(error)} cursor-pointer ${className}`}
          {...props}
        >
          {children}
        </select>
        {error && (
          <p className="text-xs font-medium text-red-500" role="alert">
            {error}
          </p>
        )}
      </div>
    );
  }
);
