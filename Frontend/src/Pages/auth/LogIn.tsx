import { useState } from "react";
import "./pages.css";
import {
  doSignInWithGoogle,
  doSignInWithEmailAndPassword,
} from "../../auth/auth";

function LogInPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  async function handleLogin(): Promise<void> {
    if (!email) {
      alert("Email Not Found");
      return;
    }
    if (!password) {
      alert("Password Not found");
      return;
    }
    try {
      const result = await doSignInWithEmailAndPassword(email, password);
      console.log(result);
    } catch (error) {
      alert(error);
    }
  }

  async function handleGoogleLogin(): Promise<void> {
    try {
      const result = await doSignInWithGoogle();
      console.log(result);
    } catch (error) {
      alert(error);
    }
  }

  return (
    <form
      className="login-page"
      onSubmit={(event) => {
        event.preventDefault();
        handleLogin();
      }}
    >
      <h1>Login-Page</h1>
      <div className="login-page-inputs">
        <label htmlFor="email">Enter Email-Id:</label>
        <input
          id="email"
          value={email}
          type="email"
          onChange={(event) => {
            setEmail(event.target.value);
            event.target.reportValidity();
          }}
        />
        <label htmlFor="pass">Enter Password:</label>
        <input
          id="pass"
          type="password"
          value={password}
          onChange={(event) => {
            setPassword(event.target.value);
          }}
        />
      </div>
      <div className="login-page-options">
        <button type="submit">Log-In</button>
        <button
          onClick={() => {
            handleGoogleLogin();
          }}
          type="button"
        >
          Login With Google
        </button>
      </div>
    </form>
  );
}

export default LogInPage;
