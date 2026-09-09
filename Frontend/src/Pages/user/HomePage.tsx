import { Box, Button } from "@mui/material";
import { useAuth } from "../../auth/fireBaseContext";

function HomePage() {
  const auth = useAuth();

  async function sendAuthDetails(): Promise<void> {
    if (!auth) {
      alert("Authentication Error");
      return;
    }
    if (!auth.userLoggedIn) {
      alert("User Not Logged In!");
      return;
    }
    if (!auth.currentUser) {
      alert("Authentication Error");
      return;
    }

    const token = await auth.currentUser.getIdToken();

    const result = await fetch("http://localhost:8080/api/test", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    console.log(result.status);
  }

  return (
    <Box sx={{ height: "100%", width: "100%", borderRadius: "20px" }}>
      <Button
        onClick={() => {
          sendAuthDetails();
        }}
      >
        Send JWT Test
      </Button>
    </Box>
  );
}

export default HomePage;
