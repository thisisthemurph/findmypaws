import { createBrowserRouter } from "react-router-dom";

import HomePage from "@/pages/home";
import RootLayout from "@/layouts/root-layout.tsx";
import ClerkSignInPage from "@/pages/clerk/ClerkSignInPage.tsx";
import ClerkProfilePage from "@/pages/clerk/ClerkProfilePage.tsx";
import DashboardPage from "@/pages/dashboard";
import SignedInLayout from "@/layouts/signed-in-layout.tsx";
import PetPage from "@/pages/pet";
import ConversationPage from "@/pages/messenger/ConversationPage.tsx";
import ConversationListingPage from "@/pages/messenger/ConversationListingPage.tsx";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    children: [
      {
        path: "/",
        element: <HomePage />,
      },
      {
        path: "/sign-in/*",
        element: <ClerkSignInPage />,
      },
      {
        path: "/user-profile",
        element: <ClerkProfilePage />,
      },
      {
        path: "dashboard",
        element: <SignedInLayout />,
        children: [
          {
            path: "/dashboard",
            element: <DashboardPage />,
          },
        ],
      },
      {
        path: "/pet",
        children: [
          {
            path: "/pet/:id",
            element: <PetPage />,
          },
        ],
      },
      {
        path: "/conversations",
        element: <ConversationListingPage />,
      },
      {
        path: "/conversations/:conversationIdentifier",
        element: <ConversationPage />,
      },
    ],
  },
]);
