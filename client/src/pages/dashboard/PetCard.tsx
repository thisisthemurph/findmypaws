import { Pet } from "@/api/types.ts";
import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Link } from "react-router-dom";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar.tsx";

interface PetCardProps {
  pet: Pet;
}

function PetCard({ pet }: PetCardProps) {
  return (
    <Link to={`/pet/${pet.id}`}>
      <Card className="flex hover:shadow-lg">
        <CardHeader className="flex items-center justify-center p-4">
          <Avatar>
            <AvatarImage src={`${import.meta.env.VITE_BASE_URL}/${pet.avatar}`} />
            <AvatarFallback>{pet.name[0]}</AvatarFallback>
          </Avatar>
        </CardHeader>
        <CardContent className="p-4 pl-0">
          <p className="font-semibold">{pet.name}</p>
          <p className="text-sm text-slate-600">This is an example of a description...</p>
        </CardContent>
      </Card>
    </Link>
  );
}

export default PetCard;
