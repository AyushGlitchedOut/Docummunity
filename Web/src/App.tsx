import { Box } from "@mui/material";
import "./App.css";
import Navbar from "./Components/Navbar";

function App() {
  return (
    <Box sx={{ margin: 0, bgcolor: "background.default", height: "100vh" }}>
      <Navbar isSignedIn={true}></Navbar>
    </Box>
  );
}

export default App;
