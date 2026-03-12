import { useState } from 'react'
import './App.css'

function App() {

  const [message, setMessage] = useState("loading...")


  async function loadMessage():Promise<void> {
    const response = await fetch("http://localhost:8080/")
    const data = await response.json()
    setMessage(data.message)
  }

  return (
    <>
      {message}
      <button onClick={()=> loadMessage()}>LOAD DATA FROM SERVER</button>
    </>
  )
}

export default App
