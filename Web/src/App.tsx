import {useState} from 'react'
import "./App.css";
import axios from 'axios';


function App() {

  const [message, setMessage] = useState("")

  async function getData():Promise<void> {
    const response = await axios.get("http://localhost:8000/test")
    setMessage(response.data)
  }

  return (
    <>
      <h1>{message}</h1>
      <button onClick={() => getData()}>Get Message</button>
    </>
  );
}

export default App;
