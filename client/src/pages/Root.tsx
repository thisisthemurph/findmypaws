import { Outlet } from "react-router-dom";
import Navigation from "../components/Navigation.tsx";

function Root() {
  return (
    <>
      <Navigation />
      <main>
        <Outlet/>
      </main>
    </>
  )
}

export default Root;
