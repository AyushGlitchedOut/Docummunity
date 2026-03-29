import { BrowserRouter } from "react-router-dom";
import Navbar from "./Components/Navbar";
import "./App.css";

function App() {
  return (
    <div className="home">
      <BrowserRouter>
        <Navbar />
      </BrowserRouter>
    </div>
  );
}

export default App;
