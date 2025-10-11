# Commit Craft

Commit Craft is a Go-based tool that leverages AI to automatically generate descriptive commit messages from your staged changes. It streamlines your git workflow by creating conventional commit messages for you.

<https://github.com/user-attachments/assets/d68fd98e-cef1-4827-987f-79cd2c36a438>

## Features

- **AI-Powered Commit Messages:** Automatically generates descriptive commit messages from your staged changes using AI.
- **Interactive UI:** Allows you to accept, modify, or quit the commit process.
- **Multiple AI Providers:** Supports both Gemini and Ollama for commit message generation.
- **Conventional Commits:** Generates messages that follow the conventional commit format.

## Prerequisites

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) must be installed.
- You need an API key for your chosen AI provider (e.g., Gemini).

## Installation

1. **Download the binary:**
   Download the latest binary for your operating system from the [Releases page](https://github.com/starrick2001/commit-craft/releases).

   ```bash
   curl -O https://github.com/Starrick2001/commit-craft/releases/latest/commit-craft
   ```

2. **Make it executable:**

   ```bash
   chmod +x commit-craft
   ```

3. **Move to your PATH (Optional):**
   For easy access, move the binary to a directory in your system's PATH.

   ```bash
   sudo mv commit-craft /usr/local/bin/
   ```

### Arch Linux (from source)
Alternatively, Arch Linux users can build and install the package using the provided PKGBUILD file:
```bash
git clone https://github.com/Starrick2001/commit-craft.git
cd commit-craft
makepkg -si
```

## Configuration

Commit Craft supports both Gemini and Ollama as AI providers. You can configure the provider and API key using environment variables.

### Gemini

To use Gemini, set the following environment variables:

```bash
export COMMIT_CRAFT_PROVIDER='gemini'
export COMMIT_CRAFT_GEMINI_KEY='YOUR_API_KEY'
```

### Ollama

To use Ollama, set the following environment variables:

```bash
export COMMIT_CRAFT_PROVIDER='ollama'
export COMMIT_CRAFT_OLLAMA_MODEL='your-ollama-model'
export COMMIT_CRAFT_OLLAMA_HOST='http://localhost:11434' # Optional: Defaults to this value
```

> **Note:** To make this permanent, add the export command to your shell's configuration file (e.g., `~/.bashrc`, `~/.zshrc`).

## Usage

1. **Stage your changes:**
   Use `git add` to stage the files you want to commit.

   ```bash
   git add .
   ```

2. **Run Commit Craft:**
   Execute the program to generate a commit message.

   ```bash
   commit-craft
   ```

3. **Choose an action:**
   The tool will display the generated commit message and prompt you to:
   - **Commit:** Applies the generated message and creates the commit.
   - **Modify:** Allows you to edit the message before committing.
   - **Quit:** Exits the program without creating a commit.

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

3. **Run the tool:

   ```bash
   go run main.go
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
