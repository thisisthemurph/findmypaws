import { useAuth } from "@/hooks/useAuth.tsx";
import { Navigate, Outlet } from "react-router-dom";

function ProtectedRoute() {
  const auth = useAuth();
  if (!auth.isAuthenticated()) {
    return <Navigate to="/login" />;
  }
  return <Outlet />;
}

export default ProtectedRoute;
