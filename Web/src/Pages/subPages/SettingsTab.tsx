import { Button, Container } from "@mui/material";
import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useFirebase } from "../../services/firebase";

function SettingsTab() {
  const navigator = useNavigate();
  const firebase = useFirebase();
  useEffect(() => {
    if (!firebase.isLoggedIn) {
      navigator("/");
    }
  }, [firebase]);

  async function handleSignOut() {
    await firebase.LogOut();
  }
  return (
    <Container>
      Settings
      <Button
        onClick={() => {
          handleSignOut();
        }}
      >
        Log Out
      </Button>
    </Container>
  );
}

export default SettingsTab;
