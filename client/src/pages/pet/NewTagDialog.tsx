import { Button } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
  DialogDescription,
} from "@/components/ui/dialog.tsx";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { ReactNode, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useAuth } from "@/hooks/useAuth.tsx";
import { Pet } from "@/api/types.ts";
import { useToast } from "@/hooks/use-toast.ts";

const addTagFormSchema = z.object({
  key: z.string(),
  value: z.string(),
});

type AddTagFormInputs = z.infer<typeof addTagFormSchema>;

interface NewTagDialogProps {
  pet: Pet;
  children: ReactNode;
}

async function createNewTag(petId: string, key: string, value: string, token: string): Promise<Pet> {
  return fetch(`${import.meta.env.VITE_API_BASE_URL}/pets/${petId}/tag`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ key, value }),
  }).then((res) => res.json());
}

function NewTagDialog({ pet, children }: NewTagDialogProps) {
  const { session } = useAuth();
  const { toast } = useToast();
  const [open, setOpen] = useState(false);

  const form = useForm<AddTagFormInputs>({
    resolver: zodResolver(addTagFormSchema),
    defaultValues: {
      key: "",
      value: "",
    },
  });

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (data: AddTagFormInputs) =>
      createNewTag(pet.id, data.key, data.value, session?.access_token ?? ""),
    onSuccess: (created: Pet) => {
      queryClient.invalidateQueries({ queryKey: ["pet"] });
      setOpen(false);
      form.reset();
      toast({
        title: "Success",
        description: `A new tag has been added for ${created.name}!`,
      });
    },
    onError: () => {
      toast({
        title: "Something went wrong",
        description: "There has been an issue adding a tag for your pet.",
      });
    },
  });

  async function onSubmit(values: AddTagFormInputs) {
    mutation.mutate(values);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="w-[95%] rounded-lg">
        <DialogTitle>New a new tag</DialogTitle>
        <DialogDescription>
          Add a new tag such as the breed, age, or any other information you would like people to know about{" "}
          {pet.name}.
        </DialogDescription>
        <Form {...form}>
          <form className="flex flex-col gap-4 justify-end" onSubmit={form.handleSubmit(onSubmit)}>
            <FormField
              control={form.control}
              name="key"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Label</FormLabel>
                  <FormControl>
                    <Input placeholder="Breed" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="value"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Value</FormLabel>
                  <FormControl>
                    <Input placeholder="Labrador" {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <Button type="submit">Add</Button>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

export default NewTagDialog;
