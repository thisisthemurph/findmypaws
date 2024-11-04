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

export type Notification = {
  id: string;
  type: "spotted_pet";
  message: string;
  link: string;
  seen: boolean;
  created_at: string;
};
