import { Loader2 } from "lucide-react";
import type { ButtonHTMLAttributes } from "react";

interface AuthButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  isSubmitting?: boolean;
}

export function AuthButton({
  isSubmitting,
  children,
  className = "",
  disabled,
  ...props
}: AuthButtonProps) {
  return (
    <button
      type="submit"
      disabled={disabled ?? isSubmitting}
      className={`mt-1 flex h-[52px] w-full items-center justify-center gap-2 rounded-[13px] bg-teal text-sm font-semibold text-white transition-all duration-150 hover:bg-teal-dark active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-60 disabled:active:scale-100 ${className}`}
      {...props}
    >
      {isSubmitting && <Loader2 size={16} className="animate-spin" />}
      {children}
    </button>
  );
}
