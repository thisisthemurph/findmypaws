import { Wrapper } from "@/components/Wrapper.tsx";
import { useAuth } from "@/hooks/useAuth.tsx";
import NewPetForm from "@/pages/profile/NewPetForm.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import PetCard from "@/pages/profile/PetCard.tsx";

function ProfilePage() {
  const { user, session } = useAuth();
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const { isPending, data } = useQuery<Pet[]>({
    queryKey: ["pets"],
    queryFn: () =>
      fetch(`${import.meta.env.VITE_API_BASE_URL}/pets`, {
        method: "GET",
        headers: { Authorization: `Bearer ${session?.access_token}` },
      }).then((res) => {
        if (!res.ok) {
          throw new Error("There was a problem fetching your pets");
        }
        return res.json();
      }),
  });

  return (
    <Wrapper className="flex flex-col gap-6">
      <h1>Welcome back {user?.name}!</h1>

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
          {isPending && <p>Walking your pets</p>}
          {data && data.map((pet) => <PetCard key={pet.id} pet={pet} />)}
        </div>
      </section>
    </Wrapper>
  );
}

export default ProfilePage;
