import { useNavigate, useParams } from "react-router-dom";
import { useApi } from "@/hooks/useApi.ts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import PetAvatar from "@/pages/pet/PetAvatar.tsx";
import Tag from "@/pages/pet/Tag.tsx";
import NewTagDialog from "@/pages/pet/NewTagDialog.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { Button } from "@/components/ui/button.tsx";
import DetailsForm from "@/pages/pet/DetailsForm.tsx";
import { useToast } from "@/hooks/use-toast.ts";
import { ToastAction } from "@/components/ui/toast.tsx";

export default function PetPage() {
  const { id } = useParams();
  const { toast } = useToast();
  const navigate = useNavigate();
  const api = useApi();
  const queryClient = useQueryClient();

  const { data: pet, isLoading } = useQuery<Pet>({
    queryKey: ["pet"],
    queryFn: () => api<Pet>(`/pets/${id}`),
  });

  const deleteMutation = useMutation({
    mutationFn: async () => {
      if (!pet) {
        throw new Error("No pet found");
      }
      return await api<void>(`/pets/${pet.id}`, { method: "DELETE" });
    },
    onSuccess: () => {
      toast({
        title: "Deleted",
        description: pet ? (
          <>
            <strong>{pet.name}</strong> has been deleted from your kennel.
          </>
        ) : (
          "Your pet has been deleted from your kennel."
        ),
      });
      navigate("/dashboard");
    },
    onError: (error: Error) =>
      toast({
        title: "Something went wrong",
        description: error?.message || "There has been an issue deleting your pet.",
        variant: "destructive",
      }),
  });

  const deleteTagMutation = useMutation({
    mutationFn: async (deleteRequest: { petId: string; key: string }) =>
      await api<void>(`/pets/${deleteRequest.petId}/tag/${deleteRequest.key}`, { method: "DELETE" }),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["pet"] });
      toast({
        title: "Deleted",
        description: `The tag is no longer associated with ${pet?.name || "your pet"}.`,
      });
    },
    onError: (error: Error) => {
      toast({
        title: "Something went wrong",
        description: error?.message || "There has been an issue deleting the tag.",
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

      await api<{ avatar_url: string }>(`/pets/${id}/avatar`, {
        method: "PUT",
        body: formData,
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["pet"] });
      toast({
        title: "Success",
        description: "Avatar updated successfully.",
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

      <Button
        variant="destructive"
        className="w-full mt-2"
        onMouseDown={() => {
          toast({
            title: "Delete",
            description: (
              <>
                Are you sure you want to delete <strong>{pet.name}</strong> from your kennel?
              </>
            ),
            action: (
              <div className="flex flex-col gap-1">
                <ToastAction altText="cancel deletion">Keep</ToastAction>
                <ToastAction altText="confirm deletion" onClick={() => deleteMutation.mutate()}>
                  Delete
                </ToastAction>
              </div>
            ),
          });
        }}
      >
        Delete
      </Button>
    </>
  );
}
