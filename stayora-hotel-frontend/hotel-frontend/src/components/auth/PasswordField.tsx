import { useId, useState, type InputHTMLAttributes } from "react";
import { Eye, EyeOff } from "lucide-react";
import { FormField } from "@/components/auth/FormField";

interface PasswordFieldProps
  extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  label: string;
  error?: string;
}

export function PasswordField({
  label,
  error,
  className = "",
  ...props
}: PasswordFieldProps) {
  const [visible, setVisible] = useState(false);
  const toggleId = useId();

  return (
    <div className="relative">
      <FormField
        label={label}
        type={visible ? "text" : "password"}
        error={error}
        className={`pr-11 ${className}`}
        aria-describedby={toggleId}
        {...props}
      />
      <button
        type="button"
        id={toggleId}
        onClick={() => setVisible((v) => !v)}
        className="absolute right-3.5 top-[38px] cursor-pointer text-muted transition-colors hover:text-ink"
        aria-label={visible ? "Hide password" : "Show password"}
      >
        {visible ? <EyeOff size={17} /> : <Eye size={17} />}
      </button>
    </div>
  );
}
