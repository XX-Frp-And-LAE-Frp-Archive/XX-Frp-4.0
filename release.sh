#!/bin/bash

# 清理之前的构建结果
rm -rf build

# 创建输出目录
mkdir -p build

# 定义目标操作系统和体系架构列表
targets=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
)

# 循环构建不同目标平台的可执行文件
for target in "${targets[@]}"; do
    # 解析目标操作系统和体系架构
    os=$(echo "$target" | cut -d'/' -f1)
    arch=$(echo "$target" | cut -d'/' -f2)

    # 设置构建目标
    export GOOS="$os"
    export GOARCH="$arch"

    # 构建可执行文件
    output_name="build/ME-Frp-$os-$arch"
    go build -o "$output_name" cmd/main.go

    # 打印构建完成信息
    echo "已生成 $output_name"
done

# 清理临时环境变量
unset GOOS
unset GOARCH
