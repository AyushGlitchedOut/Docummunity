import { AppBar, Box, Button, Toolbar, Typography } from "@mui/material";
import { useAuth } from "../auth/fireBaseContext";
import AppIcon from "../assets/Docummunity.png";
import { useNavigate } from "react-router-dom";
function Navbar() {
  const auth = useAuth();
  var userLoggedIn: boolean = false;
  const navigator = useNavigate();
  if (auth && auth.userLoggedIn) {
    userLoggedIn = true;
  }
  return (
    <AppBar
      position="fixed"
      sx={(theme) => ({
        height: "7vh",
        boxShadow: `2px 2px 2px ${theme.palette.divider}`,
        backgroundColor: theme.palette.background.paper,
        display: "flex",
        alignItems: "center",
        flexDirection: "row",
        justifyContent: "space-between",
        zIndex: theme.zIndex.drawer + 1,
      })}
    >
      <Toolbar>
        <Box
          component="img"
          src={AppIcon}
          sx={(theme) => ({
            height: 40,
            marginRight: 2,
            borderRadius: "100%",
            border: `1px solid ${theme.palette.divider}`,
          })}
        />
        <Typography
          variant="h4"
          sx={(theme) => ({
            letterSpacing: 4,
            fontWeight: 500,
            color: theme.palette.text.primary,
            cursor: "pointer",
          })}
          onClick={() => {
            navigator("/");
          }}
        >
          DOCUMMUNITY
        </Typography>
      </Toolbar>
      {userLoggedIn ? <LoggedInToolBar /> : <LoggedOutToolBar />}
    </AppBar>
  );
}

function LoggedInToolBar() {
  const navigator = useNavigate();
  return (
    <Toolbar sx={{ width: "30%", display: "flex", justifyContent: "right" }}>
      <Button
        sx={(theme) => ({
          cursor: "pointer",
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.text.primary,
          fontSize: "clamp(50%, 100%, 150%)",
          width: "30%",
          margin: 2,
        })}
        onClick={() => {
          navigator("/download");
        }}
      >
        Download
      </Button>
      <Button
        sx={(theme) => ({
          cursor: "pointer",
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.text.primary,
          fontSize: "clamp(50%, 100%, 150%)",
          width: "30%",
          margin: 2,
        })}
        onClick={() => {
          navigator("/about");
        }}
      >
        About
      </Button>
    </Toolbar>
  );
}
function LoggedOutToolBar() {
  const navigator = useNavigate();
  return (
    <Toolbar sx={{ width: "40%", display: "flex", justifyContent: "right" }}>
      <Button
        sx={(theme) => ({
          cursor: "pointer",
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.text.primary,
          fontSize: "clamp(50%, 100%, 150%)",
          width: "30%",
          margin: 2,
        })}
        onClick={() => {
          navigator("/download");
        }}
      >
        DOWNLOAD
      </Button>
      <Button
        sx={(theme) => ({
          cursor: "pointer",
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.text.primary,
          fontSize: "clamp(50%, 100%, 150%)",
          width: "30%",
          margin: 2,
        })}
        onClick={() => {
          navigator("/login");
        }}
      >
        LOG-IN
      </Button>
      <Button
        sx={(theme) => ({
          cursor: "pointer",
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.text.primary,
          fontSize: "clamp(50%, 100%, 150%)",
          width: "30%",
          margin: 2,
        })}
        onClick={() => {
          navigator("/signUp");
        }}
      >
        SIGN-UP
      </Button>
      <Button
        sx={(theme) => ({
          cursor: "pointer",
          backgroundColor: theme.palette.primary.main,
          color: theme.palette.text.primary,
          fontSize: "clamp(50%, 100%, 150%)",
          width: "30%",
          margin: 2,
        })}
        onClick={() => {
          navigator("about");
        }}
      >
        ABOUT US
      </Button>
    </Toolbar>
  );
}

export default Navbar;
