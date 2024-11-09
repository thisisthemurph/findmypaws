import { useState } from "react";

const AnonymousUserIdKey = "anonymousUserId";

type UseAnonymousUserId = [string | null, (anonymousUserId: string | null) => void];

export default function useAnonymousUser(): UseAnonymousUserId {
  const [anonymousUserId, setAnonymousUserId] = useState<string | null>(() =>
    localStorage.getItem(AnonymousUserIdKey)
  );

  const update = (anonymousUserId: string | null) => {
    if (!anonymousUserId) {
      localStorage.removeItem(AnonymousUserIdKey);
      setAnonymousUserId(null);
      return;
    }

    localStorage.setItem(AnonymousUserIdKey, anonymousUserId);
    setAnonymousUserId(anonymousUserId);
  };

  return [anonymousUserId, update];
}
