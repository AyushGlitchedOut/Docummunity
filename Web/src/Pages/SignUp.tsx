import {
  GitHub,
  Google,
  Microsoft,
  Visibility,
  VisibilityOff,
} from "@mui/icons-material";
import {
  Box,
  Button,
  Container,
  Fade,
  FormControlLabel,
  FormHelperText,
  IconButton,
  InputAdornment,
  Switch,
  TextField,
  Typography,
} from "@mui/material";
import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useFirebase } from "../services/firebase";

function SignUpPage() {
  const [lightTheme, setLightTheme] = useState(true);
  const [capsLockOn, setCapsLockOn] = useState<boolean>();
  const [passwordShow, setpasswordShow] = useState<boolean>(false);
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [confirmedPassword, setconfirmedPassword] = useState<string>("");
  const firebase = useFirebase();
  const navigator = useNavigate();

  function handleCapsLock(event: any) {
    setCapsLockOn(event.getModifierState("CapsLock"));
  }

  async function handleSignUp(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (password == confirmedPassword) {
      if (password.length < 6) {
        alert("Password must be atleats 6 charcters long!");
        return;
      }
      try {
        await firebase.SignUpWithEmailAndPassword(email, password);
      } catch (error: any) {
        if (error.code == "auth/email-already-in-use") {
          alert("Email Already in Use! (Try Signing-In)");
        }
        return;
      }
    } else {
      alert("Please match the password and confirmed password fields");
    }
  }

  async function handleSignInWithGoogle() {
    try {
      await firebase.LogInWithGoogleAccount();
    } catch (error: any) {
      alert("Something Went Wrong!");
      return;
    }
  }
  async function handleSignInWithGithub() {
    try {
      await firebase.LogInWithGithubAccount();
    } catch (error: any) {
      alert("Something Went Wrong!");
      return;
    }
  }

  useEffect(() => {
    if (firebase.isLoggedIn) {
      navigator("/home");
    }
  }, [firebase]);

  return (
    <>
      <FormControlLabel
        label={
          <Typography variant="h6" sx={{ color: "text.primary" }}>
            Light Mode
          </Typography>
        }
        sx={{ position: "absolute", margin: "5px" }}
        control={
          <Switch
            checked={lightTheme}
            onChange={() => {
              setLightTheme(!lightTheme);
            }}
            size="medium"
            color="info"
          />
        }
      />
      <Box
        sx={{
          height: "91vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <Container
          sx={{
            bgcolor: "red",
            height: "80vh",
            display: "flex",
            flexDirection: "column",
            position: "absolute",
            left: "47%",
            top: "53%",
            width: "35vw",
            transform: "translate(-50%,-50%)",
            borderRadius: "10px",
            backgroundColor: "background.paper",
            boxShadow: "5px 3px 1px grey, inset 3px 3px 10px grey",
            textAlign: "center",
            padding: "1%",
          }}
        >
          <Typography
            variant="h4"
            sx={{
              width: "100%",
              color: "text.secondary",
              letterSpacing: "2px",
              fontWeight: "800",
            }}
          >
            SIGN-UP
          </Typography>
          <Box
            component="form"
            onSubmit={(event) => handleSignUp(event)}
            sx={{
              display: "flex",
              flexDirection: "column",
              justifyContent: "space-between",
              alignItems: "start",
              width: "100%",
              height: "80%",
              margin: "10px",
            }}
          >
            {/* Email-Id */}
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "start",
                width: "100%",
              }}
            >
              <Typography variant="subtitle2">Email-id:</Typography>
              <TextField
                label="Email"
                variant="filled"
                fullWidth={true}
                required={true}
                value={email}
                onChange={(event) => setEmail(event.target.value)}
              />
            </Box>
            {/* Password */}
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "start",
                width: "100%",
              }}
            >
              <Typography variant="subtitle2">Password:</Typography>
              <TextField
                label="Password"
                variant="filled"
                type={passwordShow ? undefined : "password"}
                fullWidth={true}
                value={password}
                required={true}
                onChange={(event) => setPassword(event.target.value)}
                onKeyDown={handleCapsLock}
                onKeyUp={handleCapsLock}
                slotProps={{
                  input: {
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          onClick={() => setpasswordShow(!passwordShow)}
                        >
                          {passwordShow ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  },
                }}
              />
              {/* Note: kept the fade and helper outside the input tag without using the helper text prop as it will lead to nested <p> tags which give hydration error */}
              <Fade in={capsLockOn} timeout={500}>
                <FormHelperText component={"div"} error>
                  CAPS LOCK IS ON!!
                </FormHelperText>
              </Fade>
            </Box>
            {/* Confirm Password */}
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "start",
                width: "100%",
              }}
            >
              <Typography variant="subtitle2">Confirm Password:</Typography>
              <TextField
                required={true}
                label="Re-enter Password"
                variant="filled"
                type={passwordShow ? undefined : "password"}
                fullWidth={true}
                value={confirmedPassword}
                onChange={(event) => {
                  setconfirmedPassword(event.target.value);
                }}
                onKeyDown={handleCapsLock}
                onKeyUp={handleCapsLock}
                slotProps={{
                  input: {
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          onClick={() => setpasswordShow(!passwordShow)}
                        >
                          {passwordShow ? <VisibilityOff /> : <Visibility />}
                        </IconButton>
                      </InputAdornment>
                    ),
                  },
                }}
              />
              {/* Note: kept the fade and helper outside the input tag without using the helper text prop as it will lead to nested <p> tags which give hydration error */}
              <Fade in={capsLockOn} timeout={500}>
                <FormHelperText error>CAPS LOCK IS ON!!</FormHelperText>
              </Fade>
            </Box>
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                width: "100%",
              }}
            >
              <Button
                variant="contained"
                sx={{ width: "40%", margin: "1%" }}
                type="submit"
              >
                <Typography variant="h6">Sign-up</Typography>
              </Button>
              <Link to="/log_in">Already have an account? Sign-In!</Link>
            </Box>
            {/* Other Sign-In Options */}

            <Box
              sx={{
                width: "100%",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <IconButton
                size="large"
                onClick={() => {
                  handleSignInWithGoogle();
                }}
              >
                <Google fontSize="large" />
              </IconButton>
              <IconButton
                size="large"
                onClick={() => {
                  handleSignInWithGithub();
                }}
              >
                <GitHub fontSize="large" />
              </IconButton>
            </Box>
          </Box>
        </Container>
      </Box>
    </>
  );
}

export default SignUpPage;
