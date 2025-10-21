import { Box, Button } from "@mui/material";
import SidebarButton from "../Components/SidebarButton";
import {
  AccountCircle,
  AddCircleOutline,
  Bookmarks,
  Home,
  Inventory,
  Search,
  Settings,
} from "@mui/icons-material";
import { useFirebase } from "../services/firebase";
import { useNavigate } from "react-router-dom";
import { useEffect } from "react";

function LoggedInHomepage() {
  const firebase = useFirebase();
  const navigator = useNavigate();

  async function handleSignOut() {
    await firebase.LogOut();
  }

  useEffect(() => {
    if (!firebase.isLoggedIn) {
      navigator("/");
    }
  }, [firebase]);

  return (
    <Box
      sx={{
        height: "92vh",
        width: "100%",
        display: "flex",
        flexDirection: "row",
        justifyContent: "space-evenly",
        alignItems: "center",
      }}
    >
      {/* Sidebar */}
      <Sidebar />

      {/* Main Content */}
      <Box
        sx={{
          backgroundColor: "background.paper",
          height: "90%",
          width: "80%",
          boxShadow: "2px 2px 1px black, inset 2px 2px 5px grey",
        }}
      >
        <Button
          onClick={() => {
            handleSignOut();
          }}
        >
          Log Out
        </Button>
      </Box>
    </Box>
  );
}

function Sidebar() {
  return (
    <Box
      sx={{
        backgroundColor: "background.paper",
        height: "90%",
        width: "18%",
        boxShadow: "2px 2px 1px black, inset 2px 2px 5px grey",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <SidebarButton text="Home" Icon={<Home fontSize="large" />} />
      <SidebarButton text="Account" Icon={<AccountCircle fontSize="large" />} />
      <SidebarButton text="Search" Icon={<Search fontSize="large" />} />
      <SidebarButton
        text="Create"
        Icon={<AddCircleOutline fontSize="large" />}
      />
      <SidebarButton text="Library" Icon={<Bookmarks fontSize="large" />} />
      <SidebarButton
        text="Manage Uploads"
        Icon={<Inventory fontSize="large" />}
      />
      <SidebarButton text="Settings" Icon={<Settings fontSize="large" />} />
    </Box>
  );
}

export default LoggedInHomepage;
