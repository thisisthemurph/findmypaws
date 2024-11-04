import { UserProfile } from "@clerk/clerk-react";

export default function ClerkProfilePage() {
  return (
    <section className="flex justify-center my-12">
      <UserProfile path="/user-profile" />
    </section>
  );
}
