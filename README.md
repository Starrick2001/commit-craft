# Commit Craft

Commit Craft is a Go-based tool that leverages the Gemini API to automatically generate descriptive commit messages from your staged changes. It streamlines your git workflow by creating conventional commit messages for you.

## Prerequisites

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) must be installed.
- You need a Gemini API key.

## Installation

1. **Download the binary:**
    Download the latest binary for your operating system from the [Releases page](https://github.com/starrick2001/commit-craft/releases).

2. **Make it executable:**

    ```bash
    chmod +x commit-craft
    ```

3. **Move to your PATH (Optional):**
    For easy access, move the binary to a directory in your system's PATH.

    ```bash
    sudo mv commit-craft /usr/local/bin/
    ```

## Setup

Set up your Gemini API Key as an environment variable:

```bash
export COMMIT_CRAFT_GEMINI_KEY='YOUR_API_KEY'
```

> **Note:** To make this permanent, add the export command to your shell's configuration file (e.g., `~/.bashrc`, `~/.zshrc`).

## Usage

1. **Stage your changes:**
    Use `git add` to stage the files you want to commit.

    ```bash
    git add .
    ```

2. **Run Commit Craft:**
    Execute the program to generate a commit message and apply it.

    ```bash
    commit-craft
    ```

    The tool will use the staged diff to generate a commit message and apply it automatically.

---

## Building from Source

If you prefer to build the tool from the source code, follow these steps.

### Prerequisites

- [Go](https://golang.org/doc/install)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Steps

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/commit-craft.git
    cd commit-craft
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Run the tool:**

    ```bash
    go run main.go
    ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

