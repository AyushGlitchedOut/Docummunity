import { AppBar, Box, Button, Toolbar, Typography } from "@mui/material";
import { useAuth } from "../auth/fireBaseContext";
import AppIcon from "../assets/Docummunity.png";
function Navbar() {
  const auth = useAuth();
  var userLoggedIn: boolean = false;
  if (auth && auth.userLoggedIn) {
    userLoggedIn = true;
  }
  return (
    <AppBar
      position="static"
      sx={(theme) => ({
        height: "7vh",
        boxShadow: `2px 2px 2px ${theme.palette.divider}`,
        backgroundColor: theme.palette.background.paper,
        display: "flex",
        alignItems: "center",
        flexDirection: "row",
        justifyContent: "space-between",
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
        >
          DOCUMMUNITY
        </Typography>
      </Toolbar>
      {userLoggedIn ? <LoggedInToolBar /> : <LoggedOutToolBar />}
    </AppBar>
  );
}

function LoggedInToolBar() {
  return <Toolbar>Logged In</Toolbar>;
}
function LoggedOutToolBar() {
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
      >
        ABOUT US
      </Button>
    </Toolbar>
  );
}

export default Navbar;
