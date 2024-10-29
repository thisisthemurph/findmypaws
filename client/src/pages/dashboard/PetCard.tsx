import { Card, CardContent, CardHeader } from "@/components/ui/card.tsx";
import { Link } from "react-router-dom";
import { ReactNode } from "react";
import { Skeleton } from "@/components/ui/skeleton.tsx";

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

function PetCardSkeleton() {
  return (
    <div className="flex items-center space-x-4 border p-3 rounded-lg">
      <Skeleton className="h-12 w-12 rounded-full" />
      <div className="space-y-2">
        <Skeleton className="h-4 w-[250px]" />
        <Skeleton className="h-4 w-[200px]" />
      </div>
    </div>
  );
}

PetCard.Header = PetCardHeader;
PetCard.Content = PetCardContent;
PetCard.Skeleton = PetCardSkeleton;
