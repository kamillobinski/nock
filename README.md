<a name="readme-top"></a>



<!-- PROJECT LOGO -->
<br />
<div align="center">
  <img src="docs/readme/icon.png" alt="Logo" width="80" height="80">
  <h3 align="center">Nock</h3>
  <p align="center">
    A personal CLI tool for everyday terminal tasks.
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Nock is my personal CLI tool for the small things I do every day in the terminal. Jumps to directories I've saved, runs scripts from one directory and prints system info. Written in Go. Does exactly what I need and nothing more.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Nock is built from source and needs Go.

1. Clone the repository
   ```sh
   git clone https://github.com/kamillobinski/nock.git
   cd nock
   ```

2. Build and install the binary
   ```sh
   go build -o nock . && sudo mv nock /usr/local/bin/nock
   ```

3. Set up shell integration (zsh or bash)
   ```sh
   nock init zsh
   source ~/.zshrc
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
