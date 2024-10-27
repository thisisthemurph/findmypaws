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

export default function Header() {
  return (
    <Sheet>
      <header className="flex justify-between items-center p-4 w-full">
        <SignedIn>
          <UserButton />
        </SignedIn>
        <Link to="/" className="mb-0 text-lg text-slate-700 hover:text-blue-600">
          findmypaws
        </Link>
        <SheetTrigger asChild>
          <Button variant="ghost" size="sm">
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
      <SheetContent side="top">
        <SheetHeader>
          <div className="flex">
            <SignedIn>
              <UserButton />
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
          <SheetDescription>This is the navigation.</SheetDescription>
        </SheetHeader>
        <nav className="flex flex-col justify-center">
          <NavLink to="/">Home</NavLink>
          <SignedIn>
            <NavLink to="/dashboard">Dashboard</NavLink>
          </SignedIn>
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