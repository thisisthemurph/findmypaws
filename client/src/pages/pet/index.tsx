import { useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import { useAuth } from "@/hooks/useAuth.tsx";
import { Wrapper } from "@/components/Wrapper.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import NewTagDialog from "@/pages/pet/NewTagDialog.tsx";

function PetPage() {
  const { id } = useParams();
  const { session } = useAuth();

  const {
    isPending,
    error,
    data: pet,
  } = useQuery<Pet>({
    queryKey: ["pet"],
    queryFn: () =>
      fetch(`${import.meta.env.VITE_API_BASE_URL}/pets/${id}`, {
        method: "GET",
        headers: { Authorization: `Bearer ${session?.access_token}` },
      }).then((res) => res.json()),
  });

  if (error) {
    return <pre>{JSON.stringify(error, null, 2)}</pre>;
  }

  if (isPending || !pet) {
    return <p>Walking your pets to you...</p>;
  }

  return (
    <Wrapper>
      <section className="flex flex-col gap-4 items-center justify-center">
        <img
          className="shadow-2xl rounded-full"
          src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCHo3CkaH0oRY3MvrEN0xgn-x_Lsn3Lm3lVQ&s"
          alt={pet.name}
        />
        <p className="text-4xl capitalize font-semibold text-center">{pet.name}</p>
      </section>
      <section className="flex justify-center gap-2">
        <div className="flex justify-center gap-2">
          {Object.entries(pet.tags ?? {}).map(([key, value]) => (
            <Badge key={key} title={`tag type: ${key}`} variant="outline" size="lg">
              {value}
            </Badge>
          ))}
          {Object.entries(pet.tags ?? {}).length > 0 && (
            <NewTagDialog petName={pet.name}>
              <button>
                <Badge title="Add a new tag" variant="secondary" size="lg" className="hover:shadow-lg">
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
                </Badge>
              </button>
            </NewTagDialog>
          )}
        </div>
        {Object.entries(pet.tags ?? {}).length === 0 && (
          <NewTagDialog petName={pet.name}>
            <Button size="sm" variant="outline">
              Add a new tag
            </Button>
          </NewTagDialog>
        )}
      </section>
    </Wrapper>
  );
}

export default PetPage;
