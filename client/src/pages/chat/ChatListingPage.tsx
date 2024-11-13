import { Link, useParams } from "react-router-dom";
import { useApi } from "@/hooks/useApi.ts";
import { useQuery } from "@tanstack/react-query";
import { Conversation } from "@/api/types.ts";
import { differenceInDays, format, formatDistanceToNow } from "date-fns";
import { PageHeading } from "@/components/PageHeading.tsx";

export default function ChatListingPage() {
  const api = useApi();
  const { identifier } = useParams();

  const { data: chats, isLoading } = useQuery({
    queryKey: ["conversations", identifier],
    queryFn: async () => await api<Conversation[]>(`/conversations`),
  });

  function formatTimeAgo(date: Date) {
    if (differenceInDays(new Date(), date) > 3) {
      return format(date, "PPP");
    }
    return formatDistanceToNow(date, { addSuffix: true });
  }

  if (isLoading || !chats) {
    return <p>Loading</p>;
  }

  return (
    <>
      <PageHeading
        heading="Your chats"
        subheading={
          chats.length > 0
            ? undefined
            : "You don't have any chats at the moment. If one of your pets is lost and someone finds them, a chat will be started and shown here."
        }
      />

      <section>
        {chats
          .sort(
            (a, b) =>
              new Date(a.lastMessageAt ?? a.createdAt).getTime() -
              new Date(b.lastMessageAt ?? b.createdAt).getTime()
          )
          .map((conversation) => (
            <Link
              to={`/conversations/${conversation.identifier}`}
              key={conversation.id}
              className="block p-4 bg-slate-50 border-b"
            >
              <p>
                A chat with an {conversation.otherParticipant.name} about{" "}
                <span className="font-semibold">{conversation.pet.name}</span>.
              </p>
              <p
                className="text-right text-sm text-slate-800 mt-2"
                title={format(new Date(conversation.createdAt), "PPP")}
              >
                {formatTimeAgo(new Date(conversation.lastMessageAt ?? conversation.createdAt))}
              </p>
            </Link>
          ))}
      </section>
    </>
  );
}
