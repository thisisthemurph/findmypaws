import { ButtonHTMLAttributes, forwardRef } from "react";
import { PetCard } from "@/pages/dashboard/PetCard.tsx";
import { Avatar, AvatarFallback } from "@/components/ui/avatar.tsx";

const NewPetButton = forwardRef<HTMLButtonElement, ButtonHTMLAttributes<HTMLButtonElement>>(
  ({ children, ...props }, ref) => {
    return (
      <button ref={ref} {...props}>
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
            <p className="font-semibold">{children}</p>
          </PetCard.Content>
        </PetCard>
      </button>
    );
  }
);

export default NewPetButton;
