import { useState } from "react";
import "./pages.css";
import {
  doCreateUserWithEmailAndPassword,
  doSignInWithGoogle,
} from "../../auth/auth";

function SignUpPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [passwordConfirm, setPasswordConfirm] = useState<string>("");

  async function handleSignUp(): Promise<void> {
    if (!email) {
      alert("Email Not Found");
      return;
    }
    if (!password || !passwordConfirm) {
      alert("Password Not Found");
      return;
    }
    if (password != passwordConfirm) {
      alert("Passwords dont match");
      return;
    }
    if (password.length < 6) {
      alert("Too short password");
      return;
    }
    try {
      const result = await doCreateUserWithEmailAndPassword(email, password);
      console.log(result);
    } catch (error) {
      alert(error);
    }
  }

  async function handleGoogleSignUp(): Promise<void> {
    try {
      const result = await doSignInWithGoogle();
      console.log(result);
    } catch (error) {
      alert(error);
    }
  }

  return (
    <form
      className="signup-page"
      onSubmit={(event) => {
        event.preventDefault();
        handleSignUp();
      }}
    >
      <h1>SignUp-Page</h1>
      <div className="signup-page-inputs">
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
        <label htmlFor="confirm-pass">Confirm Password:</label>
        <input
          id="confirm-pass"
          type="password"
          value={passwordConfirm}
          onChange={(event) => {
            setPasswordConfirm(event.target.value);
          }}
        />
      </div>
      <div className="signup-page-options">
        <button type="submit">Log-In</button>
        <button
          onClick={() => {
            handleGoogleSignUp();
          }}
          type="button"
        >
          Login With Google
        </button>
      </div>
    </form>
  );
}

export default SignUpPage;
