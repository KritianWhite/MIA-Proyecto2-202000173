import React from "react";
import "../Styles/Principal.css";

const CodeEditor = () => {
  const [code, setCode] = React.useState("");

  const handleCodeChange = (event) => {
    setCode(event.target.value);
  };

  return (
    <div className="code-editor">
      <h2>Editor de código</h2>

      <textarea
        rows={25}
        cols={80}
        value={code}
        onChange={handleCodeChange}
        className="code-textarea"
      />
    </div>
  );
};

const Console = () => {
  const [output, setOutput] = React.useState("");

  const handleOutputChange = (event) => {
    setOutput(event.target.value);
  };

  return (
    <div className="console">
      <h2>Consola</h2>
      <textarea
        rows={25}
        cols={80}
        value={output}
        onChange={handleOutputChange}
        className="console-textarea"
      />
    </div>
  );
};

const Consoles = () => {
  return (
    <div className="principal">
      <div className="container">
        <CodeEditor />
        <Console />
      </div>
    </div>
  );
};

export default Consoles;
