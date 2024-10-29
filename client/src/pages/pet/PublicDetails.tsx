import { Pet } from "@/api/types.ts";

interface PublicDetailsProps {
  pet: Pet;
}

export default function PublicDetails({ pet }: PublicDetailsProps) {
  return (
    <section className="">
      {pet.blurb && (
        <div>
          <p className="text-2xl leading-loose">A little about {pet.name}</p>
          <p className="text-lg tracking-wide leading-relaxed text-slate-800">{pet.blurb}</p>
        </div>
      )}
    </section>
  );
}
