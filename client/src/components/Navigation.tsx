import {Link} from "react-router-dom";

function Navigation() {
  return (
    <nav className="p-4 bg-purple-300 shadow mb-4">
      <h1 className="text-2xl">Find my paws</h1>
      <ul className="flex gap-4">
        <li><Link to="/">Home</Link></li>
        <li><Link to="/profile">Pet Profile</Link></li>
      </ul>
    </nav>
  )
}

export default Navigation;
