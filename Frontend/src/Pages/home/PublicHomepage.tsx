import { Box } from "@mui/material";
import { useAuth } from "../../auth/fireBaseContext";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

function PublicHomePage() {
  const auth = useAuth();
  const navigator = useNavigate();
  useEffect(() => {
    if (!auth) {
      return;
    }
    if (auth.loading) {
      return;
    }
    if (auth && auth.userLoggedIn) {
      navigator("/home");
    }
  }, [auth]);
  return (
    <Box sx={{ height: "93vh", width: "100%", marginTop: "7vh" }}>Hello</Box>
  );
}

export default PublicHomePage;
