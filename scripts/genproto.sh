#!/usr/bin/env bash

set -euo pipefail

#递归匹配子目录。
shopt -s globstar

if ! [[ "$0" =~ genproto.sh ]]; then
    echo "$0"
    echo "Error: This script must be run from the root directory"
    exit 1
fi

source ./scripts/lib.sh

API_ROOT="./api"

function dirs {
    dirs=()
    while IFS= read -r dir; do
        dirs+=("$dir")
    done < <(find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u)
    echo "${dirs[@]}"
}


function proto_files {
    pb_files=$(find . -type f -name "*.proto")
    echo "${pb_files[@]}"
    
}

function gen_for_modules() {
    local go_out="internal/common/genproto"
    if [[ ! -d "$go_out" ]]; then
        log_warn "genproto: creating directory $go_out"
        run rm -rf "$go_out"
    fi
    for dir in $(dirs); do
        local service="${dir:0:${#dir}-2}"
        local pb_file="${service}.proto"
    
        if [ -d "$go_out/$dir" ]; then
            log_warning "cleaning all files under $go_out/$dir"
            run rm -rf "$go_out/$dir"/*
            else
            run mkdir -p "$go_out/$dir"
            fi
            log_info "generating code for $service"

            PROTOC_INCLUDE=""
                if [ -d "$(go env GOPATH)/src" ]; then
                    PROTOC_INCLUDE="-I=$(go env GOPATH)/src"
                elif [ -d "$(go env GOPATH)/pkg/mod" ]; then
                    # 尝试从 Go 模块中查找
                    PROTOC_INCLUDE="-I=$(go env GOPATH)/pkg/mod"
                fi

            run protoc \
                ${PROTOC_INCLUDE} \
                -I="${API_ROOT}" \
                "--go_out=${go_out}" --go_opt=paths=source_relative \
                --go-grpc_opt=require_unimplemented_servers=false \
                "--go-grpc_out=internal/common/genproto" --go-grpc_opt=paths=source_relative \
                "${API_ROOT}/${dir}/$pb_file"
    done
}

echo "directories containting protos to be built: $(dirs)"
echo "found pb_files: $(proto_files)"
gen_for_modules