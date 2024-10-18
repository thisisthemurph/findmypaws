import { Pet } from "@/api/types.ts";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card.tsx";
import { Badge } from "@/components/ui/badge.tsx";

interface PetCardProps {
  pet: Pet;
}

function PetCard({ pet }: PetCardProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{pet.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>This is only an example!</p>
      </CardContent>
      {pet.tags && Object.entries(pet.tags).length > 0 && (
        <CardFooter className="flex gap-1">
          {Object.entries(pet.tags).map(([k, v]) => (
            <Badge variant="outline" key={k}>
              {v}
            </Badge>
          ))}
        </CardFooter>
      )}
    </Card>
  );
}

export default PetCard;
