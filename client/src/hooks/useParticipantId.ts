import { useAuth } from "@clerk/clerk-react";
import useAnonymousUser from "@/hooks/useAnonymousUser.ts";

export default function useParticipantId() {
  const { userId, isLoaded } = useAuth();
  const [anonymousUserId] = useAnonymousUser();

  if (!isLoaded) return undefined;
  return userId ?? anonymousUserId ?? undefined;
}
