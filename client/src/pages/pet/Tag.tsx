import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { ReactNode } from "react";

interface TagProps {
  identifier: string;
  children?: ReactNode;
  handleDelete: (key: string) => void;
}

function Tag({ identifier, children, handleDelete }: TagProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Badge key={identifier} title={`tag type: ${identifier}`} variant="outline" size="lg">
          {children}
        </Badge>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>Tag actions</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => handleDelete(identifier)}>Delete</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export default Tag;
