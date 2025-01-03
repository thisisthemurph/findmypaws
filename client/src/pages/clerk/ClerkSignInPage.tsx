import { SignIn, useUser } from "@clerk/clerk-react";
import { useNavigate } from "react-router-dom";

export default function ClerkSignInPage() {
  const navigate = useNavigate();
  const { user } = useUser();

  if (!user) {
    return (
      <section className="flex justify-center my-12">
        <SignIn path="/sign-in" />
      </section>
    );
  }

  navigate("/");
}
