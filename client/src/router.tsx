import { createBrowserRouter } from "react-router-dom";

import Root from "@/pages/Root.tsx";
import Home from "@/pages/Home.tsx";
import LogInPage from "@/pages/auth/LogInPage.tsx";
import SignUpPage from "@/pages/auth/SignUpPage.tsx";
import ProfilePage from "@/pages/profile";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "/login",
        element: <LogInPage />,
      },
      {
        path: "/signup",
        element: <SignUpPage />,
      },
      {
        path: "/profile",
        element: <ProfilePage />,
      },
    ],
  },
]);
