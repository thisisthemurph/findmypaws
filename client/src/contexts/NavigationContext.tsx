import {createContext, ReactNode, useState} from "react";

const NavigationContext = createContext<NavigationContextProps | undefined>(undefined);

export type NavigationContextProps = [boolean, (open: boolean) => void];

const NavigationProvider = ({ children } : {children: ReactNode}) => {
  const [open, setOpen] = useState(false);

  return (
    <NavigationContext.Provider value={[open, (open: boolean) => setOpen(open)]}>
      {children}
    </NavigationContext.Provider>
  )
}

export { NavigationContext, NavigationProvider };
