import React from "react";

export default function Editor(props) {
  //* Manejador de eventos para actualizar el valor del text area 
  const handlerChangeEditor = (evt) => {
    props.handlerChange(evt.target.value);
  };

  return (
    <>
      <h2>{props.Tittle}</h2>
      <textarea
        rows="20"
        cols="80"
        id="code-textarea"
        class="code-textarea"
        value={props.contenido}
        onChange={handlerChangeEditor}

        readOnly={props.readOnly}
      ></textarea>
    </>
  );
}
