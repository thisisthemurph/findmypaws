import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form.tsx";
import useChat from "@/pages/chat/hooks/useChat.ts";
import MessageBucket from "@/pages/chat/MessageBucket.tsx";
import { useParams } from "react-router-dom";
import ChatSubmitButton from "@/pages/chat/ChatSubmitButton.tsx";
import { useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button.tsx";
import useDraftMessage from "@/pages/chat/hooks/useDraftMessage.ts";

const formSchema = z.object({
  text: z.string().min(1, "Enter a message"),
});

type FormInputs = z.infer<typeof formSchema>;

export default function ChatPage() {
  const { conversationIdentifier } = useParams();
  const [useLargeInput, setUseLargeInput] = useState(false);
  const [getDraft, setDraft, deleteDraft] = useDraftMessage();
  const [openEmojiBarMessageId, setOpenEmojiBarMessageId] = useState<number | undefined>();

  const {
    conversation,
    participantId,
    bucketedMessages,
    sendMessage,
    emojiReact,
    messageCount,
    isLoaded: isChatLoaded,
    otherParticipantIsTyping,
    handleTypingDetection,
  } = useChat(conversationIdentifier!);

  const canSendMessage = isChatLoaded;
  const chatSectionRef = useRef<HTMLDivElement>(null);

  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      text: getDraft(conversationIdentifier!),
    },
  });

  const currentMessage = form.watch("text");

  function onSubmit(data: FormInputs) {
    sendMessage(data.text);
    if (conversationIdentifier) {
      deleteDraft(conversationIdentifier);
    }
    form.setValue("text", "");
    form.setFocus("text");
  }

  useEffect(() => {
    if (!conversationIdentifier) return;
    setDraft(conversationIdentifier, currentMessage);
  }, [conversationIdentifier, currentMessage]);

  useEffect(() => {
    if (isChatLoaded && chatSectionRef.current) {
      chatSectionRef.current.scrollTop = chatSectionRef.current.scrollHeight;
      form.setFocus("text");
    }
  }, [isChatLoaded, messageCount]);

  return (
    <div className="flex flex-col h-[calc(100vh-5rem)] bg-slate-50">
      <section className="flex justify-center py-4">
        <p className="font-semibold">{isChatLoaded && conversation ? conversation.title : "Loading"}</p>
      </section>

      <section ref={chatSectionRef} className="flex-grow p-4 overflow-y-auto">
        <div className="flex flex-col gap-2">
          {isChatLoaded &&
            participantId &&
            bucketedMessages.map((bucket) => (
              <MessageBucket
                key={bucket.key}
                name={bucket.key}
                messages={bucket.messages}
                currentUserId={participantId}
                openEmojiBarMessageId={openEmojiBarMessageId}
                onOpenEmojiBar={(messageId: number) => setOpenEmojiBarMessageId(messageId)}
                onCloseEmojiBar={() => setOpenEmojiBarMessageId(undefined)}
                handleEmojiReact={(messageId, emojiKey) => emojiReact(messageId, emojiKey)}
              />
            ))}
          {otherParticipantIsTyping && (
            <div className="flex items-center px-4 rounded-xl bg-white shadow w-fit border">
              <span className="loading loading-dots w-4 text-slate-600"></span>
            </div>
          )}
        </div>
      </section>

      <section className="flex justify-end mx-4">
        <Button
          type="button"
          variant="ghost"
          title={`Use ${useLargeInput ? "normal" : "large"} input size`}
          className={`py-0 text-slate-400 hover:bg-transparent ${useLargeInput && "rotate-180"}`}
          onMouseDown={() => setUseLargeInput(!useLargeInput)}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth="2"
            stroke="currentColor"
            className="size-4"
          >
            <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 15.75 7.5-7.5 7.5 7.5" />
          </svg>
        </Button>
      </section>

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className={`relative px-4 pb-4 ${useLargeInput && "h-[28rem]"}`}
        >
          <FormField
            control={form.control}
            name="text"
            render={({ field }) => (
              <FormItem className="space-y-0 h-full">
                <FormControl>
                  <textarea
                    placeholder={isChatLoaded ? "Write a message..." : "Loading..."}
                    className={`p-4 rounded-lg w-full h-full resize-none shadow ${!canSendMessage && "disabled:bg-slate-100 shadow-none"}`}
                    disabled={!canSendMessage}
                    {...field}
                    onKeyDown={(e) => {
                      if ((e.ctrlKey || e.metaKey) && e.key === "Enter") {
                        e.preventDefault();
                        form.handleSubmit(onSubmit)();
                      }
                    }}
                    onChange={(e) => {
                      handleTypingDetection(e);
                      field.onChange(e);
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
