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
import { ReactNode } from "react";

const addTagFormSchema = z.object({
  key: z.string(),
  value: z.string(),
});

type AddTagFormInputs = z.infer<typeof addTagFormSchema>;

interface NewTagDialogProps {
  petName: string;
  children: ReactNode;
}

function NewTagDialog({ petName, children }: NewTagDialogProps) {
  const form = useForm<AddTagFormInputs>({
    resolver: zodResolver(addTagFormSchema),
    defaultValues: {
      key: "",
      value: "",
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="w-[95%] rounded-lg">
        <DialogTitle>New a new tag</DialogTitle>
        <DialogDescription>
          Add a new tag such as the breed, age, or any other information you would like people to know about{" "}
          {petName}.
        </DialogDescription>
        <Form {...form}>
          <form className="flex flex-col gap-4 justify-end">
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
