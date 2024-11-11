import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form.tsx";
import { Button } from "@/components/ui/button.tsx";
import useChat from "@/components/useChat.ts";
import MessageBucket from "@/pages/messenger/MessageBucket.tsx";
import { useParams } from "react-router-dom";

const formSchema = z.object({
  text: z.string().min(1, "Enter a message"),
});

type FormInputs = z.infer<typeof formSchema>;

export default function ConversationPage() {
  const { conversationIdentifier } = useParams();
  const {
    title,
    participantId,
    bucketedMessages,
    sendMessage,
    isLoaded: isChatLoaded,
  } = useChat(conversationIdentifier!);

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

  const canSendMessage = isChatLoaded;

  return (
    <div className="flex flex-col h-[calc(100vh-5rem)] bg-slate-50">
      <section className="flex justify-center py-4">
        <p className="font-semibold">{isChatLoaded ? title : "Loading"}</p>
      </section>

      <section className="flex-grow p-4 overflow-y-auto">
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
          {form.getValues("text") && (
            <Button type="submit" variant="ghost" size="icon" className="absolute top-5 right-5 rounded-full">
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
          )}
        </form>
      </Form>
    </div>
  );
}
