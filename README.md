# Go WhatsApp Bot

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://go.dev/)

A powerful and extensible WhatsApp bot written in Go, designed for easy interaction and customization.

> **Version Note:**
>
> - If you want to use the version with "case" features, use the code from commit [cc622cd](https://github.com/seaavey/Aemy-go/blob/cc622cd9baefa1748c21a9b7b6c9fe2ab1b5436e) or [v0.2.0-beta](https://github.com/seaavey/Aemy-go/releases/tag/v0.2.0-beta).
> - If you don't need the "case" features, use the latest version from this repository.

## Features

- **WhatsApp Client**: Seamlessly connect and interact with the WhatsApp platform.
- **Message Handling**: Robust system for processing incoming and outgoing messages.
- **Customizable Configuration**: Easily configure the bot's behavior through simple settings.
- **Utility Suite**: Includes helpers for serialization, logging, and other common tasks.

## Prerequisites

Before you begin, ensure you have Go installed on your system.

### Windows (Recommended: Latest)

1.  **Download**: Get the latest Go installer from the [official Go website](https://go.dev/dl/).
2.  **Install**: Run the installer and follow the on-screen instructions.
3.  **Verify**: Open a new terminal (Command Prompt or PowerShell) and run:
    ```cmd
    go version
    ```

### Linux (Recommended: 1.22.0)

1.  **Download the Archive**:
    ```bash
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
    ```
2.  **Extract the Archive**:
    ```bash
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
    ```
3.  **Update PATH**: Add the following line to your `~/.profile` or `~/.bashrc`:
    ```bash
    export PATH=$PATH:/usr/local/go/bin
    ```
    Then, reload your shell profile: `source ~/.profile`.
4.  **Verify Installation**:
    ```bash
    go version
    ```

## Getting Started

1.  **Clone the Repository**:

    ```bash
    git clone https://github.com/seaavey/Aemy-go.git
    cd Aemy-go
    ```

2.  **Install Dependencies**:

    ```bash
    go get
    ```

3.  **Run the Application**:
    ```bash
    go run main.go
    ```

## Deployment (Running 24/7)

For production, it is highly recommended to run the bot on a **Linux** server for better stability, performance, and tooling.

### A Note on Pterodactyl (Not Recommended)

While Pterodactyl is a popular panel for hosting applications, it is **not recommended** for this bot. Pterodactyl is primarily designed for game servers and its file management system can interfere with how this bot saves session data (`session.db`, `qrcode.png`). This can lead to frequent QR code scans and session loss every time the server restarts. For a stable experience, please use one of the recommended methods below.

### Method 1: Linux (Using `screen`)

`screen` is a terminal multiplexer that allows you to run a process in the background and detach from the session, keeping it running even after you log out.

1.  **Install `screen`** (if not already installed):

    ```bash
    # For Debian/Ubuntu
    sudo apt update && sudo apt install screen

    # For CentOS/RHEL
    sudo yum install screen
    ```

2.  **Build the Application**:
    It's better to run a compiled binary for performance.

    ```bash
    go build
    ```

    This will create an executable file (e.g., `Aemy-go`).

3.  **Start a `screen` Session**:

    ```bash
    screen -S bot-session
    ```

4.  **Run the Bot Inside `screen`**:
    Execute the compiled binary.

    ```bash
    ./Aemy-go
    ```

5.  **Detach from the Session**:
    Press `Ctrl + A`, then `D`. The bot is now running in the background.

6.  **Re-attach to the Session** (to view logs or stop the bot):
    ```bash
    screen -r bot-session
    ```

### Method 2: Windows (Using a Service Manager)

On Windows, you can use a tool like **NSSM (the Non-Sucking Service Manager)** to run the bot as a Windows service.

1.  **Build the Application**:

    ```bash
    go build
    ```

    This creates `Aemy-go.exe`.

2.  **Download NSSM**:
    Get the latest release from the [NSSM website](https://nssm.cc/download).

3.  **Install the Service**:

    - Extract `nssm.exe` and place it in a known location (or add it to your `PATH`).
    - Open Command Prompt as an administrator.
    - Run the installer GUI:
      ```cmd
      nssm install AemyBot
      ```
    - In the GUI:
      - **Path**: Browse to your compiled `Aemy-go.exe`.
      - **Startup directory**: Set it to the folder where your bot is located.
      - Click **Install service**.

4.  **Start the Service**:

    ```cmd
    nssm start AemyBot
    ```

    Your bot is now running as a background service. You can manage it using `nssm` commands (`stop`, `restart`, `status`).

## Troubleshooting

### Error: `go-sqlite3 requires cgo`

If you encounter the following error, it means the `go-sqlite3` package requires a C compiler, which may be disabled in your Go environment.

```
[Client ERROR] DB error: failed to upgrade database: failed to check if foreign keys are enabled: Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
```

#### Linux Solution

Install the necessary build tools:

```bash
sudo apt update
sudo apt install build-essential gcc
```

#### Windows Solution

You need to install a GCC toolchain, such as [MinGW-w64](https://www.mingw-w64.org/). After installation, ensure the MinGW `bin` directory is added to your system's `PATH`.

## Author

- **Seaavey** [Seaavey](https://github.com/seaavey)

See also the list of [contributors](https://github.com/seaavey/Aemy-go/contributors) who participated in this project.

## Contact

For questions or support, you can reach out via:

- **Email**: [seaavey@gmail.com](mailto:seaavey@gmail.com)
- **Telegram**: [@seaaveyz](https://t.me/@seaaveyz)
- **WhatsApp Channel**: [Join Channel](https://whatsapp.com/channel/0029Vb3gqjpHltYCf4WOSt3u)

## API

You can also access our REST API at [https://api.seaavey.my.id/](https://api.seaavey.my.id/)

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any feedback or find a bug.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
