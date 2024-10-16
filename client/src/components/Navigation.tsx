import {Link, useNavigate} from "react-router-dom";
import {useAuth} from "../hooks/useAuth.tsx";

function Navigation() {
  const { loggedIn, logout } = useAuth();
  const navigate = useNavigate();

  function handleLogout() {
    logout().then(() => {
      navigate("/");
    }).catch((e) => {
      console.error(e);
      alert("Could not log out");
    });
  }

  return (
    <nav className="p-4 bg-purple-300 shadow mb-4">
      <h1 className="text-2xl">Find my paws</h1>
      <ul className="flex gap-4">
        <li><Link to="/">Home</Link></li>
        <li><Link to="/profile">Pet Profile</Link></li>
        {!loggedIn && <li><Link to="/login">Log in</Link></li>}
        {!loggedIn && <li><Link to="/signup">Sign up</Link></li>}
        {loggedIn && <button onClick={handleLogout}>Log out</button>}
      </ul>
    </nav>
  )
}

export default Navigation;
