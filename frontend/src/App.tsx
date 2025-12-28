import { useState } from "react";
import "./App.css";
import { Greet } from "../wailsjs/go/main/App";
import MainPage from "./components/MainPage";
import Header from "./components/Header";

function App() {
  const [resultText, setResultText] = useState(
    "Please enter your name below 👇",
  );
  const [name, setName] = useState("");
  const updateName = (e: any) => setName(e.target.value);
  const updateResultText = (result: string) => setResultText(result);

  function greet() {
    Greet(name).then(updateResultText);
  }

  return (
    <div id="App">
      <Header />
      <MainPage />
    </div>
  );
}

export default App;
