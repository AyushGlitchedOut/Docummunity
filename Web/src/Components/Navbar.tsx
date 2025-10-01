import {
  Box,
  Button,
  Container,
  IconButton,
  SvgIcon,
  Typography,
} from "@mui/material";
import "@fontsource/roboto/400";
import Logo from "../assets/Docummunity.png";
import InfoOutlineIcon from "@mui/icons-material/InfoOutline";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";

interface NavbarProps {
  isSignedIn: boolean;
}

function Navbar(args: NavbarProps) {
  return (
    // Main Navbar begins here
    <Container
      sx={{
        backgroundColor: "background.paper",
        height: "8vh",
        maxWidth: "100vw",
        minWidth: "100vw",
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-between",
      }}
    >
      {/* Beginning of Navbar */}
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-evenly",
          flexDirection: "row",
          width: "30%",
          paddingTop: "3px",
        }}
      >
        <Box>
          {" "}
          <img src={Logo} style={{ height: "7vh", width: "7vh" }}></img>
        </Box>
        <Typography
          variant="h3"
          sx={{ color: "primary.main", letterSpacing: "5px" }}
        >
          DOCUMMUNITY
        </Typography>
      </Box>

      {/* Later part of Navbar */}
      {args.isSignedIn ? (
        <NavbarRightForSignedIn />
      ) : (
        <NavbarRightForAnonymous />
      )}
    </Container>
  );
}

function NavbarRightForSignedIn() {
  return (
    <Box
      sx={{
        width: "25%",
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        flexDirection: "row",
      }}
    >
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-evenly",
          flexDirection: "row",
          width: "70%",
        }}
      >
        <Button
          variant="outlined"
          sx={{
            color: "text.primary",
            borderColor: "divider",
            bgcolor: "background.default",
            ":hover": { boxShadow: "inset 0 0 20px rgba(128,128,128,0.3)" },
          }}
        >
          <Typography variant="button">Home</Typography>
        </Button>
        <Button
          variant="outlined"
          sx={{
            color: "text.primary",
            borderColor: "divider",
            bgcolor: "background.default",
            ":hover": { boxShadow: "inset 0 0 20px rgba(128,128,128,0.3)" },
          }}
        >
          <Typography variant="button">Download</Typography>
        </Button>
      </Box>
      <Box
        sx={{
          height: "6vh",
          width: "6vh",
          borderRadius: "50%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <SvgIcon component={AccountCircleIcon} sx={{ fontSize: "7vh" }} />
      </Box>
    </Box>
  );
}

function NavbarRightForAnonymous() {
  return (
    <Box
      sx={{
        width: "35%",
        display: "flex",
        justifyContent: "space-evenly",
        alignItems: "center",
        flexDirection: "row",
      }}
    >
      <Button
        variant="outlined"
        sx={{
          color: "text.primary",
          borderColor: "divider",
          bgcolor: "background.default",
          ":hover": { boxShadow: "inset 0 0 20px rgba(128,128,128,0.3)" },
        }}
      >
        <Typography variant="button">Log-In</Typography>
      </Button>
      <Button
        variant="outlined"
        sx={{
          color: "text.primary",
          borderColor: "divider",
          bgcolor: "background.default",
          ":hover": { boxShadow: "inset 0 0 20px rgba(128,128,128,0.3)" },
        }}
      >
        <Typography variant="button">Sign-Up</Typography>
      </Button>
      <Button
        variant="outlined"
        sx={{
          color: "text.primary",
          borderColor: "divider",
          bgcolor: "background.default",
          ":hover": { boxShadow: "inset 0 0 20px rgba(128,128,128,0.3)" },
        }}
      >
        <Typography variant="button">Download</Typography>
      </Button>
      <IconButton aria-label="Info" color="secondary">
        <SvgIcon component={InfoOutlineIcon} fontSize={"large"} />
      </IconButton>
    </Box>
  );
}

export default Navbar;
