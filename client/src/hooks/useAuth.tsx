import { useContext } from "react";
import { AuthContext, AuthContextProps } from "../contexts/AuthContext.tsx";

export const useAuth = (): AuthContextProps => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within the AuthProvider");
  }
  return context;
}