#!/bin/sh

set -e

DIR="$(CDPATH= cd "$(dirname "${0}")/../../.." && pwd)"
cd "${DIR}"

if [ -z "${1}" ]; then
  echo "usage: ${0} out_dir" >&2
  exit 1
fi

OUT_DIR="${1}"
rm -rf "${OUT_DIR}"
mkdir -p "${OUT_DIR}/bin"
mkdir -p "${OUT_DIR}/etc/bash_completion.d"
mkdir -p "${OUT_DIR}/share/fish/vendor_completions.d"
mkdir -p "${OUT_DIR}/share/zsh/site-functions"

for binary in buf protoc-gen-buf-breaking protoc-gen-buf-lint; do
  echo go build -ldflags \"-s -w\" -trimpath -o \"${OUT_DIR}/bin/${binary}\" \"cmd/${binary}/main.go\"
  go build -ldflags "-s -w" -trimpath -o "${OUT_DIR}/bin/${binary}" "cmd/${binary}/main.go"
done
echo \"${OUT_DIR}/bin/buf\" bash-completion \> \"${OUT_DIR}/etc/bash_completion.d/buf\"
"${OUT_DIR}/bin/buf" bash-completion > "${OUT_DIR}/etc/bash_completion.d/buf"
echo \"${OUT_DIR}/bin/buf\" fish-completion \> \"${OUT_DIR}/share/fish/vendor_completions.d/buf.fish\"
"${OUT_DIR}/bin/buf" fish-completion > "${OUT_DIR}/share/fish/vendor_completions.d/buf.fish"
echo \"${OUT_DIR}/bin/buf\" zsh-completion \> \"${OUT_DIR}/share/zsh/site-functions/_buf\"
"${OUT_DIR}/bin/buf" zsh-completion > "${OUT_DIR}/share/zsh/site-functions/_buf"
echo cp \"LICENSE\" \"${OUT_DIR}/LICENSE\"
cp "LICENSE" "${OUT_DIR}/LICENSE"
