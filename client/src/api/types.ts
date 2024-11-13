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

export type Message = {
  id: number;
  conversationId: number;
  senderId: number;
  recipientId: number;
  text: string;
  createdAt: string;
  readAt: string;
  outgoing: boolean;
};

type ConversationPetDetail = {
  name: string;
  type: string;
};

type ConversationParticipant = {
  id: string;
  name: string;
};

export type Conversation = {
  id: number;
  identifier: string;
  pet: ConversationPetDetail;
  primaryParticipantId: string;
  secondaryParticipantId: string;
  participant: ConversationParticipant;
  otherParticipant: ConversationParticipant;
  lastMessageAt: string | null;
  createdAt: string;
  title: string;
};

export interface AnonymousUser {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}
