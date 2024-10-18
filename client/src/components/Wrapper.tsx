import { HTMLAttributes } from "react";
import { cn } from "@/lib/utils.ts";

export const Wrapper = ({ children, className, ...props }: HTMLAttributes<HTMLDivElement>) => {
  return (
    <div className={cn("p-4", className)} {...props}>
      {children}
    </div>
  );
};
