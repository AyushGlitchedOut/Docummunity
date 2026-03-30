import { useState } from "react";
import "./Pages.css";

function CreatePage() {
  const [file, setFile] = useState<File | null>(null);
  const [ID, setID] = useState<string>("");

  async function handleSubmit(): Promise<void> {
    if (!file) {
      alert("No File Selected");
      return;
    }
    if (ID == "") {
      alert("Enter a valid ID");
      return;
    }
    const form = new FormData();
    form.append("file", file);
    form.append("ID", ID);

    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: form,
      });
      if (response.status == 200) {
        alert("File Uploaded");
      } else {
        alert("Error!");
      }
    } catch (error) {
      alert("Something Went Wrong");
      console.log(error);
    }
  }

  return (
    <div className="create-tab">
      <input
        type="file"
        hidden={true}
        id="upload-input"
        onChange={(event) => {
          if (!event.target.files) return;
          setFile(event.target.files[0]);
        }}
      />
      <label htmlFor="upload-input" className="upload-container">
        SELECT FILE
      </label>
      <p>File Chosen: {file ? file.name : "No File Selected"}</p>
      <input
        placeholder="Enter Id here..."
        className="upload-ID"
        maxLength={40}
        value={ID}
        onChange={(event) => {
          setID(event.target.value);
        }}
      />
      <button className="upload-submit" onClick={() => handleSubmit()}>
        Upload File
      </button>
    </div>
  );
}
export default CreatePage;
