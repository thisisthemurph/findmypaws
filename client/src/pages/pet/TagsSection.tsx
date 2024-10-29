import { Pet } from "@/api/types.ts";
import Tag from "@/pages/pet/Tag.tsx";
import NewTagDialog from "@/pages/pet/NewTagDialog.tsx";
import { Badge } from "@/components/ui/badge.tsx";
import { Button } from "@/components/ui/button.tsx";

interface TagsSectionProps {
  pet: Pet;
  editable: boolean;
  onDelete: ({ petId, key }: { petId: string; key: string }) => void;
}

export default function TagsSection({ pet, editable, onDelete }: TagsSectionProps) {
  function handleDeleteTag(petId: string, key: string) {
    if (!editable) return;
    onDelete({ petId, key });
  }

  const tags = Object.entries(pet.tags ?? {});

  return (
    <section className="flex flex-wrap justify-center gap-2 pb-12">
      {tags.map(([key, value]) =>
        editable ? (
          <Tag key={key} identifier={key} handleDelete={() => handleDeleteTag(pet.id, key)}>
            {value}
          </Tag>
        ) : (
          <Badge key={key} title={`tag type: ${key}`} variant="outline" size="lg">
            {value}
          </Badge>
        )
      )}
      {editable && (
        <NewTagDialog pet={pet}>
          {tags.length > 0 ? (
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
          ) : (
            <Button size="sm" variant="outline">
              Add a new tag
            </Button>
          )}
        </NewTagDialog>
      )}
    </section>
  );
}
