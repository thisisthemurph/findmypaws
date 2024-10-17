import { Outlet } from "react-router-dom";
import Navigation from "../components/Navigation.tsx";
import {NavigationProvider} from "../contexts/NavigationContext.tsx";

function Root() {
  return (
    <>
      <NavigationProvider>
        <Navigation />
      </NavigationProvider>
      <main>
        <Outlet/>
      </main>
    </>
  )
}

export default Root;
