import { useState } from "react";
import "./pages.css";
import {
  doSignInWithGoogle,
  doSignInWithEmailAndPassword,
} from "../../auth/authFunctions";
import { Box, Button, Container, TextField, Typography } from "@mui/material";
import GoogleIcon from "../../assets/google_logo.svg";
import { useNavigate } from "react-router-dom";

function LogInPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const navigator = useNavigate();

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
      if (!result.user) {
        return;
      }
      navigator("/home");
    } catch (error) {
      alert(error);
    }
  }

  async function handleGoogleLogin(): Promise<void> {
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
        handleLogin();
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
            LOG-IN
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
          </Container>
          <Container sx={{ margin: "10px" }}>
            <Button fullWidth variant="contained" type="submit">
              LOG-IN
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
                handleGoogleLogin();
              }}
            >
              <img src={GoogleIcon} alt="Google" width={30} height={30} />
              <Typography variant="button" sx={{ margin: "5px" }}>
                SIGN-IN WITH GOOGLE
              </Typography>
            </Button>
          </Container>
        </Box>
      </Box>
    </form>
  );
}

export default LogInPage;
