import { useState } from "react";
import {
  doCreateUserWithEmailAndPassword,
  doSignInWithGoogle,
} from "../../auth/authFunctions";
import GoogleIcon from "../../assets/google_logo.svg";
import { Box, Typography, Container, TextField, Button } from "@mui/material";
import { useNavigate } from "react-router-dom";

function SignUpPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [passwordConfirm, setPasswordConfirm] = useState<string>("");
  const navigator = useNavigate();

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

      if (!result.user) {
        return;
      }
      navigator("/home");
    } catch (error) {
      alert(error);
    }
  }

  async function handleGoogleSignUp(): Promise<void> {
    try {
      const result = await doSignInWithGoogle();
      if (!result.user) {
        return;
      }
      navigator("/home");
    } catch (error) {
      alert(error);
    }
  }

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        handleSignUp();
      }}
    >
      <Box
        sx={{
          height: "93vh",
          marginTop: "7vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Box
          sx={(theme) => ({
            height: "90%",
            width: "40%",
            backgroundColor: theme.palette.secondary.main,
            border: `1px solid ${theme.palette.divider}`,
            borderRadius: "20px",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            flexDirection: "column",
            margin: "10px",
          })}
        >
          <Typography
            variant="h3"
            sx={(theme) => ({
              fontWeight: 600,
              color: theme.palette.text.secondary,
              marginTop: "3%",
            })}
          >
            SIGN-UP
          </Typography>
          <Container>
            <Container>
              <label htmlFor="email-input">
                <Typography variant="h6">Email:</Typography>
              </label>
              <TextField
                id="email-input"
                label="e.g. name123@mail.com"
                variant="outlined"
                sx={{}}
                fullWidth
                color="primary"
                value={email}
                onChange={(event) => {
                  setEmail(event.target.value);
                }}
              />
            </Container>
            <Container>
              <label htmlFor="password-input">
                <Typography variant="h6">Password:</Typography>
              </label>
              <TextField
                type="password"
                id="password-input"
                label="e.g. first12@#$last"
                variant="outlined"
                sx={{}}
                fullWidth
                color="primary"
                value={password}
                onChange={(event) => {
                  setPassword(event.target.value);
                }}
              />
            </Container>
            <Container>
              <label htmlFor="password-confirm-input">
                <Typography variant="h6">Password:</Typography>
              </label>
              <TextField
                type="password"
                id="password-confirm-input"
                label="e.g. first12@#$last"
                variant="outlined"
                sx={{}}
                fullWidth
                color="primary"
                value={passwordConfirm}
                onChange={(event) => {
                  setPasswordConfirm(event.target.value);
                }}
              />
            </Container>
          </Container>
          <Container sx={{ margin: "10px" }}>
            <Button fullWidth variant="contained" type="submit">
              SIGN-UP
            </Button>
            <Button
              sx={{
                display: "flex",
                flexDirection: "row",
                alignItems: "center",
                justifyContent: "center",
                backgroundColor: "white",
                marginTop: "10px",
              }}
              fullWidth
              type="button"
              onClick={() => {
                handleGoogleSignUp();
              }}
            >
              <img src={GoogleIcon} alt="Google" width={30} height={30} />
              <Typography variant="button" sx={{ margin: "5px" }}>
                SIGN-UP WITH GOOGLE
              </Typography>
            </Button>
          </Container>
        </Box>
      </Box>
    </form>
  );
}

export default SignUpPage;
