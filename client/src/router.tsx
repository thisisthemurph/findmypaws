import { createBrowserRouter } from "react-router-dom";

import Root from "./pages/Root.tsx";
import PetProfile from "./pages/PetProfile.tsx";
import Home from "./pages/Home.tsx";
import LogIn from "./pages/auth/LogIn.tsx";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      {
        path: "/",
        element: <Home />
      },
      {
        path: "/login",
        element: <LogIn />
      },
      {
        path: "/profile/:id",
        element: <PetProfile />,
      }
    ],
  }
]);
