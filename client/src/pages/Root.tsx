import { Outlet } from "react-router-dom";
import Navigation from "../components/Navigation.tsx";
import { NavigationProvider } from "../contexts/NavigationContext.tsx";
import { Toaster } from "@/components/ui/toaster.tsx";

function Root() {
  return (
    <>
      <NavigationProvider>
        <Navigation />
      </NavigationProvider>
      <main>
        <Outlet />
      </main>
      <Toaster />
    </>
  );
}

export default Root;
