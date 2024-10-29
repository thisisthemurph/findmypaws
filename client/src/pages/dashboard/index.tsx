import { useApi } from "@/hooks/useApi.ts";
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
import { PetCard } from "@/pages/dashboard/PetCard.tsx";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar.tsx";
import NewPetButton from "@/pages/dashboard/NewPetButton.tsx";
import { useNavigate } from "react-router-dom";

export default function DashboardPage() {
  const api = useApi();
  const navigate = useNavigate();

  const { isLoading, data } = useQuery<Pet[]>({
    queryKey: ["pets"],
    queryFn: () => api("/pets"),
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
                    <AvatarImage src={`${import.meta.env.VITE_API_BASE_URL}/pets/${pet.id}/avatar`} />
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

          <Dialog>
            <DialogTrigger asChild>
              <NewPetButton>
                {(data?.length || 0) === 0 ? "Add your first pet" : "Add a new pet"}
              </NewPetButton>
            </DialogTrigger>
            <DialogContent className="w-[95%] rounded-lg" aria-description="add bew pet from">
              <DialogHeader className="text-left">
                <DialogTitle>Add a new pet</DialogTitle>
                <DialogDescription>Add a new pet to your kennel...</DialogDescription>
              </DialogHeader>
              <NewPetForm onCreated={(created) => navigate(`/pet/${created.id}`)} />
            </DialogContent>
          </Dialog>
        </div>
      </section>
    </>
  );
}
