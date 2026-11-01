# See here for image contents: https://github.com/devcontainers/images/tree/main/src/go

ARG VARIANT=1.27-trixie
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

RUN \
    apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y bash-completion \
    && echo 'source /usr/share/bash-completion/bash_completion' >> /root/.bashrc \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://claude.ai/install.sh | bash \
    && curl -fsSL https://opencode.ai/install | bash -s -- --no-modify-path

ENV PATH="/root/.local/bin:/root/.opencode/bin:${PATH}"
