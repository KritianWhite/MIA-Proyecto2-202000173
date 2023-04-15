import React from "react";

function NavBar() {
  return (
    <>
      <nav class="navbar navbar-expand-lg bg-body-tertiary">
        <div class="container-fluid">
          <a class="navbar-brand" href="#">
            Home
          </a>
          <div class="collapse navbar-collapse" id="navbarScroll">
            <ul
              class="navbar-nav me-auto my-2 my-lg-0 navbar-nav-scroll"
            >
              <li class="nav-item">
              <a class="nav-link" href="#">
                  Principal
                </a>
              </li>
              <li class="nav-item">
                <a class="nav-link" href="#">
                  Login
                </a>
              </li>
              
            </ul>
          </div>
        </div>
      </nav>
    </>
  );
}

export default NavBar;
