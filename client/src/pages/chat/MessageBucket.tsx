import { Message } from "@/components/useChat.ts";
import MessageBubble from "@/pages/chat/MessageBubble";

interface MessageBucketProps {
  name: string;
  messages: Message[];
  currentUserId: string;
}

export default function MessageBucket({ name, messages, currentUserId }: MessageBucketProps) {
  return (
    <>
      <div className="text-center my-2 text-slate-700">{name}</div>
      {messages.map((message) => (
        <MessageBubble
          key={message.timestamp}
          message={message}
          direction={message.senderId === currentUserId ? "outgoing" : "incoming"}
        />
      ))}
    </>
  );
}
