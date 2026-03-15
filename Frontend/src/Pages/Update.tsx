import { useState } from "react";
import "./Pages.css";
function UpdatePage() {
  const [ID, setID] = useState<string>("");
  const [file, setFile] = useState<File | null>(null);

  async function handleUpdate(): Promise<void> {
    if (!ID) {
      alert("No ID given");
      return;
    }
    if (!file) {
      alert("No File Selected");
      return;
    }

    const fileForm = new FormData();
    fileForm.append("file", file);

    try {
      const response = await fetch("http://localhost:8080/update/" + ID, {
        method: "PUT",
        body: fileForm,
      });

      alert(response.status);
    } catch (error) {
      alert("Something Went Wrong");
    }
  }

  return (
    <div className="update-tab">
      <input
        placeholder="Enter ID to be updated"
        className="update-ID"
        value={ID}
        onChange={(event) => setID(event.target.value)}
      />
      <input
        type="file"
        hidden={true}
        id="update-file"
        maxLength={40}
        onChange={(event) => {
          if (!event.target.files) return;
          setFile(event.target.files[0]);
        }}
      />
      <label htmlFor="update-file" className="update-container">
        SELECT FILE
      </label>
      <p>File Chosen: {file ? file.name : "No File Selected"}</p>
      <button className="update-submit" onClick={() => handleUpdate()}>
        UPDATE
      </button>
    </div>
  );
}
export default UpdatePage;
