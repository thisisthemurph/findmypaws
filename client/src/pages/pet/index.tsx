import { useParams } from "react-router-dom";
import { useFetch } from "@/hooks/useFetch.ts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import PetAvatar from "@/pages/pet/PetAvatar.tsx";
import Tag from "@/pages/pet/Tag.tsx";
import NewTagDialog from "@/pages/pet/NewTagDialog.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { Button } from "@/components/ui/button.tsx";
import DetailsForm from "@/pages/pet/DetailsForm.tsx";
import { useToast } from "@/hooks/use-toast.ts";

export default function PetPage() {
  const { id } = useParams();
  const { toast } = useToast();
  const fetch = useFetch();
  const queryClient = useQueryClient();

  const { data: pet, isLoading } = useQuery<Pet>({
    queryKey: ["pet"],
    queryFn: () => fetch<Pet>(`/pets/${id}`),
  });

  const deleteTagMutation = useMutation({
    mutationFn: async (deleteRequest: { petId: string; key: string }) =>
      await fetch<void>(`/pets/${deleteRequest.petId}/tag/${deleteRequest.key}`, { method: "DELETE" }),
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
    mutationFn: async (file: File) => {
      const allowedMimeTypes = ["image/jpeg", "image/png"];
      if (!allowedMimeTypes.includes(file.type)) {
        throw new Error("Only JPEG and PNG files are allowed");
      }

      const formData = new FormData();
      formData.append("file", file);
      formData.append("fileName", file.name);

      await fetch<{ avatar_url: string }>(`/pets/${id}/avatar`, {
        method: "PUT",
        body: formData,
      });
    },
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

  async function onAvatarChange(file: File) {
    updateAvatarMutation.mutate(file);
  }

  if (isLoading || !pet) {
    return "Loading";
  }

  return (
    <>
      <PetAvatar pet={pet} changeAvatar={onAvatarChange} />
      <section className="flex justify-center gap-2">
        <div className="flex flex-wrap justify-center gap-2">
          {Object.entries(pet.tags ?? {}).map(([key, value]) => (
            <Tag
              key={key}
              identifier={key}
              handleDelete={() => deleteTagMutation.mutate({ petId: pet.id, key })}
            >
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
      <DetailsForm pet={pet} />
    </>
  );
}
