import type { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import type { UserRole } from "@/types/auth";

export function ProtectedRoute({
  children,
  roles,
}: {
  children: ReactNode;
  // When set, only these roles may view the route - anyone else is sent
  // home instead of to /login, since they ARE authenticated, just not
  // authorized for this specific page.
  roles?: UserRole[];
}) {
  const { isAuthenticated, isLoading, user } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-bg">
        <Loader2 size={28} className="animate-spin text-teal" />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (roles && (!user || !roles.includes(user.role))) {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}
