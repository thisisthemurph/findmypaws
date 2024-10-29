import { useNavigate, useParams } from "react-router-dom";
import { useApi } from "@/hooks/useApi.ts";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Pet } from "@/api/types.ts";
import PetAvatar from "@/pages/pet/PetAvatar.tsx";
import { Button } from "@/components/ui/button.tsx";
import DetailsForm from "@/pages/pet/DetailsForm.tsx";
import { useToast } from "@/hooks/use-toast.ts";
import { ToastAction } from "@/components/ui/toast.tsx";
import { useAuth } from "@clerk/clerk-react";
import PublicDetails from "@/pages/pet/PublicDetails.tsx";
import TagsSection from "@/pages/pet/TagsSection.tsx";

export default function PetPage() {
  const { id } = useParams();
  const { toast } = useToast();
  const navigate = useNavigate();
  const auth = useAuth();
  const api = useApi();
  const queryClient = useQueryClient();

  const { data: pet, isLoading } = useQuery<Pet>({
    queryKey: ["pet"],
    queryFn: () => api<Pet>(`/pets/${id}`),
  });

  const userIsOwner = auth.userId === pet?.user_id;

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
      <TagsSection pet={pet} editable={userIsOwner} onDelete={deleteTagMutation.mutate} />

      <section className="bg-slate-50 p-8">
        {!userIsOwner && (
          <div className="mb-4">
            <p className="text-2xl leading-loose">Have you seen me?</p>
            <p className="text-lg tracking-wide leading-relaxed text-slate-800">
              If you think you have seen {pet.name}, or you have them in your possession, please let the owner
              know by sending them a message.
            </p>
          </div>
        )}
        {userIsOwner ? <DetailsForm pet={pet} /> : <PublicDetails pet={pet} />}
        {userIsOwner && (
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
        )}
      </section>
    </>
  );
}
