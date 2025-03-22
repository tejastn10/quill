<p align="center">
  <img src="logo.svg" alt="Logo">
</p>

# Quill ✒️

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tejastn10/quill?logo=go)
[![Unit Tests](https://github.com/tejastn10/quill/actions/workflows/unit-test.yml/badge.svg?logo=github)](https://github.com/tejastn10/quill/actions/workflows/unit-test.yml)
[![Security Audit](https://github.com/tejastn10/quill/actions/workflows/security-audit.yml/badge.svg?logo=github)](https://github.com/tejastn10/quill/actions/workflows/security-audit.yml)
![License](https://img.shields.io/badge/License-MIT-yellow?logo=open-source-initiative&logoColor=white)

**Quill** is a lightweight Git clone built from scratch to understand the inner workings of Git. It allows you to explore how Git handles version control, including file diffs, staging, commits, and object storage. Designed with precision and simplicity in mind.

---

## Features 🌟

- **Repository Initialization**: Set up a new version-controlled repository.
- **Object Storage**: Understand Git's internal storage of blobs, trees, and commits.
- **Staging Area**: Track changes and prepare files for commits.
- **Commit History**: Record changes with metadata and parent tracking.
- **Diffing**: Compare file versions to highlight changes.

---

## Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/doc/install) installed on your machine.
- [Git](https://git-scm.com/) for comparison and testing purposes (optional).

### Installation ⚙️

1. Clone this repository:

    ```bash
    git clone https://github.com/tejastn10/quill.git
    cd quill
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Run the project:

    ```bash
    go run main.go
    ```

---

## Usage

### Initialize a Repository

Run the following command to initialize a new Quill repository in your current directory:

```bash
./quill init
```

This will create a .quill directory, which includes subdirectories for storing objects, references, and configuration.

Example Workflow

1. Add a file:
   Add a file to the staging area:

   ```bash
   ./quill add <file>
   ```

2. Commit changes:
   Commit staged changes with a message:

   ```bash
   ./quill commit -m "Initial commit"
   ```

3. View commit log:
   Display the commit history:

   ```bash
   ./quill log
   ```

---

## Project Structure 📂

```md
quill/
├── cmd/            # CLI commands (user-facing commands like init, add, commit, etc.)
├── pkg/            # Core functionality
│   ├── hash/       # Hashing algorithms and utilities
│   ├── objects/    # Blob, tree, commit handling
│   ├── refs/       # Branch and HEAD management
│   ├── repo/       # Initialization of .quill directory
│   ├── index/      # Staging area implementation
│   └── storage/    # Low-level File I/O operations
├── internal/       # Internal utilities and helpers (e.g., logging, config)
├── .gitignore      # Ignore build artifacts
├── main.go         # Entry point for the CLI application
└── README.md       # Documentation
```

---

## License 📜

This project is licensed under the MIT License. See the [LICENSE](LICENSE.md) file for details.

---

## Acknowledgments 🙌

- Inspired by the concept of **precision and versioning** in writing.
- Built with ❤️ and Go.
