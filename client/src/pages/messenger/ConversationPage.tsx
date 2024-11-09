import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
// import { useParams } from "react-router-dom";
// import { ConversationWithMessages, Message } from "@/api/types.ts";
// import { useToast } from "@/hooks/use-toast.ts";
import { useAuth } from "@clerk/clerk-react";
import { useEffect, useState } from "react";
import useAnonymousUser from "@/hooks/useAnonymousUser.ts";

const formSchema = z.object({
  text: z.string().min(1, "Enter a message"),
});

type FormInputs = z.infer<typeof formSchema>;

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

type SendMessage = z.infer<typeof SendMessageSchema>;
type Message = z.infer<typeof MessageSchema>;
type MessageEvent = z.infer<typeof MessageEventSchema>;

export default function ConversationPage() {
  // const { conversationId } = useParams();
  // const { toast } = useToast();
  const { userId, isLoaded: isUserLoaded } = useAuth();
  const [ws, setWS] = useState<WebSocket | undefined>();
  const [messages, setMessages] = useState<Message[]>([]);
  const [anonymousUserId] = useAnonymousUser();

  useEffect(() => {
    console.log("setting web socket for user", userId);
    // if (!userId) return;
    if (!isUserLoaded) return;
    const participantId = userId ?? anonymousUserId;
    if (!participantId) return;

    const webSocketUrl = `ws://localhost:42096/room?r=49b6d8d8-816c-4628-9238-fba78ab18c90&pid=${participantId}`;
    console.log(webSocketUrl);
    const webSocket = new WebSocket(webSocketUrl);
    webSocket.onopen = (event) => {
      console.log("WebSocket opened", event);
    };

    // New message received from the ws.
    webSocket.onmessage = (event) => {
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

    webSocket.onclose = (event) => {
      console.log("Closing ws", event);
    };

    webSocket.onerror = (event) => {
      console.error("error on ws", event);
    };

    setWS(webSocket);

    return () => {
      webSocket.close();
    };
  }, [userId, isUserLoaded, anonymousUserId]);

  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      text: "",
    },
  });

  function onSubmit(data: FormInputs) {
    if (!ws) {
      console.error("No web socket");
      return;
    }

    const participantId = userId ?? anonymousUserId;
    if (!participantId) {
      console.error("no user id set");
      return;
    }

    const event: MessageEvent = {
      type: "send_message",
      payload: {
        text: data.text,
        senderId: participantId,
      },
    };

    ws.send(JSON.stringify(event));
    form.reset();
  }

  if (!isUserLoaded) {
    return <p>Loading user</p>;
  }

  const effectiveUserId = userId ?? anonymousUserId;
  if (!effectiveUserId) {
    return <p>Cannot determine the effective user ID</p>;
  }

  return (
    <div className="flex flex-col h-[calc(100vh-5rem)] bg-slate-50">
      <section className="flex justify-center py-4">
        <p className="font-semibold">Mike Murphy</p>
      </section>

      <section className="flex-grow p-4 overflow-y-auto">
        <div className="flex flex-col gap-1">
          {messages.map((m) => (
            <MessageBubble key={m.timestamp} message={m} currentUserId={effectiveUserId} />
          ))}
        </div>
      </section>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="relative p-4">
          <FormField
            control={form.control}
            name="text"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input placeholder="message" className="w-full rounded-full" {...field} />
                </FormControl>
              </FormItem>
            )}
          />
          <Button type="submit" variant="ghost" size="icon" className="absolute top-4 right-4">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="w-10 h-10"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5"
              />
            </svg>
          </Button>
        </form>
      </Form>
    </div>
  );
}

function MessageBubble({ message, currentUserId }: { message: SendMessage; currentUserId: string }) {
  const { text, senderId } = message;
  const outgoing = senderId === currentUserId;

  return (
    <p
      className={`px-4 py-3 text-sm max-w-[80%] w-fit ${outgoing && "self-end"} ${outgoing ? "bg-[#7F00FF] text-white" : "bg-white"} rounded-full shadow`}
    >
      {text}
    </p>
  );
}
