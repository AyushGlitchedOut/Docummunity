import {useState} from 'react'
import "./App.css";
import axios from 'axios';
import { Button } from '@mui/material';


function App() {

  const [message, setMessage] = useState("")

  async function getData():Promise<void> {
    const response = await axios.get("http://localhost:8000/test")
    setMessage(response.data)
  }

  return (
    <>
      <h1>{message}</h1>
      <Button onClick={() => getData()} variant='contained' color='secondary'>Get Message</Button>
    </>
  );
}

export default App;
