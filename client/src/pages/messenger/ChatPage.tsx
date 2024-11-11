import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form.tsx";
import useChat from "@/components/useChat.ts";
import MessageBucket from "@/pages/messenger/MessageBucket.tsx";
import { useParams } from "react-router-dom";
import ChatSubmitButton from "@/pages/messenger/ChatSubmitButton.tsx";
import { useEffect, useRef } from "react";

const formSchema = z.object({
  text: z.string().min(1, "Enter a message"),
});

type FormInputs = z.infer<typeof formSchema>;

export default function ChatPage() {
  const { conversationIdentifier } = useParams();
  const {
    title,
    participantId,
    bucketedMessages,
    sendMessage,
    isLoaded: isChatLoaded,
  } = useChat(conversationIdentifier!);

  const canSendMessage = isChatLoaded;
  const chatSectionRef = useRef<HTMLDivElement>(null);

  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      text: "",
    },
  });

  function onSubmit(data: FormInputs) {
    sendMessage(data.text);
    form.reset();
  }

  useEffect(() => {
    if (isChatLoaded && chatSectionRef.current) {
      chatSectionRef.current.scrollTop = chatSectionRef.current.scrollHeight;
    }
  }, [isChatLoaded, bucketedMessages]);

  return (
    <div className="flex flex-col h-[calc(100vh-5rem)] bg-slate-50">
      <section className="flex justify-center py-4">
        <p className="font-semibold">{isChatLoaded ? title : "Loading"}</p>
      </section>

      <section ref={chatSectionRef} className="flex-grow p-4 overflow-y-auto">
        <div className="flex flex-col gap-1">
          {isChatLoaded &&
            participantId &&
            bucketedMessages.map((bucket) => (
              <MessageBucket
                key={bucket.key}
                name={bucket.key}
                messages={bucket.messages}
                currentUserId={participantId}
              />
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
                  <textarea
                    placeholder={isChatLoaded ? "Write a message..." : "Loading..."}
                    className={`p-4 rounded-lg w-full resize-none shadow ${!canSendMessage && "disabled:bg-slate-100 shadow-none"}`}
                    disabled={!canSendMessage}
                    {...field}
                    onKeyDown={(e) => {
                      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
                        e.preventDefault();
                        form.handleSubmit(onSubmit)();
                      }
                    }}
                  ></textarea>
                </FormControl>
              </FormItem>
            )}
          />
          <ChatSubmitButton show={form.getValues("text") !== ""} />
        </form>
      </Form>
    </div>
  );
}
