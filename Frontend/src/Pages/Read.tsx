import { useState } from "react";
import "./Pages.css";

function ReadPage() {
  const [ID, setID] = useState<string>("");

  function parseFileNameHeader(arg: string | null): string {
    if (!arg) {
      return "download";
    }
    const match = arg.match(/filename="(.+?)"/);
    return match ? match[1] : "download";
  }

  async function downloadFile(): Promise<void> {
    if (!ID) {
      alert("No ID provided");
      return;
    }

    const res = await fetch(
      "http://localhost:8080/download/" + encodeURIComponent(ID),
    );
    if (!res.ok) {
      alert("Something Went Wrong");
      return;
    }

    const dataBlob = await res.blob();
    const dataBlobUrl = URL.createObjectURL(dataBlob);

    const a = document.createElement("a");
    a.href = dataBlobUrl;

    const fileName = parseFileNameHeader(
      res.headers.get("Content-Disposition"),
    );
    a.download = fileName;

    document.body.appendChild(a);
    a.click();
    a.remove();

    URL.revokeObjectURL(dataBlobUrl);
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
