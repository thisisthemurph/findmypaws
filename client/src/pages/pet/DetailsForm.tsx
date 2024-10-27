import { z } from "zod";
import { useToast } from "@/hooks/use-toast.ts";
import { Pet } from "@/api/types.ts";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form.tsx";
import { Textarea } from "@/components/ui/textarea.tsx";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover.tsx";
import { Button } from "@/components/ui/button.tsx";
import { cn } from "@/lib/utils.ts";
import { format } from "date-fns";
import { Calendar } from "@/components/ui/calendar.tsx";
import { Input } from "@/components/ui/input.tsx";

const updateFormSchema = z.object({
  name: z.string().min(1, "A name must be provided"),
  blurb: z.string(),
  dob: z.date().optional(),
});

type UpdateFormSchema = z.infer<typeof updateFormSchema>;

interface DetailsFormProps {
  pet: Pet;
}

function DetailsForm({ pet }: DetailsFormProps) {
  const { toast } = useToast();

  function onSubmitUpdate(data: UpdateFormSchema) {
    console.log({ submitted: data });
    toast({
      title: "You submitted the following values:",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        </pre>
      ),
    });
  }

  const form = useForm<UpdateFormSchema>({
    resolver: zodResolver(updateFormSchema),
    defaultValues: {
      name: "",
      blurb: pet?.blurb ?? "",
      dob: undefined,
    },
  });

  const { isDirty } = form.formState;

  useEffect(() => {
    if (pet) {
      form.reset({
        name: pet.name ?? "",
        blurb: pet.blurb ?? "",
        dob: pet.dob ? new Date(pet.dob) : undefined,
      });
    }
  }, [pet, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmitUpdate)} className="flex flex-col gap-4 mt-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input
                  placeholder={`What is your ${pet?.type.toLocaleLowerCase() ?? "pet"} called?`}
                  {...field}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="blurb"
          render={({ field }) => (
            <FormItem>
              <FormLabel>About</FormLabel>
              <FormControl>
                <Textarea placeholder={`What would you like to say about ${pet.name}?`} {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="dob"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Date of birth</FormLabel>
              <div className="flex gap-2">
                <Popover>
                  <PopoverTrigger asChild>
                    <FormControl>
                      <Button
                        variant="outline"
                        className={cn(
                          "w-[240px] flex justify-between pl-3 text-left font-normal",
                          !field.value && "text-muted-foreground"
                        )}
                      >
                        {field.value ? format(field.value, "PPP") : <span>Pick a date</span>}
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          fill="none"
                          viewBox="0 0 24 24"
                          strokeWidth={1.5}
                          stroke="currentColor"
                          className="w-6 h-6"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5"
                          />
                        </svg>
                      </Button>
                    </FormControl>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={field.value}
                      onSelect={field.onChange}
                      disabled={(date) => {
                        const today = new Date();
                        const thirtyYearsAgo = new Date(
                          today.getFullYear() - 30,
                          today.getMonth(),
                          today.getDate()
                        );
                        return date > today || date < thirtyYearsAgo;
                      }}
                      initialFocus
                    />
                  </PopoverContent>
                </Popover>
                <Button type="button" variant="outline">
                  Clear
                </Button>
              </div>
            </FormItem>
          )}
        />
        <Button type="submit" disabled={!isDirty}>
          Update
        </Button>
      </form>
    </Form>
  );
}

export default DetailsForm;