import {useAuth} from "../hooks/useAuth.tsx";
import {Wrapper} from "../components/Wrapper.tsx";

function Home() {
  const { loggedIn, user } = useAuth()

  return (
    <Wrapper>
      <h1>Home</h1>
      {loggedIn && <p>Hi {user?.name}</p>}
    </Wrapper>
  )
}

export default Home;
