import { useEffect, useState } from "react";

type DraftsCollection = {
  [key: string]: string;
};

type UseDraftMessage = [
  (identifier: string) => string,
  (identifier: string, message: string) => void,
  (identifier: string) => void,
];

export default function useDraftMessage(): UseDraftMessage {
  const [drafts, setDrafts] = useState<DraftsCollection>(() => {
    const draftsJson = localStorage.getItem("drafts");
    return draftsJson ? (JSON.parse(draftsJson) as DraftsCollection) : {};
  });

  useEffect(() => {
    localStorage.setItem("drafts", JSON.stringify(drafts));
  }, [drafts]);

  function getDraft(identifier: string): string {
    return drafts[identifier] ?? "";
  }

  function setDraft(identifier: string, message: string): void {
    setDrafts((prev) => ({ ...prev, [identifier]: message }));
  }

  function deleteDraft(identifier: string): void {
    setDrafts((prev) => {
      const updatedDrafts = { ...prev };
      delete updatedDrafts[identifier];
      return updatedDrafts;
    });
  }

  return [getDraft, setDraft, deleteDraft];
}
