import { Wrapper } from "../components/Wrapper.tsx";
import {useEffect, useState} from "react";

type Tag = {
  id: string;
  key: string;
  value: string;
}

type Pet = {
  id: string;
  name: string;
  image: string;
  tags: Tag[];
  missing?: boolean;
}

const pets: Pet[] = [
  {
    id: "1",
    name: "doris kiff",
    image: "/images/Doris.png",
    missing: false,
    tags: [
      {
        id: "1",
        key: "breed",
        value: "Staffy mastif",
      },
      {
        id: "2",
        key: "age",
        value: "5 years old",
      },
      {
        id: "3",
        key: "temperament",
        value: "Playful",
      }
    ],
  }
]

const emptyPet: Pet = {
  id: "",
  name: "",
  image: "",
  tags: [],
}

function PetProfile() {
  const [pet, setPet] = useState<Pet>(emptyPet);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    const pet = pets.find((pet) => pet.id === "1");
    if (!pet) {
      return;
    }

    setPet(pet);
    setIsLoading(false);
  }, [])

  return (
    <Wrapper>
      {
        isLoading
          ? (<p className="text-2xl text-center animate-pulse">Walking...</p>)
          : (
            <section className="flex flex-col gap-2">
              {pet?.missing && (<p className="text-red-600 font-bold text-2xl text-center mb-4">MISSING</p>)}
              <section className="flex flex-col gap-4 items-center justify-center">
                <img
                  className="shadow-2xl rounded-full"
                  src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTCHo3CkaH0oRY3MvrEN0xgn-x_Lsn3Lm3lVQ&s" alt={pet.name} />
                <p className="text-4xl capitalize font-semibold text-center">{pet.name}</p>
              </section>
              <section className="flex justify-center gap-2">
                {pet.tags.map((t) => (<p key={t.id} title={t.key} className="inline text-slate-800 text-sm bg-slate-200 rounded-full shadow px-2 py-1">{t.value}</p>))}
              </section>
            </section>
          )
      }
    </Wrapper>
  )
}

export default PetProfile;
