import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem } from "@/components/ui/form.tsx";
import { Button } from "@/components/ui/button.tsx";
import { useAuth } from "@clerk/clerk-react";
import useAnonymousUser from "@/hooks/useAnonymousUser.ts";
import useChat from "@/components/useChat.ts";
import MessageBucket from "@/pages/messenger/MessageBucket.tsx";

const formSchema = z.object({
  text: z.string().min(1, "Enter a message"),
});

type FormInputs = z.infer<typeof formSchema>;

export default function ConversationPage() {
  const { userId, isLoaded: isUserLoaded } = useAuth();
  const [anonymousUserId] = useAnonymousUser();
  const { bucketedMessages, sendMessage } = useChat("49b6d8d8-816c-4628-9238-fba78ab18c90");

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

  if (!isUserLoaded) {
    return <p>Loading user</p>;
  }

  const participantId = userId ?? anonymousUserId;
  if (!participantId) {
    return <p>Cannot determine the effective user ID</p>;
  }

  return (
    <div className="flex flex-col h-[calc(100vh-5rem)] bg-slate-50">
      <section className="flex justify-center py-4">
        <p className="font-semibold">Mike Murphy</p>
      </section>

      <section className="flex-grow p-4 overflow-y-auto">
        <div className="flex flex-col gap-1">
          {bucketedMessages.map((bucket) => (
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
                    placeholder="Write a message..."
                    className="p-4 rounded-lg w-full resize-none shadow"
                    {...field}
                  ></textarea>
                </FormControl>
              </FormItem>
            )}
          />
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
        </form>
      </Form>
    </div>
  );
}
