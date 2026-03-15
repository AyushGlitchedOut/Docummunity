import { useState } from "react";
import "./Pages.css";

function DeletePage() {
  const [ID, setID] = useState<string>("");
  async function handleDelete(): Promise<void> {
    if (!ID) {
      alert("No Id Provided");
      return;
    }
    try {
      const response = await fetch("http://localhost:8080/delete/" + ID, {
        method: "DELETE",
      });
      alert(response.status);
    } catch (error) {
      alert("Something Went Wrong");
    }
  }
  return (
    <div className="delete-tab">
      <input
        placeholder="ID for Deleting"
        className="delete-input"
        value={ID}
        onChange={(event) => setID(event.target.value)}
      />
      <button
        className="delete-submit"
        onClick={() => {
          handleDelete();
        }}
      >
        DELETE
      </button>
    </div>
  );
}
export default DeletePage;
