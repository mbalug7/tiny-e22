FROM mcr.microsoft.com/vscode/devcontainers/go:1.19
RUN wget https://github.com/tinygo-org/tinygo/releases/download/v0.26.0/tinygo_0.26.0_amd64.deb
RUN dpkg -i tinygo_0.26.0_amd64.deb
RUN adduser vscode dialout
RUN adduser vscode tty
USER vscode