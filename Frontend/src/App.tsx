import { BrowserRouter, Route, Routes } from "react-router-dom";
import Navbar from "./Components/Navbar";
import "./App.css";
import HomePage from "./Pages/auth/Homepage";
import LogInPage from "./Pages/auth/LogIn";
import SignUpPage from "./Pages/auth/SignUp";
import AboutPage from "./Pages/auth/AboutPage";
import UserHomePage from "./Pages/Home/userHomePage";

function App() {
  return (
    <div className="home">
      <BrowserRouter>
        <Navbar />
        <Routes>
          <Route element={<HomePage />} path="/" />
          <Route element={<LogInPage />} path="/login" />
          <Route element={<SignUpPage />} path="/signUp" />
          <Route element={<AboutPage />} path="/about" />
          <Route element={<UserHomePage />} path="/home"></Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
