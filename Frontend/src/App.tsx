import { useState, type JSX } from "react";
import "./App.css";
import LogInPage from "./Pages/auth/LogIn";
import SignUpPage from "./Pages/auth/SignUp";
import { useAuth } from "./auth/fireBaseContext";
import { doSignOut } from "./auth/auth";

function App() {
  const [page, setPage] = useState<JSX.Element>(<SignUpPage />);
  const auth = useAuth();

  return (
    <div className="home">
      {/* <Navbar setPage={setPage} /> */}
      <div className="auth-button-container">
        <button
          className="auth-button"
          onClick={() => {
            setPage(<LogInPage />);
          }}
        >
          LOGIN
        </button>
        <button
          className="auth-button"
          onClick={() => {
            setPage(<SignUpPage />);
          }}
        >
          SIGN-UP
        </button>
        <button
          className="auth-button"
          onClick={() => {
            console.log(auth);
          }}
        >
          INFO
        </button>
        <button
          className="auth-button"
          onClick={() => {
            doSignOut();
          }}
        >
          LOGOUT
        </button>
      </div>
      <div className="page-display">{page}</div>
    </div>
  );
}

export default App;
