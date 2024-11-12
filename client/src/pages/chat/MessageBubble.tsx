import { format } from "date-fns";
import { Message } from "@/components/useChat.ts";

interface MessageBubbleProps {
  message: Message;
  direction: "incoming" | "outgoing";
}

export default function MessageBubble({ message, direction }: MessageBubbleProps) {
  const outgoing = direction == "outgoing";

  return (
    <p
      className={`flex flex-col px-4 py-3 text-sm max-w-[80%] w-fit rounded-xl shadow whitespace-pre-line ${outgoing ? "self-end bg-[#7F00FF] text-white" : "bg-white"}`}
    >
      <span>{message.text}</span>
      <span className={`text-xs ${outgoing ? "text-white text-right" : "text-slate-700"}`}>
        {format(new Date(message.timestamp), "HH:mm")}
      </span>
    </p>
  );
}
