import {
  AppBar,
  Avatar,
  Box,
  Button,
  Toolbar,
  Typography,
} from "@mui/material";
import { useAuth } from "../auth/fireBaseContext";
import AppIcon from "../assets/Docummunity.png";
import { useNavigate } from "react-router-dom";
import type { User } from "firebase/auth";
import { doSignOut } from "../auth/auth";
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
      {userLoggedIn ? (
        <LoggedInToolBar user={auth ? auth.currentUser : null} />
      ) : (
        <LoggedOutToolBar />
      )}
    </AppBar>
  );
}

type LoggedInToolBarProps = {
  user: User | null;
};

function LoggedInToolBar({ user }: LoggedInToolBarProps) {
  return (
    <Toolbar
      sx={{ width: "40%", display: "flex", justifyContent: "space-between" }}
    >
      <Typography
        variant="h5"
        sx={(theme) => ({
          fontWeight: 600,
          color: theme.palette.text.secondary,
        })}
        noWrap
      >
        Welcome Back,{" "}
        {user ? (user.displayName ? user.displayName : "User") : "User"}!
      </Typography>
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
          doSignOut();
        }}
      >
        LOG-OUT
      </Button>
      {user ? (
        user.photoURL ? (
          <Avatar src={user.photoURL} />
        ) : user.displayName ? (
          <Avatar>{user.displayName}</Avatar>
        ) : (
          <Avatar src="" />
        )
      ) : (
        <Avatar src="" />
      )}
    </Toolbar>
  );
}
function LoggedOutToolBar() {
  const navigator = useNavigate();
  return (
    <Toolbar
      sx={{ width: "40%", display: "flex", justifyContent: "space-between" }}
    >
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
