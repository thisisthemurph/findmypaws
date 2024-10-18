import { Pet } from "@/api/types.ts";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { Link } from "react-router-dom";

interface PetCardProps {
  pet: Pet;
}

function PetCard({ pet }: PetCardProps) {
  return (
    <Link to={`/pet/${pet.id}`}>
      <Card className="group hover:bg-accent hover:shadow-lg">
        <CardHeader>
          <CardTitle>{pet.name}</CardTitle>
        </CardHeader>
        <CardContent>
          <p>This is only an example!</p>
        </CardContent>
        {pet.tags && Object.entries(pet.tags).length > 0 && (
          <CardFooter className="flex flex-wrap gap-1">
            {Object.entries(pet.tags).map(([k, v]) => (
              <Badge variant="outline" key={k} className="bg-white">
                {v}
              </Badge>
            ))}
          </CardFooter>
        )}
      </Card>
    </Link>
  );
}

export default PetCard;
