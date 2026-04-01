import { BrowserRouter, Route, Routes } from "react-router-dom";
import Navbar from "./Components/Navbar";
import PublicHomePage from "./Pages/home/PublicHomepage";
import LogInPage from "./Pages/home/LogIn";
import SignUpPage from "./Pages/home/SignUp";
import AboutPage from "./Pages/home/AboutPage";
import UserSpace from "./Pages/userHome/UserSpace";
import HomePage from "./Pages/userHome/HomePage";
import SearchPage from "./Pages/userHome/SearchPage";
import SavedPage from "./Pages/userHome/SavedPage";
import CreatePage from "./Pages/userHome/CreatePage";
import ManageUploadsPage from "./Pages/userHome/ManageUploadsPage";
import SettingsPage from "./Pages/userHome/SettingsPage";
import DownloadPage from "./Pages/home/DownloadPage";

function App() {
  return (
    <div className="home">
      <BrowserRouter>
        <Navbar />
        <Routes>
          <Route element={<PublicHomePage />} path="/" />
          <Route element={<DownloadPage />} path="/download" />
          <Route element={<LogInPage />} path="/login" />
          <Route element={<SignUpPage />} path="/signUp" />
          <Route element={<AboutPage />} path="/about" />
          <Route element={<UserSpace />} path="/home">
            <Route element={<HomePage />} path="" />
            <Route element={<SearchPage />} path="search" />
            <Route element={<SavedPage />} path="saved" />
            <Route element={<CreatePage />} path="create" />
            <Route element={<ManageUploadsPage />} path="manage_uploads" />
            <Route element={<SettingsPage />} path="settings" />
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
