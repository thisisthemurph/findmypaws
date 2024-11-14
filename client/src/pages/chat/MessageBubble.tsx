import { format } from "date-fns";
import { Message } from "@/pages/chat/hooks/useChat.ts";
import { useEffect, useState } from "react";
import EmojiReactionButton from "@/pages/chat/EmojiReactionButton.tsx";
import EmojiBar, { AllowedEmojis, EmojiKey } from "@/pages/chat/EmojiBar.tsx";

interface MessageBubbleProps {
  message: Message;
  direction: "incoming" | "outgoing";
  emojiBarOpen: boolean;
  onOpenEmojiBar: (messageId: number) => void;
  onCloseEmojiBar: () => void;
  onUpdateEmoji: (messageId: number, emoji: string) => void;
}

export default function MessageBubble({ message, direction, ...props }: MessageBubbleProps) {
  const [emoji, setEmoji] = useState<string>(message.emoji ?? "");
  const outgoing = direction === "outgoing";

  function handleSetEmoji(emojiKey: EmojiKey) {
    if (Object.keys(AllowedEmojis).includes(emojiKey)) {
      setEmoji(AllowedEmojis[emojiKey]);
      props.onUpdateEmoji(message.id, emojiKey);
    }
    props.onCloseEmojiBar();
  }

  function handleClearEmoji() {
    setEmoji("");
    props.onUpdateEmoji(message.id, "");
    props.onCloseEmojiBar();
  }

  useEffect(() => {
    setEmoji(message.emoji ?? "");
  }, [message.emoji]);

  return (
    <div className={`relative group flex items-center w-full ${outgoing ? "flex-row-reverse" : ""}`}>
      <div
        className={`flex flex-col px-4 py-3 text-sm max-w-[80%] w-fit rounded-xl shadow whitespace-pre-line ${outgoing ? "bg-[#7F00FF] text-white" : "bg-white"}`}
      >
        {props.emojiBarOpen ? (
          <EmojiBar
            currentEmoji={emoji}
            handleClearEmoji={handleClearEmoji}
            handleSetEmoji={handleSetEmoji}
          />
        ) : (
          <>
            <p>{message.text}</p>
            <span className={`text-xs ${outgoing ? "text-white text-right" : "text-slate-700"}`}>
              {format(new Date(message.timestamp), "HH:mm")}
            </span>
          </>
        )}
      </div>
      <EmojiReactionButton
        messageId={message.id}
        emoji={emoji}
        outgoing={outgoing}
        emojiBarOpen={props.emojiBarOpen}
        onOpenEmojiBar={props.onOpenEmojiBar}
        onCloseEmojiBar={props.onCloseEmojiBar}
      />
    </div>
  );
}
