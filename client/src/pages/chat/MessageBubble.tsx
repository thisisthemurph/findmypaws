import { format } from "date-fns";
import { Message } from "@/components/useChat.ts";
import { Button } from "@/components/ui/button.tsx";
import { useEffect, useState } from "react";

interface MessageBubbleProps {
  message: Message;
  direction: "incoming" | "outgoing";
  onUpdateEmoji: (messageId: number, emoji: string) => void;
}

const allowedEmojis = {
  "thumbs-up": "👍",
  "thumbs-down": "👎",
  "smiling-face": "😊",
  "laughing-face": "😆",
  "crying-face": "😭",
} as const;

type EmojiKey = keyof typeof allowedEmojis;

export default function MessageBubble({ message, direction, onUpdateEmoji }: MessageBubbleProps) {
  const [showEmojiBar, setShowEmojiBar] = useState(false);
  const [emoji, setEmoji] = useState<string>(message.emoji ?? "");
  const outgoing = direction === "outgoing";

  function handleSetEmoji(emojiKey: EmojiKey) {
    if (Object.keys(allowedEmojis).includes(emojiKey)) {
      setEmoji(allowedEmojis[emojiKey]);
      onUpdateEmoji(message.id, emojiKey);
    }
    setShowEmojiBar(false);
  }

  useEffect(() => {
    setEmoji(message.emoji ?? "");
  }, [message.emoji]);

  return (
    <div className={`relative group flex items-center w-full ${outgoing ? "flex-row-reverse" : ""}`}>
      <div
        className={`flex flex-col px-4 py-3 text-sm max-w-[80%] w-fit rounded-xl shadow whitespace-pre-line ${outgoing ? "bg-[#7F00FF] text-white" : "bg-white"}`}
      >
        {showEmojiBar ? (
          <section className={`flex gap-2 px-2 py-2 bg-slate-500 rounded`}>
            {Object.keys(allowedEmojis).map((emojiKey) => (
              <button key={emojiKey} onMouseDown={() => handleSetEmoji(emojiKey as EmojiKey)}>
                {allowedEmojis[emojiKey as EmojiKey]}
              </button>
            ))}
          </section>
        ) : (
          <>
            <p>{message.text}</p>
            <span className={`text-xs ${outgoing ? "text-white text-right" : "text-slate-700"}`}>
              {format(new Date(message.timestamp), "HH:mm")}
            </span>
          </>
        )}
      </div>
      <Button
        variant="ghost"
        size="icon"
        className={`${emoji ? "flex" : "hidden"} text-slate-600 group-hover:flex hover:bg-transparent hover:text-black hover:scale-125`}
        onMouseDown={() => setShowEmojiBar(!showEmojiBar)}
      >
        {emoji ? (
          <span>{emoji}</span>
        ) : (
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15.182 15.182a4.5 4.5 0 0 1-6.364 0M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z"
            />
          </svg>
        )}
      </Button>
    </div>
  );
}
