import { Message } from "@/pages/chat/hooks/useChat.ts";
import MessageBubble from "@/pages/chat/MessageBubble";

interface MessageBucketProps {
  name: string;
  messages: Message[];
  currentUserId: string;
  openEmojiBarMessageId: number | undefined;
  onOpenEmojiBar: (messageId: number) => void;
  onCloseEmojiBar: () => void;
  handleEmojiReact: (messageId: number, emojiKey: string) => void;
}

export default function MessageBucket({
  name,
  messages,
  currentUserId,
  openEmojiBarMessageId,
  handleEmojiReact,
  ...props
}: MessageBucketProps) {
  return (
    <>
      <div className="text-center text-slate-700">{name}</div>
      {messages.map((message) => (
        <MessageBubble
          key={message.id}
          message={message}
          direction={message.senderId === currentUserId ? "outgoing" : "incoming"}
          emojiBarOpen={openEmojiBarMessageId === message.id}
          onOpenEmojiBar={props.onOpenEmojiBar}
          onCloseEmojiBar={props.onCloseEmojiBar}
          onUpdateEmoji={(messageId, emoji) => {
            handleEmojiReact(messageId, emoji);
          }}
        />
      ))}
    </>
  );
}
