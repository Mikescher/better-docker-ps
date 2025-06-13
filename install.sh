#!/usr/bin/env bash
set -euo pipefail

# This script handles the installation of 'dops' by detecting the OS
# and architecture, downloading the appropriate binary, and configuring the shell PATH.

# Reset
Color_Off=''

# Regular Colors
Red=''
Green=''
Dim='' # White

# Bold
Bold_Green=''
Bold_White=''

if [[ -t 1 ]]; then
    # Reset
    Color_Off='\033[0m' # Text Reset

    # Regular Colors
    Red='\033[0;31m'   # Red
    Green='\033[0;32m' # Green
    Dim='\033[0;2m'    # White

    # Bold
    Bold_Green='\033[1;32m' # Bold Green
    Bold_White='\033[1m'    # Bold White
fi

error() {
    echo -e "${Red}error${Color_Off}:" "$@" >&2
    exit 1
}

info() {
    echo -e "${Dim}$@ ${Color_Off}"
}

success() {
    echo -e "${Green}$@ ${Color_Off}"
}

info_bold() {
    echo -e "${Bold_White}$@ ${Color_Off}"
}

# Check for curl
command -v curl >/dev/null ||
    error 'curl is required to install dops'

REPO="Mikescher/better-docker-ps"
BINARY_NAME=""

# Platform detection
OS=$(uname -s)
ARCH=$(uname -m)

info "Detecting platform: ${OS}/${ARCH}..."

case "$OS" in
    Linux)
        if [ "$ARCH" = "x86_64" ]; then
            BINARY_NAME="dops_linux-amd64-static"
        elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
            BINARY_NAME="dops_linux-arm64-static"
        fi
        ;;
    Darwin)
        if [ "$ARCH" = "arm64" ]; then
            BINARY_NAME="dops_macos-arm64"
        elif [ "$ARCH" = "x86_64" ]; then
            error "Intel-based Macs are not supported."
        fi
        ;;
esac

if [ -z "$BINARY_NAME" ]; then
    error "Unsupported OS or Architecture: ${OS}/${ARCH}"
fi

DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${BINARY_NAME}"

install_env=DOPS_INSTALL
bin_env=\$$install_env/bin

install_dir=${!install_env:-$HOME/.dops}
bin_dir=$install_dir/bin
exe=$bin_dir/dops

if [[ ! -d $bin_dir ]]; then
    mkdir -p "$bin_dir" ||
        error "Failed to create install directory \"$bin_dir\""
fi

info "Downloading 'dops' from ${DOWNLOAD_URL}..."

curl --fail --location --progress-bar --output "$exe" "$DOWNLOAD_URL" ||
    error "Failed to download dops from \"$DOWNLOAD_URL\""

chmod +x "$exe" ||
    error 'Failed to set permissions on dops executable'

tildify() {
    if [[ $1 = $HOME/* ]]; then
        local replacement=\~/
        echo "${1/$HOME\//$replacement}"
    else
        echo "$1"
    fi
}

success "dops was installed successfully to $Bold_Green$(tildify "$exe")"

if command -v dops >/dev/null; then
    echo "Run 'dops --help' to get started"
    exit
fi

refresh_command=''

tilde_bin_dir=$(tildify "$bin_dir")
quoted_install_dir=\"${install_dir//\"/\\\"}\"

if [[ $quoted_install_dir = \"$HOME/* ]]; then
    quoted_install_dir=${quoted_install_dir/$HOME\//\$HOME/}
fi

echo

case $(basename "$SHELL") in
fish)
    commands=(
        "set --export $install_env $quoted_install_dir"
        "set --export PATH $bin_env \$PATH"
    )

    fish_config=$HOME/.config/fish/config.fish
    tilde_fish_config=$(tildify "$fish_config")

    if [[ -w $fish_config ]]; then
        {
            echo -e '\n# dops'
            for command in "${commands[@]}"; do
                echo "$command"
            done
        } >>"$fish_config"
        info "Added \"$tilde_bin_dir\" to \$PATH in \"$tilde_fish_config\""
        refresh_command="source $tilde_fish_config"
    else
        echo "Manually add the directory to $tilde_fish_config (or similar):"
        for command in "${commands[@]}"; do
            info_bold "  $command"
        done
    fi
    ;;
zsh)
    commands=(
        "export $install_env=$quoted_install_dir"
        "export PATH=\"$bin_env:\$PATH\""
    )

    zsh_config=$HOME/.zshrc
    tilde_zsh_config=$(tildify "$zsh_config")

    if [[ -w $zsh_config ]]; then
        {
            echo -e '\n# dops'
            for command in "${commands[@]}"; do
                echo "$command"
            done
        } >>"$zsh_config"
        info "Added \"$tilde_bin_dir\" to \$PATH in \"$tilde_zsh_config\""
        refresh_command="exec $SHELL"
    else
        echo "Manually add the directory to $tilde_zsh_config (or similar):"
        for command in "${commands[@]}"; do
            info_bold "  $command"
        done
    fi
    ;;
bash)
    commands=(
        "export $install_env=$quoted_install_dir"
        "export PATH=\"$bin_env:\$PATH\""
    )

    bash_configs=(
        "$HOME/.bashrc"
        "$HOME/.bash_profile"
    )

    if [[ ${XDG_CONFIG_HOME:-} ]]; then
        bash_configs+=(
            "$XDG_CONFIG_HOME/.bash_profile"
            "$XDG_CONFIG_HOME/.bashrc"
            "$XDG_CONFIG_HOME/bash_profile"
            "$XDG_CONFIG_HOME/bashrc"
        )
    fi

    set_manually=true
    for bash_config in "${bash_configs[@]}"; do
        tilde_bash_config=$(tildify "$bash_config")

        if [[ -w $bash_config ]]; then
            {
                echo -e '\n# dops'
                for command in "${commands[@]}"; do
                    echo "$command"
                done
            } >>"$bash_config"
            info "Added \"$tilde_bin_dir\" to \$PATH in \"$tilde_bash_config\""
            refresh_command="source $bash_config"
            set_manually=false
            break
        fi
    done

    if [[ $set_manually = true ]]; then
        echo "Manually add the directory to your shell configuration file (or similar):"
        for command in "${commands[@]}"; do
            info_bold "  $command"
        done
    fi
    ;;
*)
    echo 'Manually add the directory to your shell configuration file (or similar):'
    info_bold "  export $install_env=$quoted_install_dir"
    info_bold "  export PATH=\"$bin_env:\$PATH\""
    ;;
esac

echo
info "To get started, run:"
echo

if [[ $refresh_command ]]; then
    info_bold "  $refresh_command"
fi

info_bold "  dops --help"

echo
info "To use 'dops' as a drop-in replacement for 'docker ps',"
info "add the following function to your shell configuration file (e.g., ~/.zshrc, ~/.bashrc):"
echo
info_bold 'docker() {'
info_bold '  case $1 in'
info_bold '    ps)'
info_bold '      shift'
info_bold '      command dops "$@"'
info_bold '      ;;'
info_bold '    *)'
info_bold '      command docker "$@";;'
info_bold '  esac'
info_bold '}'
