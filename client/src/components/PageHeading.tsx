import { ReactNode } from "react";

interface PageHeadingProps {
  heading: string;
  subheading?: string;
  children?: ReactNode;
}

export function PageHeading({ heading, subheading, children }: PageHeadingProps) {
  return (
    <section className="flex justify-between items-start mb-6">
      <div>
        <h1 className="mb-0">{heading}</h1>
        {subheading && <p className="text-slate-700 text-sm">{subheading}</p>}
      </div>
      {children && <div>{children}</div>}
    </section>
  );
}
