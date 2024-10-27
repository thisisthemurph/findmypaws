import { useEffect } from "react";
import { useAuth } from "@clerk/clerk-react";
import { Outlet, useLocation, useNavigate } from "react-router-dom";

export default function SignedInLayout() {
  const { userId, isLoaded } = useAuth();
  const navigate = useNavigate();
  const { pathname } = useLocation();

  useEffect(() => {
    if (isLoaded && !userId) {
      if (!pathname) navigate("/sign-in");
      navigate(`sign-in?redirect_url=${pathname}`);
    }
  }, [isLoaded]);

  if (!isLoaded) return "Loading...";

  return <Outlet />;
}
