# Installer zh-CN Release Checksums — v1.4.2-zh.3 (RC)

> 本文件属于 **yepyepos/Installer**（安装器仓库）。
> Vencord 主项目（yepyepos/Vencord）的 dist 校验记录在其自己的 `docs/RELEASE_CHECKSUMS.md`，两者独立维护。

## Release

```text
Tag:            v1.4.2-zh.3（Pre-release / Release Candidate）
GitHub Release: https://github.com/yepyepos/Installer/releases/tag/v1.4.2-zh.3
Tag 指向 commit: 6e9688a076e15da86f5a64d01334332620055530
```

注：本文档所在提交晚于 tag commit（发布后补记 checksum 属正常流程）；
tag 与两个 EXE 的 build commit 均为 `6e9688a`（`-version` 输出可证）。

## CI 构建产物（正式资产）

构建于 GitHub Actions `Release` workflow（run 36842319537），windows-latest + MSYS2：

| 文件 | 大小 | SHA256 |
|---|---|---|
| VencordInstaller.exe | 11,853,312 B | `5ee2c5e3f9cd7a2ea708dab148f9d782d4df9ad12fa15ecbc2ff7acfd1bcd74d` |
| VencordInstallerCli.exe | 8,504,832 B | `2f64c9cd1318de820d38a8b8f0b5b0a07ac958f795e466880068d39c69088214` |
| SHA256SUMS.txt | 179 B | （即本表来源，随 Release 附带） |

- VencordInstaller.exe：GUI，windows/amd64，`-H=windowsgui`
- VencordInstallerCli.exe：CLI，windows/386（官方工作流的 CLI 目标架构）

## 构建环境（CI）

```text
Go:            MSYS2 mingw-w64-x86_64-go（runner 上当时版本）
编译器:        MSYS2 mingw-w64-x86_64-gcc（CGO，GUI 需要）
go-winres:     github.com/tc-hib/go-winres@d743268d7ea168077ddd443c4240562d4f5e8c3e (v0.3.3，官方 pinned)
版本注入:      -X buildinfo.InstallerTag=v1.4.2-zh.3
               -X buildinfo.InstallerGitHash=6e9688a
资源生成:      go-winres make --product-version "git-tag"
触发:          push tag v1.4.2-zh.3（release.yml）
```

## 本地对照构建（非发布资产）

同 commit 本地复现构建亦验证通过（供复现参考，不随 Release 发布）：

```text
Go:        go1.27.1 windows/amd64 (scoop)
MinGW:     x86_64-posix-seh-rev1, MinGW-Builds gcc 16.2.0 (scoop mingw)
make:      GNU Make 4.4.1
命令:      make GUI=1 / make（Makefile 官方 target，tag 本地存在时自动取版本）
产物:      VencordInstaller.exe 12,279,296 B（amd64）/ VencordInstallerCli.exe 9,406,976 B（amd64）
说明:      与 CI 产物哈希不同（工具链差异），行为一致；CI 产物为发布来源
```

## 版本资源说明

go-winres v0.3.3（官方 pinned 版本）生成的 VS_VERSIONINFO 经 Windows
`FileVersionInfo` API 读取为空——**官方 v1.4.2 VencordInstaller.exe 同样如此**，
属构建工具链既有行为，本 fork 与官方保持一致，未做偏离性"修复"。
资源实体存在且内容正确（.rsrc 内 UTF-16 字符串可检索：ProductName、中文描述、
GPL 版权、zh-CN Community 标识）；图标资源正常（任务栏/标题栏可见）。

## 签名与杀软

- 官方 Windows 版无 Authenticode 签名，本 fork 同样未签名（不伪装官方签名身份）
- Windows Defender 自定义扫描两个 EXE：0 威胁
- SmartScreen 对未签名 exe 可能有提示；用户可校验上方 SHA256
