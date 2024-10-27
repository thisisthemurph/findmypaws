import { UserProfile } from "@clerk/clerk-react";

export default function ClerkProfilePage() {
  return <UserProfile path="/user-profile" />;
}
