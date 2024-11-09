import { Link, useParams } from "react-router-dom";
import { useApi } from "@/hooks/useApi.ts";
import { useQuery } from "@tanstack/react-query";
import { Conversation } from "@/api/types.ts";

export default function ConversationListingPage() {
  const api = useApi();
  const { identifier } = useParams();

  const { data, isLoading } = useQuery({
    queryKey: ["conversations", identifier],
    queryFn: async () => await api<Conversation[]>(`/conversations`),
  });

  if (isLoading || !data) {
    return <p>Loading</p>;
  }

  return (
    <>
      <pre>{JSON.stringify(data, null, 2)}</pre>
      <section>
        {data.map((conversation) => (
          <Link
            to={`/conversations/${conversation.id}`}
            key={conversation.id}
            className="block p-4 bg-slate-50 border-b"
          >
            <p>
              A chat with an anonymous user about <strong>{conversation.pet.name}</strong>.
            </p>
            <p className="text-right text-sm text-slate-800 mt-2">{conversation.lastMessageAt}</p>
          </Link>
        ))}
      </section>
    </>
  );
}
