import { useState } from "react";
import { doSignInWithEmailAndPassword } from "../../auth/authFunctions";
import { Box, Button, Container, Typography } from "@mui/material";
import GoogleIcon from "../../assets/google_logo.svg";
import { useNavigate } from "react-router-dom";
import { handleGoogleAuth } from "./SignUp";
import XTextField from "../../Components/XTextField";

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
            <XTextField
              required
              labelName="Email:"
              label="e.g. name123@gmail.com"
              value={email}
              onChange={(event) => {
                setEmail(event.target.value);
              }}
            />
            <XTextField
              labelName="Password"
              label="e.g. first123@#$last"
              required
              type="password"
              value={password}
              onChange={(event) => {
                setPassword(event.target.value);
              }}
            />
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
                handleGoogleAuth(navigator);
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
