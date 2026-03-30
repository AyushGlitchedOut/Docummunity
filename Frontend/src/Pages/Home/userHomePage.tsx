import { Box } from "@mui/material";
import SideBar from "../../Components/SideBar";
import { useAuth } from "../../auth/fireBaseContext";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

function UserHomePage() {
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
    }
  }, [auth]);
  return (
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
        Main Content
      </Box>
    </Box>
  );
}
export default UserHomePage;
