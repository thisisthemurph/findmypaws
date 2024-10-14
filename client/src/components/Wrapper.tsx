import {ReactNode} from "react";

export const Wrapper = ({children}: { children: ReactNode }) => {
  return (
    <section className="p-4">
      { children }
    </section>
  )
}
