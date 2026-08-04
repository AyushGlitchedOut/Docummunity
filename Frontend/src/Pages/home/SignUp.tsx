import { useMemo, useRef, useState } from "react";
import {
  doCreateUserWithEmailAndPassword,
  doSignInWithGoogle,
} from "../../auth/authFunctions";
import GoogleIcon from "../../assets/google_logo.svg";
import {
  Box,
  Typography,
  Container,
  TextField,
  Button,
  Divider,
  Avatar,
} from "@mui/material";
import { useNavigate, type NavigateFunction } from "react-router-dom";
import { auth } from "../../auth/fireBaseConfig";
import { BACKEND_URL } from "../../consts";

// A common google sign up function to handle both sign-up and log-in cases, used in both login and sign-up page
export async function handleGoogleAuth(
  navigator: NavigateFunction,
): Promise<void> {
  //Navigator passed as an argument rather than a new navigator since navigator is a hook and  cant be called outside of component
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

function SignUpPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [passwordConfirm, setPasswordConfirm] = useState<string>("");
  const [displayName, setDisplayName] = useState<string>("");
  const [bio, setBio] = useState<string>("");
  const [profileImageFile, setProfileImageFile] = useState<File | null>(null);

  const profilePictureURL = useMemo(() => {
    return profileImageFile ? URL.createObjectURL(profileImageFile) : null;
  }, [profileImageFile]);

  const hiddenFilePickerRef = useRef<HTMLInputElement>(null);

  const navigator = useNavigate();

  async function handleSignUp(): Promise<void> {
    if (password != passwordConfirm) {
      alert("Passwords dont match");
      return;
    }
    if (password.length < 6) {
      alert("Too short password");
      return;
    }
    if (!displayName) {
      alert("Please Enter a Display Name");
      return;
    }
    if (profileImageFile) {
      if (profileImageFile.size > 2 * 1000 * 1000) {
        alert("The file should be smaller than 2mb");
        return;
      }
    }

    try {
      const result = await doCreateUserWithEmailAndPassword(email, password);
      if (!result.user) {
        return;
      }

      //Creating the Form
      const form = new FormData();
      if (profileImageFile) {
        form.append("PROFILE_PIC", profileImageFile);
      }
      form.append("DISPLAY_NAME", displayName);
      form.append("BIO", bio);
      form.append("SETTINGS", "{}");

      const token = await auth.currentUser?.getIdToken();
      if (!token) {
        alert("Auth Issues");
        return;
      }

      const createAccount = await fetch(
        //DO NOT PUT /CREATE/ or else the address wont work
        BACKEND_URL + "/user/CREATE",
        {
          method: "POST",
          headers: { Authorization: `Bearer ${token}` },
          body: form,
        },
      );
      if (!createAccount.ok) {
        alert("Something Went Wrong");
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
            width: "80%",
            backgroundColor: theme.palette.secondary.main,
            border: `1px solid ${theme.palette.divider}`,
            borderRadius: "20px",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            flexDirection: "row",
            margin: "10px",
          })}
        >
          <Box
            sx={(theme) => ({
              width: "40%",
              margin: "1%",
              height: "90%",
              display: "flex",
              alignItems: "center",
              justifyContent: "space-evenly",
              flexDirection: "column",
            })}
          >
            <Container
              sx={(theme) => ({
                width: "60%",
                aspectRatio: 1,
                display: "flex",
                alignItems: "center",
                justifyContent: "space-between",
                flexDirection: "column",
              })}
            >
              <Avatar
                src={profilePictureURL ?? undefined}
                sx={(theme) => ({
                  height: 200,
                  width: 200,
                  aspectRatio: 1,
                  border: "2px solid black",
                })}
                onClick={() => hiddenFilePickerRef.current?.click()}
              />
              <Button
                variant="contained"
                sx={(theme) => ({ margin: "5%" })}
                onClick={() => hiddenFilePickerRef.current?.click()}
              >
                Change Profile Pic
              </Button>

              {/* Hidden Input for Allowing File selection */}
              <input
                hidden
                ref={hiddenFilePickerRef}
                type="file"
                accept="image/*"
                onChange={(event) => {
                  const selectedFile = event.target.files?.[0];

                  if (selectedFile) {
                    setProfileImageFile(selectedFile);
                  }
                }}
              />
            </Container>
            <Container>
              <label htmlFor="display_name">
                <Typography variant="h6">Display Name:</Typography>
              </label>
              <TextField
                required
                id="display_name"
                label="e.g. Ayush Gupta"
                variant="outlined"
                sx={{}}
                fullWidth
                color="primary"
                value={displayName}
                onChange={(event) => {
                  setDisplayName(event.target.value);
                }}
              />
            </Container>
            <Container>
              <label htmlFor="bio">
                <Typography variant="h6">Bio:</Typography>
              </label>
              <TextField
                id="bio"
                label="Tell people about yourself (Optional)"
                variant="outlined"
                sx={{}}
                fullWidth
                color="primary"
                multiline
                minRows={2}
                maxRows={4}
                value={bio}
                onChange={(event) => {
                  setBio(event.target.value);
                }}
              />
            </Container>
          </Box>
          <Divider
            sx={(theme) => ({
              backgroundColor: theme.palette.divider,
              width: "0.2%",
              height: "90%",
            })}
          />

          <Box
            sx={(theme) => ({
              height: "90%",
              width: "50%",
              margin: "5%",
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              flexDirection: "column",
            })}
          >
            <Typography
              variant="h3"
              align="center"
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
                  required
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
                  required
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
                  <Typography variant="h6">Repeat Password:</Typography>
                </label>
                <TextField
                  required
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
                  handleGoogleAuth(navigator);
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
      </Box>
    </form>
  );
}

export default SignUpPage;
