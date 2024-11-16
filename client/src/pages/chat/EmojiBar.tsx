interface EmojiBarProps {
  currentEmoji: string;
  handleClearEmoji: () => void;
  handleSetEmoji: (emojiKey: EmojiKey) => void;
}

export const AllowedEmojis = {
  "thumbs-up": "👍",
  "thumbs-down": "👎",
  "smiling-face": "😊",
  "laughing-face": "😆",
  "crying-face": "😭",
} as const;

export type EmojiKey = keyof typeof AllowedEmojis;

export default function EmojiBar(props: EmojiBarProps) {
  return (
    <section className={`flex gap-2 px-2 py-2 rounded`}>
      {props.currentEmoji && (
        <button
          onMouseDown={props.handleClearEmoji}
          title="Clear emoji"
          className="text-red-500 hover:scale-125"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="size-4"
          >
            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
          </svg>
        </button>
      )}
      {Object.keys(AllowedEmojis).map((emojiKey) => (
        <button
          key={emojiKey}
          onMouseDown={() => props.handleSetEmoji(emojiKey as EmojiKey)}
          className="hover:scale-125"
        >
          {AllowedEmojis[emojiKey as EmojiKey]}
        </button>
      ))}
    </section>
  );
}
