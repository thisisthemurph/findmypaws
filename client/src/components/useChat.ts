import { useEffect, useState } from "react";
import { z } from "zod";
import useParticipantId from "@/hooks/useParticipantId.ts";
import { useApi } from "@/hooks/useApi.ts";
import { Conversation } from "@/api/types.ts";

const SendMessageSchema = z.object({
  text: z.string(),
  senderId: z.string(),
});

const MessageSchema = z.object({
  id: z.number(),
  text: z.string(),
  emoji: z.string().nullable(),
  senderId: z.string(),
  timestamp: z.string(),
});

const EmojiReactSchema = z.object({
  conversationId: z.number(),
  messageId: z.number(),
  emojiKey: z.string(),
});

const NewEmojiReactSchema = z.object({
  messageId: z.number(),
  emoji: z.string().nullable(),
});

const MessageEventSchema = z.discriminatedUnion("type", [
  // Event for sending a new message
  z.object({
    type: z.literal("send_message"),
    payload: SendMessageSchema,
  }),
  // Event for incoming messages
  z.object({
    type: z.literal("new_message"),
    payload: MessageSchema,
  }),
  // Event for sending an emoji reaction
  z.object({
    type: z.literal("emoji_react"),
    payload: EmojiReactSchema,
  }),
  // Event for incoming emoji reactions
  z.object({
    type: z.literal("new_emoji_react"),
    payload: NewEmojiReactSchema,
  }),
]);

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

  const customOrder = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday", "Today"];

  // Convert the object to an array and sort by descending date
  const sortedGroups: GroupedMessages = Object.keys(grouped)
    .map((key) => ({
      key,
      messages: grouped[key].sort(
        (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
      ),
    }))
    .sort((a, b) => {
      const aIndex = customOrder.indexOf(a.key);
      const bIndex = customOrder.indexOf(b.key);

      if (aIndex !== -1 && bIndex !== -1) {
        return aIndex - bIndex;
      }
      return new Date(b.messages[0].timestamp).getTime() - new Date(a.messages[0].timestamp).getTime();
    });

  return sortedGroups;
};

export default function useChat(roomIdentifier: string) {
  const [webSocket, setWebSocket] = useState<WebSocket | undefined>();
  const [conversation, setConversation] = useState<Conversation | undefined>();
  const [messages, setMessages] = useState<Message[]>([]);
  const participantId = useParticipantId();
  const api = useApi();

  const [isConversationDetailsLoaded, setIsConversationDetailsLoaded] = useState(false);
  const [isWebSocketLoaded, setIsWebSocketLoaded] = useState(false);

  const bucketedMessages = groupMessagesByTime(messages);

  useEffect(() => {
    const getChatTitle = async (identifier: string) => {
      return await api<Conversation>(`/conversations/${identifier}`);
    };

    getChatTitle(roomIdentifier)
      .then((conversation) => {
        setConversation(conversation);
        setIsConversationDetailsLoaded(true);
      })
      .catch(() => {
        console.error("failed to get chat name");
        setIsConversationDetailsLoaded(false);
      });
  }, [roomIdentifier]);

  useEffect(() => {
    if (!participantId) return;

    const webSocketUrl = `ws://localhost:42096/room?r=${roomIdentifier}&pid=${participantId}`;
    const socket = new WebSocket(webSocketUrl);
    socket.onopen = () => {
      // Reset the messages to prevent loading the same ones again.
      setMessages([]);
      setIsWebSocketLoaded(true);
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
          setMessages((previousMessages) => [...previousMessages, receivedEvent.payload]);
          break;
        case "send_message":
          console.warn("SendMessage event received, but no action taken.");
          break;
        case "new_emoji_react":
          setMessages((previousMessages) =>
            previousMessages.map((message) =>
              message.id === receivedEvent.payload.messageId
                ? { ...message, emoji: receivedEvent.payload.emoji }
                : message
            )
          );
          break;
        default:
          console.error("Unsupported event type", receivedEvent);
      }
    };

    socket.onclose = (event) => {
      console.warn("Closing WebSocket", event);
      setIsWebSocketLoaded(false);
    };

    socket.onerror = (event) => {
      console.error("error on WebSocket", event);
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

  const react = (messageId: number, emojiKey: string) => {
    if (!conversation) return;
    if (!webSocket) throw new Error("WebSocket not available");
    if (!participantId) throw new Error("Participant ID is undefined");

    console.log({ conversationId: conversation.id, messageId, emojiKey });

    const event: MessageEvent = {
      type: "emoji_react",
      payload: {
        conversationId: conversation.id,
        messageId: messageId,
        emojiKey: emojiKey,
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
    react,
    conversation,
    isLoaded: isWebSocketLoaded && isConversationDetailsLoaded,
  };
}
