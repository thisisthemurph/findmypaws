import { useEffect, useState } from "react";
import { z } from "zod";
import useParticipantId from "@/hooks/useParticipantId.ts";

const SendMessageSchema = z.object({
  text: z.string(),
  senderId: z.string(),
});

const MessageSchema = z.object({
  text: z.string(),
  senderId: z.string(),
  timestamp: z.string(),
});

const MessageEventSchema = z.discriminatedUnion("type", [
  z.object({
    type: z.literal("send_message"),
    payload: SendMessageSchema,
  }),
  z.object({
    type: z.literal("new_message"),
    payload: MessageSchema,
  }),
]);

export type SendMessage = z.infer<typeof SendMessageSchema>;
export type Message = z.infer<typeof MessageSchema>;
export type MessageEvent = z.infer<typeof MessageEventSchema>;

type GroupedMessages = {
  key: string;
  messages: Message[];
}[];

const groupMessagesByTime = (messages: z.infer<typeof MessageSchema>[]) => {
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());

  const grouped: Record<string, Message[]> = {};

  messages
    .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
    .forEach((message) => {
      const messageDate = new Date(message.timestamp);
      const messageDayStart = new Date(
        messageDate.getFullYear(),
        messageDate.getMonth(),
        messageDate.getDate()
      );

      const daysAgo = Math.floor((today.getTime() - messageDayStart.getTime()) / (1000 * 60 * 60 * 24));

      let key: string;
      if (daysAgo === 0) {
        key = "Today";
      } else if (daysAgo < 7) {
        key = messageDate.toLocaleDateString("en-US", { weekday: "long" });
      } else {
        key = messageDate.toLocaleDateString("en-US", {
          day: "numeric",
          month: "long",
          year: "numeric",
        });
      }

      if (!grouped[key]) {
        grouped[key] = [];
      }
      grouped[key].push(message);
    });

  // Convert the object to an array and sort by descending date
  const sortedGroups: GroupedMessages = Object.keys(grouped)
    .map((key) => ({
      key,
      messages: grouped[key].sort(
        (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
      ),
    }))
    .sort((a, b) => {
      if (a.key === "Today") return 1;
      if (b.key === "Today") return -1;
      return new Date(b.messages[0].timestamp).getTime() - new Date(a.messages[0].timestamp).getTime();
    });

  return sortedGroups;
};

export default function useChat(roomIdentifier: string) {
  const [webSocket, setWebSocket] = useState<WebSocket | undefined>();
  const [messages, setMessages] = useState<Message[]>([]);
  const participantId = useParticipantId();

  const bucketedMessages = groupMessagesByTime(messages);

  useEffect(() => {
    if (!participantId) return;

    const webSocketUrl = `ws://localhost:42096/room?r=${roomIdentifier}&pid=${participantId}`;
    const socket = new WebSocket(webSocketUrl);
    socket.onopen = (event) => {
      // Reset the messages to prevent loading the same ones again.
      setMessages([]);
      console.log("WebSocket opened", event);
    };

    // New message received from the WebSocket.
    socket.onmessage = (event) => {
      const eventData = JSON.parse(event.data);
      const result = MessageEventSchema.safeParse(eventData);
      if (!result.success) {
        console.error("Invalid MessageEvent format", result.error);
        return;
      }

      const receivedEvent = result.data;
      switch (receivedEvent.type) {
        case "new_message":
          setMessages((prev) => [...prev, receivedEvent.payload]);
          break;
        case "send_message":
          console.warn("SendMessage event received, but no action taken.");
          break;
        default:
          console.error("Unsupported event type", receivedEvent);
      }
    };

    socket.onclose = (event) => {
      console.log("Closing ws", event);
    };

    socket.onerror = (event) => {
      console.error("error on ws", event);
    };

    setWebSocket(socket);
    return () => socket.close();
  }, [roomIdentifier, participantId]);

  const sendMessage = (text: string) => {
    if (!webSocket) throw new Error("WebSocket not available");
    if (!participantId) throw new Error("Participant ID is undefined");

    const event: MessageEvent = {
      type: "send_message",
      payload: {
        text: text,
        senderId: participantId,
      },
    };

    webSocket.send(JSON.stringify(event));
  };

  return {
    roomId: roomIdentifier,
    participantId,
    messages,
    bucketedMessages,
    sendMessage,
  };
}
