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
import NewPetForm from "@/pages/dashboard/NewPetForm.tsx";
import { useState } from "react";
import { PetCard } from "@/pages/dashboard/PetCard.tsx";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar.tsx";
import NewPetButton from "@/pages/dashboard/NewPetButton.tsx";

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

      <section>
        <h2>Your kennel</h2>
        <div className="flex flex-col sm:flex-row flex-wrap gap-4">
          {isLoading && <p>Walking your pets</p>}
          {data &&
            data.map((pet, index) => (
              <PetCard key={pet.id} petId={pet.id}>
                <PetCard.Header>
                  <Avatar>
                    <AvatarImage src={`${import.meta.env.VITE_BASE_URL}/${pet.avatar}`} />
                    <AvatarFallback>{pet.name[0]}</AvatarFallback>
                  </Avatar>
                </PetCard.Header>
                <PetCard.Content>
                  <p className="font-semibold">{pet.name}</p>
                  <p className="text-sm text-slate-600">
                    {index % 4 === 0 ? "This is an example of a description..." : "this is a test"}
                  </p>
                </PetCard.Content>
              </PetCard>
            ))}

          <Dialog open={isOpen} onOpenChange={setIsOpen}>
            <DialogTrigger asChild>
              <NewPetButton />
            </DialogTrigger>
            <DialogContent className="w-[95%] rounded-lg" aria-description="add bew pet from">
              <DialogHeader className="text-left">
                <DialogTitle>Add a new pet</DialogTitle>
                <DialogDescription>Use this form to add a new pets details...</DialogDescription>
              </DialogHeader>
              <NewPetForm onFormComplete={() => setIsOpen(false)} />
            </DialogContent>
          </Dialog>
        </div>
      </section>
    </>
  );
}
