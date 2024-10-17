import {useContext} from "react";
import {NavigationContext, NavigationContextProps} from "../contexts/NavigationContext.tsx";

export const useNavigation = (): NavigationContextProps => {
  const context = useContext(NavigationContext);
  if (!context) {
    throw new Error("useNavigation must be used within NavigationProvider");
  }
  return context;
}
