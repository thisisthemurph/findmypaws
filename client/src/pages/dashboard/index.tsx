import { useFetch } from "@/hooks/useFetch.ts";
import { Pet } from "@/api/types.ts";
import { useQuery } from "@tanstack/react-query";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog.tsx";
import { Button } from "@/components/ui/button.tsx";
import NewPetForm from "@/pages/dashboard/NewPetForm.tsx";
import { useState } from "react";
import PetCard from "@/pages/dashboard/PetCard.tsx";

export default function DashboardPage() {
  const fetch = useFetch();
  const [isOpen, setIsOpen] = useState(false);

  const { isLoading, data } = useQuery<Pet[]>({
    queryKey: ["pets"],
    queryFn: () => fetch("/pets"),
  });

  return (
    <>
      <h1>Dashboard</h1>

      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogTrigger asChild>
          <Button variant="outline">Add a new pet</Button>
        </DialogTrigger>
        <DialogContent className="w-[95%] rounded-lg" aria-description="add bew pet from">
          <DialogHeader className="text-left">
            <DialogTitle>Add a new pet</DialogTitle>
            <DialogDescription>Use this form to add a new pets details...</DialogDescription>
          </DialogHeader>
          <NewPetForm onFormComplete={() => setIsOpen(false)} />
        </DialogContent>
      </Dialog>

      <section>
        <h2>Your pets</h2>
        <div className="flex flex-col gap-4">
          {isLoading && <p>Walking your pets</p>}
          {data && data.map((pet) => <PetCard key={pet.id} pet={pet} />)}
        </div>
      </section>
    </>
  );
}
