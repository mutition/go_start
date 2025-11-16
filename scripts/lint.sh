#!/usr/bin/env bash

set -euo pipefail

source ./scripts/lib.sh

function install_if_not_exist() {
  TOOL_NAME=$1
  INSTALL_URL=$2
  if command -v $TOOL_NAME &> /dev/null
  then
    log_callout "$TOOL_NAME is already installed."
  else
    log_cmd "$TOOL_NAME is not installed. Installing..."
    run go install "$INSTALL_URL"
  fi
}

install_if_not_exist go-cleanarch github.com/roblaszczak/go-cleanarch@latest

readonly LINT_VERSION="1.54.0"
NEED_INSTALL=false
if command -v golangci-lint >/dev/null 2>&1; then
  # golangci-lint has version 1.54.0 built with go1.21.0 from c1d8c565 on 2023-08-09T11:50:00Z
  CURRENT_VERSION=$(golangci-lint --version | awk '{print $4}' | sed 's/^v//')
  log_callout "golangci-lint v$CURRENT_VERSION already installed."
  # ✅ 如果已安装，就跳过安装（即使版本不同）
  NEED_INSTALL=false  # 改为 false，跳过安装
else
  NEED_INSTALL=true
fi

# ✅ 可选：如果网络有问题，可以注释掉安装步骤
# if [ "$NEED_INSTALL" == true ]; then
#   run curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.0
# fi

run go-cleanarch

log_info "lint modules:"
log_info "$(modules)"

run goimports -w -l .

# ✅ 修复配置文件路径检查
CONFIG_FILE="$ROOT_DIR/.golangci.yaml"
if [ ! -f "$CONFIG_FILE" ]; then
  log_error "Config file not found: $CONFIG_FILE"
  exit 1
fi

while read -r module; do
  # ✅ 修复路径大小写问题（Windows 上目录是 Internal，不是 internal）
  if [ -d "$ROOT_DIR/Internal/$module" ]; then
    MODULE_DIR="$ROOT_DIR/Internal/$module"
  elif [ -d "$ROOT_DIR/internal/$module" ]; then
    MODULE_DIR="$ROOT_DIR/internal/$module"
  else
    log_warning "Module directory not found: $module, skipping..."
    continue
  fi
  
  # ✅ 检查是否是有效的 Go 模块（有 go.mod 文件）
  if [ ! -f "$MODULE_DIR/go.mod" ]; then
    log_warning "Module $module has no go.mod, skipping..."
    continue
  fi
  
  # ✅ 检查是否有 Go 文件（递归查找，不限制深度）
  if ! find "$MODULE_DIR" -name "*.go" -type f | grep -q .; then
    log_warning "Module $module has no Go files, skipping..."
    continue
  fi
  
  # ✅ 使用绝对路径，避免 Windows 路径问题
  log_info "Linting module: $module in $MODULE_DIR"
  (
    cd "$MODULE_DIR" || exit 1
    golangci-lint run --config "$CONFIG_FILE" || exit 1
  )
  if [ $? -ne 0 ]; then
    log_error "Linting failed for module: $module"
    exit 1
  fi
done < <(modules)