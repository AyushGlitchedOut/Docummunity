import { useState, type JSX } from "react";
import "./App.css";
import Navbar from "./Components/Navbar";
import CreatePage from "./Pages/Create";

function App() {
  const [page, setPage] = useState<JSX.Element>(<CreatePage />);

  return (
    <div className="home">
      <Navbar setPage={setPage} />
      <div className="page-display">{page}</div>
    </div>
  );
}

export default App;
