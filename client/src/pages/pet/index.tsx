import { useParams } from "react-router-dom";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import { useAuth } from "@/hooks/useAuth.tsx";
import { Wrapper } from "@/components/Wrapper.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import NewTagDialog from "@/pages/pet/NewTagDialog.tsx";
import { useToast } from "@/hooks/use-toast.ts";
import PetAvatar from "@/pages/pet/PetAvatar.tsx";
import Tag from "@/pages/pet/Tag.tsx";

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
    queryFn: async () =>
      fetch(`${import.meta.env.VITE_API_BASE_URL}/pets/${id}`, {
        method: "GET",
        headers: { Authorization: `Bearer ${session?.access_token}` },
      }).then((res) => {
        if (!res.ok) {
          throw new Error("There has been an issue fetching pet");
        }
        return res.json();
      }),
  });

  const queryClient = useQueryClient();
  const deleteTagMutation = useMutation({
    mutationFn: (key: string) => deleteTag(pet?.id ?? "", key, session?.access_token ?? ""),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["pet"] });
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

  const updateAvatarMutation = useMutation({
    mutationFn: async (file: File) => await updateAvatar(file),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["pet"] });
      toast({
        title: "Updated",
        description: "Your avatar has been updated.",
      });
    },
    onError: (error: Error) => {
      toast({
        title: "Something went wrong",
        description: error?.message || "There has been an issue updating the avatar.",
        variant: "destructive",
      });
    },
  });

  async function updateAvatar(file: File) {
    const allowedTypes = ["image/jpeg", "image/png"];
    if (!allowedTypes.includes(file.type)) {
      throw new Error("Only JPEG and PNG file types are allowed.");
    }

    const url = `${import.meta.env.VITE_API_BASE_URL}/pets/${id}/avatar`;
    const formData = new FormData();
    formData.append("file", file);
    formData.append("fileName", file.name);

    return await fetch(url, {
      method: "PUT",
      headers: {
        Authorization: `Bearer ${session?.access_token ?? ""}`,
      },
      body: formData,
    }).then((res) => {
      if (!res.ok) {
        throw new Error("Error updating the avatar.");
      }
      return res.json();
    });
  }

  if (error) {
    return <pre>{JSON.stringify(error, null, 2)}</pre>;
  }

  if (isPending || !pet) {
    return <p>Walking your pets to you...</p>;
  }

  return (
    <Wrapper>
      <PetAvatar pet={pet} changeAvatar={async (file: File) => updateAvatarMutation.mutate(file)} />
      <section className="flex justify-center gap-2">
        <div className="flex flex-wrap justify-center gap-2">
          {Object.entries(pet.tags ?? {}).map(([key, value]) => (
            <Tag key={key} identifier={key} handleDelete={(key) => deleteTagMutation.mutate(key)}>
              {value}
            </Tag>
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

export default PetPage;
