# Vencord zh-CN Release Contract（Installer ↔ Vencord 仓库接口契约）

> 适用版本：Phase I3 起。
> 本文定义 `yepyepos/Installer`（中文安装器）从 `yepyepos/Vencord` 获取 Desktop 构建产物的完整接口。
> 两个仓库独立维护，本契约是二者之间唯一的耦合点。

## Repository

```text
Installer 查询的仓库（唯一）:  https://github.com/yepyepos/Vencord
API 端点（主源 = 回退源）:     https://api.github.com/repos/yepyepos/Vencord/releases?per_page=20
```

- 安装器**永远不会**查询 `Vendicated/Vencord` 或 `vencord.dev`。
- 主源失败（网络/限流）→ 明确报错并停止，**不得静默改装官方英文 Vencord**。回退端点与主源相同（仅保留上游代码结构的重试语义）。

## Desktop assets

一个 Release 只有在**同时包含以下 4 个资产**（按文件名**前缀**匹配，因此 `*.map` 兄弟文件可共存）时才被视为合格的 Desktop Vencord Release：

```text
patcher.js
preload.js
renderer.js
renderer.css
```

- 下载目标：`<VencordData>/dist/`（Windows 默认 `%APPDATA%\Vencord\dist`，可用 `VENCORD_USER_DATA_DIR` 覆盖）。
- 下载校验：`Content-Length` 一致 + 至少成功 4 个文件；没有 checksum/签名校验（后续安全增强项，见文末）。

## Release 选择规则

```text
在 releases 列表（含 pre-release，不含 draft）中，
按 published_at 降序，选择第一个包含全部 4 个 Desktop assets 的 release。
列表返回顺序不可依赖（created_at 可能反映 tag 目标 commit 时间）。
```

- 纯浏览器扩展 Release（只有 `Vencord.user.js` / extension zip）**不合格**，会被跳过。
- 因此：**Stable Release 若要被安装器选中，必须附带 4 个 Desktop assets。**

## 版本识别 / Build hash

```text
Release 侧:   release body 中的一行机器可读标记（不区分大小写）:
              Vencord-Desktop-Hash: <hash>
              例: Vencord-Desktop-Hash: 4c73c063

已安装侧:     <VencordData>/dist/patcher.js 首行:
              // Vencord <hash>

判定:         A == B → 已是最新（安装时跳过下载）
              A != B → 视为过期 → 重新下载
```

- hash 来源 = Vencord 构建时嵌入 dist 文件的 git commit（`yepyepos/Vencord` 仓库构建 commit 短 hash）。
- **禁止**依赖 Release 标题中偶然出现的数字作为版本判断。
- 下载完成后安装器会校验落盘 `patcher.js` 首行与 Release 声明一致，不一致输出
  `Release contract violation` 警告（metadata 被篡改/忘记更新标记的防线）。
- 兼容性：body 无标记时回退到上游格式（Release 名称末段 = hash，如 `DevBuild <hash>`），
  保证上游 merge 冲突最小化。

## 构建侧义务（yepyepos/Vencord 发布流程）

1. Desktop dist 构建产物 `patcher.js` / `preload.js` / `renderer.js` / `renderer.css` 必须作为 assets 附到 Release。
2. Release body 必须包含 `Vencord-Desktop-Hash: <构建 commit 短hash>`，且与 dist 内嵌 hash 一致。
3. 发布后可用 `docs/RELEASE_CHECKSUMS.md` 记录 SHA256（现状：v1.15.7-zh.4 的 4 个 Desktop asset
   SHA256 与该文件记录一致）。

## Installer 侧现状（Phase I3 PoC）

- `constants.go`: `ReleaseUrl` / `ReleaseUrlFallback` → yepyepos/Vencord 列表端点。
- `github_downloader.go`: `GetLatestVencordRelease`（列表 + 选择）、`ExtractDesktopHash`
  （marker 优先、名称回退）、`HasAllDesktopAssets`（4 资产前缀匹配）、下载后 contract 校验。
- ~~Self Updater 仍指官方仓库~~（Phase I4 已完成重定向）：安装器自更新来源为
  **yepyepos/Installer** 列表端点，按版本号选最高 tag，从该 Release 的 assets 按
  官方同名（`VencordInstaller.exe` / `VencordInstallerCli.exe`）解析下载地址；
  fallback 与主源同端点，不存在任何指向官方仓库的回退路径。
  详见 `docs/INSTALLER_LOCALIZATION_ZH_CN.md` 的"Self Updater"章节与
  `self_updater.go` / `version.go`。

## 尚未实现（记录在案）

- Asset SHA256 / 签名校验（安装器侧安全增强）。
- "安装官方 Vencord" 之类的显式用户选项（默认路径永远是中文）。
- Self Updater 重定向（I4）。
