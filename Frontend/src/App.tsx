import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Navbar from "./Components/Navbar";
import PublicHomePage from "./Pages/home/PublicHomepage";
import LogInPage from "./Pages/home/LogIn";
import SignUpPage from "./Pages/home/SignUp";
import AboutPage from "./Pages/home/AboutPage";
import UserSpace from "./Pages/user/UserSpace";
import HomePage from "./Pages/user/HomePage";
import SearchPage from "./Pages/user/SearchPage";
import SavedPage from "./Pages/user/SavedPage";
import CreatePage from "./Pages/user/CreatePage";
import ManageUploadsPage from "./Pages/user/ManageUploadsPage";
import SettingsPage from "./Pages/user/SettingsPage";
import DownloadPage from "./Pages/home/DownloadPage";
import RecordViewerPage from "./Pages/user/RecordViewerPage";

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
            {/* Fallback route to stop mistakenly navigating to viewRecord/ */}
            <Route element={<Navigate to={"/home"} />} path="viewRecord" />
            <Route element={<RecordViewerPage />} path="viewRecord/:uuid" />
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
