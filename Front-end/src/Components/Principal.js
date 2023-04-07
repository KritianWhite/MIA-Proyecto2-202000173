import React, { useRef } from "react";
import "../Styles/Principal.css";

const Consoles2 = () => {
  const textAreaRef = useRef(null);

  const handleFileInputChange = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      const content = e.target.result;
      textAreaRef.current.value = content;
    };
    reader.readAsText(file);
  };

  return (
    <>
      <div className="principal">
        <div className="container">
          <div className="code-editor">
            <h2>Editor de código</h2>
            <textarea
              ref={textAreaRef}
              rows={25}
              cols={80}
              className="code-textarea"
            />
          </div>

          <div className="console">
            <h2>Consola</h2>
            <textarea rows={25} cols={80} className="console-textarea" />
          </div>

          <input
            onChange={handleFileInputChange}
            class="form-control form-control-lg"
            id="formFileLg"
            type="file"
          />

          <div class="d-grid gap-2">
            <button type="button" class="btn btn-light">
              Ejecutar
            </button>
          </div>
        </div>
      </div>
    </>
  );
};

export default Consoles2;
