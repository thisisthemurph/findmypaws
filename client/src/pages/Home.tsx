import {useAuth} from "../hooks/useAuth.tsx";

function Home() {
  const { loggedIn, user } = useAuth()

  return (
    <>
      <h1>Home</h1>
      {loggedIn && <p>Hi {user?.name}</p>}
    </>
  )
}

export default Home;
