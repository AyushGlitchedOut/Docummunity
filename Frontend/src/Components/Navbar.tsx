import { type JSX } from "react";
import "./Navbar.css";
import CreatePage from "../Pages/CRUD/Create";
import ReadPage from "../Pages/CRUD/Read";
import UpdatePage from "../Pages/CRUD/Update";
import DeletePage from "../Pages/CRUD/Delete";

function Navbar(args: {
  setPage: React.Dispatch<React.SetStateAction<JSX.Element>>;
}) {
  return (
    <div className="navbar">
      <button
        className="navbar-button"
        onClick={() => args.setPage(<CreatePage />)}
      >
        CREATE
      </button>

      <button
        className="navbar-button"
        onClick={() => args.setPage(<ReadPage />)}
      >
        READ
      </button>

      <button
        className="navbar-button"
        onClick={() => args.setPage(<UpdatePage />)}
      >
        UPDATE
      </button>

      <button
        className="navbar-button"
        onClick={() => args.setPage(<DeletePage />)}
      >
        DELETE
      </button>
    </div>
  );
}

export default Navbar;
