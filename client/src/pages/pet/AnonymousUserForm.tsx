import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useApi } from "@/hooks/useApi.ts";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";
import useAnonymousUser from "@/hooks/useAnonymousUser.ts";
import { AnonymousUser } from "@/api/types.ts";
import { useNavigate } from "react-router-dom";

const formSchema = z.object({
  anonymousUserName: z.string().min(1, "A name is required"),
});

type FormSchema = z.infer<typeof formSchema>;

interface AnonymousUserFormProps {
  conversationIdentifier: string;
}

export default function AnonymousUserForm({ conversationIdentifier }: AnonymousUserFormProps) {
  const api = useApi();
  const [anonymousUserId] = useAnonymousUser();
  const navigate = useNavigate();

  const form = useForm<FormSchema>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      anonymousUserName: "",
    },
  });

  async function onSubmit(values: FormSchema) {
    try {
      // Update the anonymous user's name.
      await api<AnonymousUser>(`/user/anonymous/${anonymousUserId}`, {
        method: "PUT",
        body: JSON.stringify({
          name: values.anonymousUserName,
        }),
      });

      // Create the conversation if it does not exist.
      await api<void>(`/conversations`, {
        method: "POST",
        body: JSON.stringify({
          identifier: conversationIdentifier,
          participantId: anonymousUserId,
        }),
      });

      navigate(`/conversations/${conversationIdentifier}`);
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-2">
        <FormField
          control={form.control}
          name="anonymousUserName"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="What would you like to be called?" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" disabled={!form.formState.isDirty || !form.formState.isValid} className="w-fit">
          Start chatting
        </Button>
      </form>
    </Form>
  );
}
