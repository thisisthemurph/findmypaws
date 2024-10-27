import { PetCard } from "@/pages/dashboard/PetCard.tsx";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";
import { ButtonHTMLAttributes } from "react";

export default function NewPetButton(props: ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button {...props}>
      <PetCard>
        <PetCard.Header>
          <Avatar>
            <AvatarFallback>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-4"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
              </svg>
            </AvatarFallback>
          </Avatar>
        </PetCard.Header>
        <PetCard.Content>
          <p className="font-semibold">Add your next pet</p>
        </PetCard.Content>
      </PetCard>
    </button>
  );
}
