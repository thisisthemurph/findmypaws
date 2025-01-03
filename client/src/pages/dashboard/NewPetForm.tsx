import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select.tsx";
import { Button } from "@/components/ui/button.tsx";
import { z } from "zod";
import { useApi } from "@/hooks/useApi.ts";
import { useToast } from "@/hooks/use-toast.ts";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";

const newPetFormSchema = z.object({
  name: z.string().min(1, "The name of your pet is required"),
  type: z.string().min(1, "Please select a type for your pet"),
});

type NewPetFormInputs = z.infer<typeof newPetFormSchema>;

interface NewPetFormProps {
  onCreated: (created: Pet) => void;
}

export default function NewPetForm({ onCreated }: NewPetFormProps) {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const api = useApi();

  const form = useForm<NewPetFormInputs>({
    resolver: zodResolver(newPetFormSchema),
    defaultValues: {
      name: "",
      type: "",
    },
    mode: "onChange",
  });

  async function handleCreatePet(newPet: NewPetFormInputs) {
    return await api<Pet>("/pets", {
      method: "POST",
      body: JSON.stringify(newPet),
    });
  }

  const createPetMutation = useMutation({
    mutationFn: async (newPet: NewPetFormInputs) => handleCreatePet(newPet),
    onSuccess: async (created: Pet) => {
      await queryClient.invalidateQueries({ queryKey: ["pets"] });
      onCreated(created);
    },
    onError: (error: Error) => {
      toast({
        title: "Something went wrong",
        description: error?.message || "There has been an issue adding your pet.",
        variant: "destructive",
      });
    },
  });

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((values) => createPetMutation.mutate(values))}
        className="flex flex-col gap-4"
      >
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
              <FormLabel htmlFor="pet-type-selector">Type</FormLabel>
              <Select onValueChange={field.onChange} defaultValue={field.value}>
                <SelectTrigger id="pet-type-selector">
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

        <Button type="submit" disabled={!form.formState.isValid}>
          Add
        </Button>
      </form>
    </Form>
  );
}
