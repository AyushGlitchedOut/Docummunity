import { Box, Button } from "@mui/material";
import { doSignOut } from "../../auth/authFunctions";

function SettingsPage() {
  return (
    <Box sx={{ height: "100%", width: "100%", borderRadius: "20px" }}>
      <Button
        variant="contained"
        onClick={() => {
          doSignOut();
        }}
      >
        LOG OUT
      </Button>
    </Box>
  );
}

export default SettingsPage;
