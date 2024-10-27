import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Link } from "react-router-dom";
import { ReactNode } from "react";

interface PetCardProps {
  petId?: string;
  children: ReactNode;
}

export function PetCard({ petId, children }: PetCardProps) {
  if (!petId) {
    return <Card className="flex hover:shadow-lg min-w-[14rem]">{children}</Card>;
  }
  return (
    <Link to={`/pet/${petId}`}>
      <Card className="flex hover:shadow-lg min-w-[14rem]">{children}</Card>
    </Link>
  );
}

function PetCardHeader({ children }: { children: ReactNode }) {
  return <CardHeader className="flex items-center justify-center p-4">{children}</CardHeader>;
}

function PetCardContent({ children }: { children: ReactNode }) {
  return <CardContent className="flex flex-col justify-center p-4 pl-0">{children}</CardContent>;
}

PetCard.Header = PetCardHeader;
PetCard.Content = PetCardContent;
