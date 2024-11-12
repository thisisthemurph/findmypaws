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
import { useCallback, useEffect } from "react";
import { Dialog, DialogContent, DialogTitle, DialogTrigger } from "@/components/ui/dialog.tsx";
import { v4 as uuidv4 } from "uuid";
import AnonymousUserForm from "@/pages/pet/AnonymousUserForm.tsx";

function getOrCreateAnonymousUserId(): string {
  let userId = localStorage.getItem("anonymousUserId");
  if (!userId) {
    userId = uuidv4();
    localStorage.setItem("anonymousUserId", userId);
  }
  return userId;
}

interface AlertRequest {
  user_id?: string;
  anonymous_user_id?: string;
}

interface AlertResponse {
  alert_created: boolean;
}

export default function PetPage() {
  const { id } = useParams();
  const { toast } = useToast();
  const navigate = useNavigate();
  const { userId } = useAuth();
  const api = useApi();
  const queryClient = useQueryClient();

  const { data: pet, isLoading } = useQuery<Pet>({
    queryKey: ["pet"],
    queryFn: () => api<Pet>(`/pets/${id}`),
  });

  const userIsOwner = userId === pet?.user_id;

  const sendAlert = useCallback(async (petId: string, userId: string | null | undefined) => {
    const body: AlertRequest = {};
    if (userId) {
      body.user_id = userId;
    } else {
      body.anonymous_user_id = getOrCreateAnonymousUserId();
    }
    return await api<AlertResponse>(`/pets/${petId}/alert`, {
      method: "POST",
      body: JSON.stringify(body),
    });
  }, []);

  useEffect(() => {
    if (!pet?.id || !pet.user_id || userId === pet.user_id) return;
    sendAlert(pet.id, userId)
      .then((resp) => {
        if (resp.alert_created) {
          toast({
            title: "Alert",
            description: "An alert has been sent to the owner to inform them you have visited this page.",
          });
        }
      })
      .catch((err: Error) => console.error(err));
  }, [userId, pet?.id, pet?.user_id, sendAlert]);

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

      {!userIsOwner && (
        <Dialog>
          <DialogTrigger asChild>
            <button className="absolute bottom-5 right-5 size-16 flex justify-center items-center bg-purple-300/50 text-slate-800 rounded-full hover:bg-purple-300 transition-colors">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="h-8 w-8"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 0 1-.825-.242m9.345-8.334a2.126 2.126 0 0 0-.476-.095 48.64 48.64 0 0 0-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0 0 11.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155"
                />
              </svg>
            </button>
          </DialogTrigger>
          <DialogContent className="w-[95%] rounded-lg">
            <DialogTitle>Start a chat for {pet.name}</DialogTitle>
            <p>Tell us your name to start the conversation...</p>
            <AnonymousUserForm conversationIdentifier={pet.id} />
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}
