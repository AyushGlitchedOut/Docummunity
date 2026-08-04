import {
  Archive,
  Bookmark,
  Create,
  HomeFilled,
  SearchOutlined,
  Settings,
} from "@mui/icons-material";
import {
  Avatar,
  Box,
  Divider,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
} from "@mui/material";
import { useEffect, useState, type JSX } from "react";
import {
  useLocation,
  useNavigate,
  type NavigateFunction,
} from "react-router-dom";
import { useAuth } from "../auth/fireBaseContext";
import type { UserInfo } from "../models/models";
import { BACKEND_URL } from "../consts";

function SideBar() {
  const location = useLocation();
  const navigator = useNavigate();
  const auth = useAuth();
  const [UserInfo, setUserInfo] = useState<UserInfo>();
  const [AvatarURL, setAvatarURL] = useState<string>();

  async function fetchUserInfo(): Promise<void> {
    if (!auth?.currentUser) return;
    const token = await auth?.currentUser?.getIdToken();
    if (!token) {
      alert("Auth Issues");
      return;
    }

    const response = await fetch(BACKEND_URL + "/user/ACCOUNT", {
      method: "GET",
      headers: { Authorization: `Bearer ${token}` },
    });

    if (!response.ok) {
      alert("Not OK");
    }
    const resJSON = await response.json();
    console.log(resJSON);
    const newUser: UserInfo = {
      BIO: resJSON.message.BIO,
      UID: resJSON.message.UID,
      DISPLAY_NAME: resJSON.message.DISPLAY_NAME,
      PROFILE_PIC: resJSON.message.PROFILE_PIC,
      SETTINGS: resJSON.message.SETTINGS,
    };
    setUserInfo(newUser);
    const fetchedAvatarURL = newUser.PROFILE_PIC.split("/");
    setAvatarURL(fetchedAvatarURL[fetchedAvatarURL.length - 1]);
  }

  useEffect(() => {
    if (!auth) return;
    fetchUserInfo();
  }, [auth]);

  return (
    <Box
      sx={(theme) => ({
        border: `1px solid ${theme.palette.divider}`,
        backgroundColor: theme.palette.secondary.main,
        borderRadius: "10px",
        margin: "1%",
        width: "25%",
        height: "80%",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
      })}
    >
      <Box
        sx={{
          height: "10%",
          width: "90%",
          display: "flex",
          flexDirection: "row",
          alignItems: "center",
          justifyContent: "left",
        }}
      >
        {AvatarURL ? (
          <Avatar src={`${BACKEND_URL}/user/PROFILE_PIC/${AvatarURL}`} />
        ) : (
          <Avatar src="" />
        )}
        <Box
          sx={{
            height: "100%",
            width: "80%",
            margin: "5%",
            display: "flex",
            flexDirection: "column",
            cursor: "pointer",
          }}
          onClick={() => {
            navigator("/home/settings");
          }}
        >
          <Typography>Welcome Back, </Typography>
          <Typography>{UserInfo ? UserInfo.DISPLAY_NAME : "User"}</Typography>
        </Box>
      </Box>
      <Divider
        sx={(theme) => ({
          backgroundColor: theme.palette.divider,
          width: "90%",
        })}
      ></Divider>
      <Box
        sx={{
          width: "90%",
        }}
      >
        <List>
          <SideBarButton
            PageName="HOMEPAGE"
            NavigationAddress="/home/"
            Icon={<HomeFilled />}
            navigator={navigator}
            location={location.pathname}
          />
          <SideBarButton
            PageName="SEARCH"
            NavigationAddress="/home/search"
            Icon={<SearchOutlined />}
            navigator={navigator}
            location={location.pathname}
          />
          <SideBarButton
            PageName="SAVED"
            NavigationAddress="/home/saved"
            Icon={<Bookmark />}
            navigator={navigator}
            location={location.pathname}
          />
        </List>
      </Box>
      <Divider
        sx={(theme) => ({
          backgroundColor: theme.palette.divider,
          width: "90%",
        })}
      ></Divider>
      <Box
        sx={{
          width: "90%",
        }}
      >
        <SideBarButton
          PageName="CREATE"
          NavigationAddress="/home/create"
          Icon={<Create />}
          navigator={navigator}
          location={location.pathname}
        />
        <SideBarButton
          PageName="MANAGE UPLOADS"
          NavigationAddress="/home/manage_uploads"
          Icon={<Archive />}
          navigator={navigator}
          location={location.pathname}
        />
      </Box>
      <Divider
        sx={(theme) => ({
          backgroundColor: theme.palette.divider,
          width: "90%",
        })}
      ></Divider>
      <Box
        sx={{
          width: "90%",
        }}
      >
        <SideBarButton
          PageName="SETTINGS"
          NavigationAddress="/home/settings"
          Icon={<Settings />}
          navigator={navigator}
          location={location.pathname}
        />
      </Box>
    </Box>
  );
}
export default SideBar;

type SideBarButtonProps = {
  PageName: string;
  NavigationAddress: string;
  Icon: JSX.Element;
  navigator: NavigateFunction;
  location: string;
};

function SideBarButton({
  PageName,
  NavigationAddress,
  Icon,
  navigator,
  location,
}: SideBarButtonProps) {
  return (
    <ListItemButton
      onClick={() => {
        navigator(NavigationAddress);
      }}
      sx={(theme) => ({
        backgroundColor:
          location == NavigationAddress
            ? theme.palette.primary.main
            : "inherit",
      })}
    >
      <ListItemIcon
        sx={{
          "& svg": {
            fontSize: "2.4rem",
          },
        }}
      >
        {Icon}
      </ListItemIcon>
      <ListItemText
        primary={
          <Typography
            variant="button"
            sx={{ fontSize: "1rem", fontWeight: "600" }}
          >
            {PageName}
          </Typography>
        }
      ></ListItemText>
    </ListItemButton>
  );
}
