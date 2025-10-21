import { Box } from "@mui/material";
import "./App.css";
import Navbar from "./Components/Navbar";
import { HashRouter, Route, Routes } from "react-router-dom";
import Homepage from "./Pages/Homepage";
import LoginPage from "./Pages/LogIn";
import SignUpPage from "./Pages/SignUp";
import DownloadPage from "./Pages/DownloadPage";
import LoggedInHomepage from "./Pages/LoggedInHomepage";

function App() {
  return (
    <Box
      sx={{
        margin: 0,
        bgcolor: "background.default",
        height: "100vh",
        overflow: "hidden",
      }}
    >
      <HashRouter>
        <Navbar isSignedIn={false}></Navbar>
        <Routes>
          <Route path="/" element={<Homepage />} />
          <Route path="/log_in" element={<LoginPage />} />
          <Route path="/sign_up" element={<SignUpPage />} />
          <Route path="/download" element={<DownloadPage />} />
          <Route path="/home" element={<LoggedInHomepage />} />
        </Routes>
      </HashRouter>
    </Box>
  );
}

export default App;
