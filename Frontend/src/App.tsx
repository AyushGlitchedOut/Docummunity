import { useState } from "react";
import "./App.css";

function App() {
  const [message, setMessage] = useState("loading...");
  const [input, setInput] = useState("");

  async function loadMessage(): Promise<void> {
    const response = await fetch("http://localhost:8080/");
    const data = await response.json();
    setMessage(data.message);
  }

  async function sendMessage(): Promise<void> {
    const response = await fetch("http://localhost:8080/message", {
      method: "POST",

      body: JSON.stringify({
        name: "Ayush",
        message: "Docummunity is HERE!!!!!!!!!!",
      }),
    });
    const data = await response.status;
    alert(data);
  }

  return (
    <div>
      {message}
      <br></br>
      <button onClick={() => loadMessage()}>LOAD DATA FROM SERVER</button>
      <br></br>
      <input
        value={input}
        onChange={(event) => {
          setInput(event.target.value);
        }}
      />
      <br></br>
      <button onClick={() => sendMessage()}>SEND DATA TO SERVER</button>
    </div>
  );
}

export default App;
