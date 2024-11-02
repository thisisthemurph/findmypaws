export type PetTags = {
  [key: string]: string; // A dictionary of tags where both keys and values are strings
};

export type Pet = {
  id: string;
  user_id: string;
  type: "Cat" | "Dog" | "Unspecified";
  name: string;
  tags: PetTags;
  dob: string | null;
  avatar: string | null;
  blurb: string | null;
  created_at: string;
  updated_at: string;
};

export type Alert = {
  id: number;
  pet_id: string;
  user_id: string;
  anonymous_user_id: string;
  created_at: string;
};

type SpottedPetLinks = {
  pet: string;
  [key: string]: string;
};

export type Notification = {
  id: string;
  type: "spottedPet";
  message: string;
  seen: boolean;
  created_at: string;
} & { type: "spottedPet"; links: SpottedPetLinks };
