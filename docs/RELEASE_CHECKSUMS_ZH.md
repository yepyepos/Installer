# Installer zh-CN Release Checksums — v1.4.2-zh.6 (Stable)

> 本文件属于 **yepyepos/Installer**（安装器仓库）。Vencord 主项目的 dist 校验记录在
> 其自己的 `docs/RELEASE_CHECKSUMS.md`，两者独立维护。

## Stable Release

```text
Tag:            v1.4.2-zh.6（Stable，GitHub latest release）
GitHub Release: https://github.com/yepyepos/Installer/releases/tag/v1.4.2-zh.6
Tag 指向 commit: be81596（I7 P1 修复：CLI 下载失败不再误报成功）
```

## 发布资产（CI 构建，run 36999705472，windows-latest + MSYS2）

| 文件 | 大小 | SHA256 |
|---|---|---|
| VencordInstaller.exe | 11853824 B | 见 Release 附带 SHA256SUMS.txt |
| VencordInstallerCli.exe | 8505856 B | 同上 |
| SHA256SUMS.txt | LF 行尾 | `sha256sum -c SHA256SUMS.txt` 直接通过（发布后实测） |

发布后重新下载实测：两个 EXE 校验 OK、Defender 0 威胁、`-version` = `v1.4.2-zh.6 (be81596)`。

## 版本沿革

```text
v1.4.2-zh.1/.2   I4 自更新链测试 Pre-release
v1.4.2-zh.3      I5/I6 RC（首个 CI 构建链）
v1.4.2-zh.4      I6 P1 修复（空壳 app-* 目录）RC，最终验收对象
v1.4.2-zh.5      I7 首发 Stable 候选 → 发现 CLI 下载失败误报成功（P1）→ 已降级 Pre-release
v1.4.2-zh.6      I7 Stable（= zh.5 + P1 修复：CLI 拒绝假成功并打印中文原因）
```

## 构建环境（CI）

```text
Go/编译器:  MSYS2 mingw-w64-x86_64-go / gcc
go-winres:  v0.3.3（官方 pinned d743268）
版本注入:   -X buildinfo.InstallerTag=v1.4.2-zh.6 -X buildinfo.InstallerGitHash=be81596
触发:       push tag v1.4.2-zh.6（release.yml，Windows-only）
校验和:     SHA256SUMS.txt 以 LF 行尾生成（release.yml 内 [IO.File]::WriteAllText）
```

## 签名与杀软（Stable 发布件实测）

- 未签名（官方 Windows 版同样未签名，社区构建身份已在 Release Notes 声明）
- SmartScreen 对带下载标记（MotW）的未签名 exe 实测拦截启动，用户点击
  「更多信息 → 仍要运行」继续；未做任何绕过
- Windows Defender 自定义扫描 Stable 下载件：0 威胁
