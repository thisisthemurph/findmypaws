import { useForm } from "react-hook-form";
import { Form, FormControl, FormField, FormItem, FormLabel } from "@/components/ui/form.tsx";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Input } from "@/components/ui/input.tsx";
import { useToast } from "@/hooks/use-toast.ts";
import { Button } from "@/components/ui/button.tsx";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select.tsx";

const newPetFormSchema = z.object({
  name: z.string().min(1, "The name of your pet is required"),
  type: z.string(),
});

type NewPetFormInputs = z.infer<typeof newPetFormSchema>;

interface NewPetFormProps {
  onFormComplete: () => void;
}

function NewPetForm({ onFormComplete }: NewPetFormProps) {
  const { toast } = useToast();
  const form = useForm<NewPetFormInputs>({
    resolver: zodResolver(newPetFormSchema),
    defaultValues: {
      name: "",
      type: "",
    },
  });

  async function onSubmit(values: NewPetFormInputs) {
    console.table([values]);
    toast({
      title: "Success",
      description: `You had added ${values.name} to your list of pets.`,
      variant: "default",
    });
    onFormComplete();
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
