import { Button } from "@/components/ui/button.tsx";

interface EmojiReactionButtonProps {
  messageId: number;
  emoji: string | undefined;
  outgoing: boolean;
  emojiBarOpen: boolean;
  onOpenEmojiBar: (messageId: number) => void;
  onCloseEmojiBar: () => void;
}

export default function EmojiReactionButton({
  messageId,
  emoji,
  outgoing,
  ...props
}: EmojiReactionButtonProps) {
  if (outgoing && !emoji) {
    return null;
  }

  if (outgoing && emoji) {
    return <span className="p-2">{emoji}</span>;
  }

  return (
    <Button
      variant="ghost"
      size="icon"
      className={`${emoji ? "flex" : "hidden"} text-slate-600 group-hover:flex hover:bg-transparent hover:text-black hover:scale-125`}
      onMouseDown={() => {
        if (props.emojiBarOpen || outgoing) {
          props.onCloseEmojiBar();
        }
        if (!props.emojiBarOpen && !outgoing) {
          props.onOpenEmojiBar(messageId);
        }
      }}
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
  );
}
