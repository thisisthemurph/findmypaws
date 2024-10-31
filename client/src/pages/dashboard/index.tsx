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
import { PageHeading } from "@/components/PageHeading.tsx";

export default function DashboardPage() {
  const api = useApi();
  const navigate = useNavigate();

  const { isLoading, data } = useQuery<Pet[]>({
    queryKey: ["pets"],
    queryFn: () => api("/pets"),
  });

  return (
    <div className="p-4">
      <PageHeading heading="Your kennel" subheading="Maintain your pets in your kennel..." />
      <section className="flex flex-col sm:flex-row flex-wrap gap-4">
        {isLoading || !data ? (
          <>
            <PetCard.Skeleton />
            <PetCard.Skeleton />
          </>
        ) : (
          data.map((pet) => (
            <PetCard key={pet.id} petId={pet.id}>
              <PetCard.Header>
                <Avatar>
                  <AvatarImage src={`${import.meta.env.VITE_API_BASE_URL}/pets/${pet.id}/avatar`} />
                  <AvatarFallback>{pet.name[0]}</AvatarFallback>
                </Avatar>
              </PetCard.Header>
              <PetCard.Content>
                <p className="font-semibold">{pet.name}</p>
                {pet.blurb && (
                  <p className="text-sm text-slate-600">
                    {pet.blurb.slice(0, 80)}
                    {pet.blurb.length > 80 && "..."}
                  </p>
                )}
              </PetCard.Content>
            </PetCard>
          ))
        )}

        <Dialog>
          <DialogTrigger asChild>
            <NewPetButton>{(data?.length || 0) === 0 ? "Add your first pet" : "Add a new pet"}</NewPetButton>
          </DialogTrigger>
          <DialogContent className="w-[95%] rounded-lg" aria-description="add bew pet from">
            <DialogHeader className="text-left">
              <DialogTitle>Add a new pet</DialogTitle>
              <DialogDescription>Add a new pet to your kennel...</DialogDescription>
            </DialogHeader>
            <NewPetForm onCreated={(created) => navigate(`/pet/${created.id}`)} />
          </DialogContent>
        </Dialog>
      </section>
    </div>
  );
}
