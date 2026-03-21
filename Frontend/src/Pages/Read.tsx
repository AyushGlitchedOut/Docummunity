import { useState } from "react";
import "./Pages.css";

function ReadPage() {
  const [ID, setID] = useState<string>("");

  async function downloadFile(): Promise<void> {
    if (!ID) {
      alert("No ID provided");
      return;
    }

    const a = document.createElement("a");
    a.href = "http://localhost:8080/download/" + encodeURIComponent(ID);
    a.click();
  }

  return (
    <div className="read-tab">
      <input
        placeholder="Enter ID here"
        className="read-input"
        value={ID}
        onChange={(event) => {
          setID(event.target.value);
        }}
      />
      <button
        className="read-download"
        onClick={() => {
          downloadFile();
        }}
      >
        DOWNLOAD
      </button>
    </div>
  );
}
export default ReadPage;
