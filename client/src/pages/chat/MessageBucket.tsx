import { Message } from "@/components/useChat.ts";
import MessageBubble from "@/pages/chat/MessageBubble";

interface MessageBucketProps {
  name: string;
  messages: Message[];
  currentUserId: string;
  handleEmojiReact: (messageId: number, emojiKey: string) => void;
}

export default function MessageBucket({
  name,
  messages,
  currentUserId,
  handleEmojiReact,
}: MessageBucketProps) {
  return (
    <>
      <div className="text-center my-2 text-slate-700">{name}</div>
      {messages.map((message) => (
        <MessageBubble
          key={message.id}
          message={message}
          direction={message.senderId === currentUserId ? "outgoing" : "incoming"}
          onUpdateEmoji={(messageId, emoji) => {
            handleEmojiReact(messageId, emoji);
          }}
        />
      ))}
    </>
  );
}
