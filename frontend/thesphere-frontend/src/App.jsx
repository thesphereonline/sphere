import React, { useEffect, useState } from 'react';

function App() {
  const [message, setMessage] = useState("");

  useEffect(() => {
    fetch("http://localhost:3001")
      .then((res) => res.json())
      .then((data) => setMessage(data.message))
      .catch((err) => console.error(err));
  }, []);

  return (
    <div>
      <h1>The Sphere Frontend</h1>
      <p>{message}</p>
    </div>
  )
}

export default App;