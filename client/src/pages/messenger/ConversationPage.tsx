import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useApi } from "@/hooks/useApi.ts";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { ConversationWithMessages, Message } from "@/api/types.ts";

const formSchema = z.object({
  message: z.string().min(1, "Enter a message"),
});

type FormInputs = z.infer<typeof formSchema>;

export default function ConversationPage() {
  const api = useApi();
  const { conversationId } = useParams();

  const { data: conversation, isLoading } = useQuery<ConversationWithMessages>({
    queryKey: ["conversation", conversationId],
    queryFn: async () => await api<ConversationWithMessages>(`/conversations/${conversationId}`),
  });

  const form = useForm<FormInputs>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      message: "",
    },
  });

  function onSubmit(data: FormInputs) {
    console.log(data);
  }

  function localTimeString(isoDate: string): string {
    const date = new Date(isoDate);
    return date.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
      hour12: false,
    });
  }

  return (
    <div className="flex flex-col h-[calc(100vh-5rem)] bg-slate-50">
      <section className="flex justify-center py-4">
        <p className="font-semibold">Mike Murphy</p>
      </section>

      <section className="flex-grow p-4 overflow-y-auto">
        {!isLoading && conversation ? (
          <>
            <div className="flex flex-col gap-3">
              {conversation.messages.map((message) => (
                <MessageBubble message={message} />
              ))}
            </div>
            <div className="text-sm text-slate-600 ml-4 mt-1">
              {localTimeString(conversation.lastMessageAt)}
            </div>
          </>
        ) : (
          <p>Loading</p>
        )}
      </section>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="relative p-4">
          <FormField
            control={form.control}
            name="message"
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

function MessageBubble({ message }: { message: Message }) {
  const { text, outgoing } = message;
  return (
    <p
      className={`px-4 py-3 text-sm max-w-[80%] w-fit ${outgoing && "self-end"} ${outgoing ? "bg-[#7F00FF] text-white" : "bg-white"} rounded-full shadow`}
    >
      {text}
    </p>
  );
}
