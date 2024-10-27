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
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { ReactNode, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import { useToast } from "@/hooks/use-toast.ts";
import { useFetch } from "@/hooks/useFetch.ts";

const addTagFormSchema = z.object({
  key: z.string().min(1, "A label must be provided."),
  value: z.string().min(1, "A value must be provided."),
});

export type AddTagFormInputs = z.infer<typeof addTagFormSchema>;

interface NewTagDialogProps {
  pet: Pet;
  children: ReactNode;
}

function NewTagDialog({ pet, children }: NewTagDialogProps) {
  // const createNewTag = useFetch(`/pets/${pet.id}/tag`, "POST");
  const fetch = useFetch();
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
    mutationFn: async (data: AddTagFormInputs) =>
      await fetch<Pet>(`/pets/${pet.id}/tag`, {
        method: "POST",
        body: JSON.stringify(data),
      }),
    onSuccess: async (created: Pet) => {
      await queryClient.invalidateQueries({ queryKey: ["pet"] });
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
        <DialogTitle>Add a new tag</DialogTitle>
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
                  <FormMessage />
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
                  <FormMessage />
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
