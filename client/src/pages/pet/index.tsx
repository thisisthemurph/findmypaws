import { useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import { useAuth } from "@/hooks/useAuth.tsx";
import { Wrapper } from "@/components/Wrapper.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import NewTagDialog from "@/pages/pet/NewTagDialog.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { useToast } from "@/hooks/use-toast.ts";

async function deleteTag(petId: string, key: string, token: string) {
  return fetch(`${import.meta.env.VITE_API_BASE_URL}/pets/${petId}/tag/${key}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` },
  }).then((res) => res.json());
}

function PetPage() {
  const { id } = useParams();
  const { session } = useAuth();
  const { toast } = useToast();

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

  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: (key: string) => deleteTag(pet?.id ?? "", key, session?.access_token ?? ""),
    onSuccess: async () => {
      queryClient.invalidateQueries({ queryKey: ["pet"] });
      toast({
        title: "Deleted",
        description: "You have deleted the tag!",
      });
    },
    onError: () => {
      toast({
        title: "Something went wrong",
        description: "There has been an issue adding a tag for your pet.",
        variant: "destructive",
      });
    },
  });

  if (error) {
    return <pre>{JSON.stringify(error, null, 2)}</pre>;
  }

  if (isPending || !pet) {
    return <p>Walking your pets to you...</p>;
  }

  return (
    <Wrapper>
      <section className="flex flex-col gap-4 items-center justify-center mb-4">
        <img
          className="shadow-2xl rounded-full"
          src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCHo3CkaH0oRY3MvrEN0xgn-x_Lsn3Lm3lVQ&s"
          alt={pet.name}
        />
        <p className="text-4xl capitalize font-semibold text-center">{pet.name}</p>
      </section>
      <section className="flex justify-center gap-2">
        <div className="flex flex-wrap justify-center gap-2">
          {Object.entries(pet.tags ?? {}).map(([key, value]) => (
            <PetBadge
              key={key}
              tagKey={key}
              tagValue={value}
              handleDelete={(key) => {
                mutation.mutate(key);
              }}
            />
          ))}
          {Object.entries(pet.tags ?? {}).length > 0 && (
            <NewTagDialog pet={pet}>
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
          <NewTagDialog pet={pet}>
            <Button size="sm" variant="outline">
              Add a new tag
            </Button>
          </NewTagDialog>
        )}
      </section>
    </Wrapper>
  );
}

function PetBadge({
  tagKey,
  tagValue,
  handleDelete,
}: {
  tagKey: string;
  tagValue: string;
  handleDelete: (key: string) => void;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Badge key={tagKey} title={`tag type: ${tagKey}`} variant="outline" size="lg">
          {tagValue}
        </Badge>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>Tag actions</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={() => handleDelete(tagKey)}>Delete</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

export default PetPage;
