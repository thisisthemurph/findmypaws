import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth.tsx";
import { ReactNode } from "react";
import { useNavigation } from "../hooks/useNavigation.tsx";

function Navigation() {
  const [navOpen, setNavOpen] = useNavigation();
  const { logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  const loggedIn = isAuthenticated();

  async function handleLogout() {
    try {
      await logout();
    } finally {
      setNavOpen(false);
      navigate("/login");
    }
  }

  return (
    <nav className="flex flex-col md:p-4 md:flex-row md:justify-between md:items-center shadow mb-4">
      <section className="p-4 md:p-0 flex justify-between">
        <h1 className="text-2xl text-purple-700">Findmypaws</h1>
        <button className="md:hidden" onClick={() => setNavOpen(!navOpen)}>
          {navOpen ? "close" : "open"}
        </button>
      </section>
      <ul
        className={`${navOpen ? "block" : "hidden"} md:flex px-4 md:px-0 pb-4 md:pb-0 flex flex-col items-center md:flex-row`}
      >
        <NavLink to="/">Home</NavLink>
        {loggedIn && <NavLink to="/profile">Profile</NavLink>}
        {!loggedIn && <NavLink to="/login">Log in</NavLink>}
        {!loggedIn && <NavLink to="/signup">Sign up</NavLink>}
        {loggedIn && (
          <li className="w-full md:w-auto">
            <button
              onClick={handleLogout}
              className="block w-full p-4 text-center text-2xl md:text-xl rounded-md hover:bg-purple-200 hover:text-purple-700 transition-colors"
            >
              Log out
            </button>
          </li>
        )}
      </ul>
    </nav>
  );
}

function NavLink({ children, to }: { children: ReactNode; to: string }) {
  const [, setNavOpen] = useNavigation();
  return (
    <li className="w-full md:w-auto">
      <Link
        onClick={() => setNavOpen(false)}
        className="inline-block w-full p-4 text-center text-2xl md:text-xl rounded-md hover:bg-purple-200 hover:text-purple-700 transition-colors"
        to={to}
      >
        {children}
      </Link>
    </li>
  );
}

export default Navigation;
