import { Pet } from "@/api/types.ts";
import { ChangeEvent, useState } from "react";
import { Dialog, DialogClose, DialogContent, DialogTitle, DialogTrigger } from "@/components/ui/dialog.tsx";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar.tsx";
import { Label } from "@/components/ui/label.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Button } from "@/components/ui/button.tsx";

function PetAvatar({ pet, changeAvatar }: { pet: Pet; changeAvatar: (file: File) => Promise<void> }) {
  const [file, setFile] = useState<File | undefined>();

  function handleFileChange(event: ChangeEvent<HTMLInputElement>) {
    event.preventDefault();
    if (!event.target.files) {
      setFile(undefined);
      return;
    }
    setFile(event.target.files[0]);
  }

  return (
    <section className="group relative flex flex-col gap-4 items-center justify-center mb-4">
      <Dialog>
        <DialogTrigger>
          <Avatar className="relative w-64 h-64 hover:shadow-lg">
            <AvatarImage src={`${import.meta.env.VITE_BASE_URL}/${pet.avatar}`} />
            <AvatarFallback className="flex flex-col gap-4">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth="1.5"
                stroke="currentColor"
                className="w-12 h-12"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m2.25 15.75 5.159-5.159a2.25 2.25 0 0 1 3.182 0l5.159 5.159m-1.5-1.5 1.409-1.409a2.25 2.25 0 0 1 3.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 0 0 1.5-1.5V6a1.5 1.5 0 0 0-1.5-1.5H3.75A1.5 1.5 0 0 0 2.25 6v12a1.5 1.5 0 0 0 1.5 1.5Zm10.5-11.25h.008v.008h-.008V8.25Zm.375 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Z"
                />
              </svg>
              <p className="text-slate-800">Upload an avatar for {pet.name}</p>
            </AvatarFallback>
          </Avatar>
        </DialogTrigger>
        <DialogContent className="w-[95%] rounded-lg">
          <DialogTitle>Upload an avatar for {pet.name}</DialogTitle>
          <form
            onSubmit={async (event) => {
              event.preventDefault();
              if (!file) {
                return;
              }
              await changeAvatar(file);
            }}
          >
            <Label htmlFor="avatar" className="hidden">
              Avatar
            </Label>
            <div className="flex gap-2">
              <Input
                id="avatar"
                name="avatar"
                onChange={handleFileChange}
                placeholder={`${pet.name}.jpg`}
                type="file"
              />
              <DialogClose type="submit" asChild>
                <Button type="submit" disabled={file === undefined}>
                  Upload
                </Button>
              </DialogClose>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <p className="text-4xl capitalize font-semibold text-center">{pet.name}</p>
    </section>
  );
}

export default PetAvatar;
