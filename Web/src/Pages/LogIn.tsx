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
import { useState } from "react";
import { Link } from "react-router-dom";

function LoginPage() {
  const [lightTheme, setLightTheme] = useState(true);
  const [capsLockOn, setCapsLockOn] = useState<boolean>();
  const [passwordShow, setpasswordShow] = useState<boolean>(false);
  function handleCapsLock(event: any) {
    setCapsLockOn(event.getModifierState("CapsLock"));
  }

  function handleSignIn() {
    alert("Sign In");
  }

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
            LOG-IN
          </Typography>
          <Box
            component="form"
            sx={{
              display: "flex",
              flexDirection: "column",
              justifyContent: "space-between",
              alignItems: "start",
              width: "100%",
              height: "70%",
              margin: "10px",
            }}
          >
            {/* Username */}
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                alignItems: "start",
                width: "100%",
              }}
            >
              <Typography variant="subtitle2">Username/Email-id:</Typography>
              <TextField label="Username" variant="filled" fullWidth={true} />
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
                id="capsLockCheck"
                fullWidth={true}
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
                helperText={
                  <Fade in={capsLockOn} timeout={500}>
                    <FormHelperText error>CAPS LOCK IS ON!!</FormHelperText>
                  </Fade>
                }
              />
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
                onClick={() => handleSignIn()}
              >
                <Typography variant="h6">Sign-In</Typography>
              </Button>
              <Link to="/sign_up">Don't have an Account? Sign up!</Link>
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
              <IconButton size="large">
                <Google fontSize="large" />
              </IconButton>
              <IconButton size="large">
                <GitHub fontSize="large" />
              </IconButton>
              <IconButton size="large">
                <Microsoft fontSize="large" />
              </IconButton>
            </Box>
          </Box>
        </Container>
      </Box>
    </>
  );
}

export default LoginPage;
