import { Box } from "@mui/material";
import SideBar from "../../Components/SideBar";
import { useAuth } from "../../auth/fireBaseContext";
import { useEffect, useState } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import type { UserInfo } from "../../models/models";
import { BACKEND_URL } from "../../consts";
import { doSignOut } from "../../auth/authFunctions";
import { UserContext } from "../../Contexts/UserContext";

function UserSpace() {
  const [user, setUser] = useState<UserInfo | null>(null);
  const auth = useAuth();
  const navigator = useNavigate();
  useEffect(() => {
    if (!auth) {
      return;
    }
    if (auth.loading) {
      return;
    }
    if (!auth.userLoggedIn) {
      navigator("/", { replace: true });
      return;
    }

    async function fetchUser(): Promise<void> {
      const response = await fetch(BACKEND_URL + "/user/ACCOUNT", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${await auth?.currentUser?.getIdToken()}`,
        },
      });

      if (!response.ok) {
        await doSignOut();
        alert("No User Found");
        return;
      }
      const resJSON = await response.json();
      const newUser: UserInfo = resJSON.message;
      setUser(newUser);
    }

    fetchUser();
  }, [auth]);
  return (
    <UserContext.Provider value={user}>
      <Box
        sx={{
          height: "93vh",
          marginTop: "7vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          flexDirection: "row",
        }}
      >
        <SideBar />
        <Box
          sx={(theme) => ({
            border: `1px solid ${theme.palette.divider}`,
            margin: "1%",
            height: "90%",
            width: "100%",
            borderRadius: "20px",
            backgroundColor: theme.palette.secondary.main,
          })}
        >
          <Outlet />
        </Box>
      </Box>
    </UserContext.Provider>
  );
}
export default UserSpace;
