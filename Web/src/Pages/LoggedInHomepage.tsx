import { Box } from "@mui/material";
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
import { useEffect, useState } from "react";
import HomeTab from "./subPages/HomeTab";
import AccountTab from "./subPages/AccountTab";
import SearchTab from "./subPages/SearchTab";
import CreateTab from "./subPages/CreateTab";
import LibraryTab from "./subPages/LibraryTab";
import ManageUploadsTab from "./subPages/ManageUploadsTab";
import SettingsTab from "./subPages/SettingsTab";
import { useNavigate } from "react-router-dom";

const stringToTab = {
  home: <HomeTab />,
  account: <AccountTab />,
  search: <SearchTab />,
  create: <CreateTab />,
  library: <LibraryTab />,
  manage_uploads: <ManageUploadsTab />,
  settings: <SettingsTab />,
};

type stringToTabKey = keyof typeof stringToTab;

function LoggedInHomepage() {
  const navigator = useNavigate();
  const [currentTab, setCurrentTab] = useState<stringToTabKey>("home");

  //handle search params (if any)
  useEffect(() => {
    const url = window.location.hash;
    const params = url.split("?")[1];
    if (params === undefined || params == "") return;
    const args = new URLSearchParams(params);
    const tabArgs = args.get("tab");

    if (tabArgs && tabArgs in stringToTab) {
      setCurrentTab(tabArgs as stringToTabKey);
    }
  }, []);

  useEffect(() => {
    if (currentTab == "home") return;
    navigator(`?tab=${currentTab}`);
  }, [currentTab]);

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
      <Sidebar activeTab={currentTab} setTab={setCurrentTab} />

      {/* Main Content */}
      <Box
        sx={{
          backgroundColor: "background.paper",
          height: "90%",
          width: "80%",
          boxShadow: "2px 2px 1px black, inset 2px 2px 5px grey",
        }}
      >
        {stringToTab[currentTab]}
      </Box>
    </Box>
  );
}

interface SidebarProps {
  activeTab: string;
  setTab: React.Dispatch<React.SetStateAction<stringToTabKey>>;
}

function Sidebar(args: SidebarProps) {
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
      <SidebarButton
        text="Home"
        Icon={<Home fontSize="large" />}
        ActiveTab={args.activeTab == "home" || args.activeTab == ""}
        callback={() => {
          args.setTab("home");
        }}
      />
      <SidebarButton
        text="Account"
        Icon={<AccountCircle fontSize="large" />}
        ActiveTab={args.activeTab == "account"}
        callback={() => {
          args.setTab("account");
        }}
      />
      <SidebarButton
        text="Search"
        Icon={<Search fontSize="large" />}
        ActiveTab={args.activeTab == "search"}
        callback={() => {
          args.setTab("search");
        }}
      />
      <SidebarButton
        text="Create"
        Icon={<AddCircleOutline fontSize="large" />}
        ActiveTab={args.activeTab == "create"}
        callback={() => {
          args.setTab("create");
        }}
      />
      <SidebarButton
        text="Library"
        Icon={<Bookmarks fontSize="large" />}
        ActiveTab={args.activeTab == "library"}
        callback={() => {
          args.setTab("library");
        }}
      />
      <SidebarButton
        text="Manage Uploads"
        Icon={<Inventory fontSize="large" />}
        ActiveTab={args.activeTab == "manage_uploads"}
        callback={() => {
          args.setTab("manage_uploads");
        }}
      />
      <SidebarButton
        text="Settings"
        Icon={<Settings fontSize="large" />}
        ActiveTab={args.activeTab == "settings"}
        callback={() => {
          args.setTab("settings");
        }}
      />
    </Box>
  );
}

export default LoggedInHomepage;
