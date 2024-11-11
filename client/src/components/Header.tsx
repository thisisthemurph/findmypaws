import { Notification } from "@/api/types.ts";
import { Link } from "react-router-dom";
import { SignedIn, SignedOut, UserButton } from "@clerk/clerk-react";
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet.tsx";
import { Button } from "@/components/ui/button.tsx";
import { ReactNode } from "react";
import * as VisuallyHidden from "@radix-ui/react-visually-hidden";
import { useApi } from "@/hooks/useApi.ts";
import NotificationMenu from "@/components/NotificationMenu.tsx";
import { useQuery } from "@tanstack/react-query";

export default function Header() {
  const api = useApi();

  const { data } = useQuery<Notification[]>({
    queryKey: ["notifications"],
    queryFn: async () => await api<Notification[]>("/user/notifications"),
  });

  return (
    <Sheet>
      <header className="flex justify-between items-center h-[5rem] px-4 w-full">
        <SignedIn>
          <NotificationMenu notifications={data ?? []} />
        </SignedIn>
        <Link to="/" className="mb-0 text-lg text-slate-700 hover:text-blue-600">
          findmypaws
        </Link>
        <SheetTrigger asChild>
          <Button variant="ghost" size="icon" className="rounded-full">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="size-6"
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 9h16.5m-16.5 6.75h16.5" />
            </svg>
          </Button>
        </SheetTrigger>
      </header>
      {/*https://discord.com/channels/856971667393609759/1300053615154565151*/}
      <SheetContent side="top">
        <SheetHeader>
          <div className="flex">
            <SignedIn>
              <UserButton
                appearance={{
                  elements: {
                    userButtonPopoverCard: "pointer-events-auto",
                  },
                }}
              />
            </SignedIn>
            <SignedOut>
              <SheetClose asChild>
                <Button variant="outline" size="sm" asChild>
                  <Link to="/sign-in">Sign in</Link>
                </Button>
              </SheetClose>
            </SignedOut>
          </div>
          <SheetTitle>Findmypaws</SheetTitle>
          <SheetDescription>
            <VisuallyHidden.Root>This is the navigation</VisuallyHidden.Root>
          </SheetDescription>
        </SheetHeader>
        <nav className="flex flex-col justify-center">
          <NavLink to="/">Home</NavLink>
          <SignedIn>
            <NavLink to="/dashboard">Your Kennel</NavLink>
          </SignedIn>
          <NavLink to="/conversations">Chats</NavLink>
        </nav>
      </SheetContent>
    </Sheet>
  );
}

interface NavLinkProps {
  children: ReactNode;
  to: string;
}

function NavLink({ to, children }: NavLinkProps) {
  return (
    <SheetClose asChild>
      <Link to={to} className="py-2 text-center text-lg">
        {children}
      </Link>
    </SheetClose>
  );
}
