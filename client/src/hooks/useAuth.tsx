import { useContext } from "react";
import { AuthContext } from "../contexts/AuthContext.tsx";
import { AuthContextProps } from "../contexts/authContextProps.ts";

export const useAuth = (): AuthContextProps => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within the AuthProvider");
  }
  return context;
}