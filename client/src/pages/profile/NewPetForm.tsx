import { useForm } from "react-hook-form";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input.tsx";
import { useToast } from "@/hooks/use-toast.ts";
import { Button } from "@/components/ui/button.tsx";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select.tsx";
import { useAuth } from "@/hooks/useAuth.tsx";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";

const newPetFormSchema = z.object({
  name: z.string().min(1, "The name of your pet is required"),
  type: z.string(),
});

type NewPetFormInputs = z.infer<typeof newPetFormSchema>;

interface NewPetFormProps {
  onFormComplete: () => void;
}

async function createPet(pet: NewPetFormInputs, token: string) {
  const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/pets`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(pet),
  });
  if (!response.ok) {
    throw new Error("Failed to create new pet");
  }
  return response.json();
}

function NewPetForm({ onFormComplete }: NewPetFormProps) {
  const { session } = useAuth();
  const { toast } = useToast();
  const form = useForm<NewPetFormInputs>({
    resolver: zodResolver(newPetFormSchema),
    defaultValues: {
      name: "",
      type: "",
    },
  });

  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: (newPet: NewPetFormInputs) => {
      return createPet(newPet, session?.access_token ?? "");
    },
    onSuccess: (created: Pet) => {
      queryClient.invalidateQueries({ queryKey: ["pets"] });
      toast({
        title: "Success",
        description:
          created.type !== "Unspecified"
            ? `Your ${created.type.toLocaleLowerCase()} ${created.name} has been added!`
            : `${created.name} has been added!`,
        variant: "default",
      });
      onFormComplete();
    },
    onError: () => {
      toast({
        title: "Something went wrong",
        description: "There has been an issue adding your pet.",
        variant: "destructive",
      });
    },
  });

  async function onSubmit(values: NewPetFormInputs) {
    mutation.mutate(values);
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-4">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="What do you call your pet?" {...field} />
              </FormControl>
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Type</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <SelectTrigger>
                  <SelectValue placeholder="Pet type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Cat">Cat</SelectItem>
                  <SelectItem value="Dog">Dog</SelectItem>
                  <SelectItem value="Unspecified">Other</SelectItem>
                </SelectContent>
              </Select>
            </FormItem>
          )}
        />

        <Button type="submit">Add</Button>
      </form>
    </Form>
  );
}

export default NewPetForm;
