# Installer 简体中文本地化（Phase I2）

> 适用版本：Phase I2 起（基于 I3 的 `76fdeb3`/`b8331e5` 之上）。
> I4 增补：Self Updater 重定向与更新策略，见文末"Self Updater（Phase I4 增补）"。
> 范围：GUI / CLI 全部用户可见文本的简体中文化 + Windows GUI 的 CJK 字体支持。

## 翻译架构

轻量静态文案层，不引入任何 i18n 框架：

```text
localization.go        UI struct（约 130 个字段）+ English + Chinese 两份实例
                       pickLanguage() 选择 locale；L 为全局选中实例
localization_test.go   完整性测试（两个 locale 所有字段非空、中文含 CJK、
                       格式占位符可用、环境变量切换）
```

- **默认语言：简体中文**。
- **英文 fallback**：`VENCORD_INSTALLER_LANG=en` 切回英文；两个 locale 字段一一对应且全部非空，任何情况都不会出现空字符串。
- 不做：JSON locale 文件、运行时远程下载、插件化语言系统。
- 复合消息使用 `%s`/`%v` 格式字段（如 `LogPatchingFmt = "正在修补 %s…"`），调用处 `fmt.Sprintf`。

## 术语表（GUI/CLI 统一）

```text
Install → 安装          Uninstall → 卸载        Patch → 修补
Repair → 修复           Reinstall → 重新安装    Update → 更新
Release → Release       Build → 构建
```

## 保持英文/不翻译的内容

- 产品名：Discord / Discord PTB / Discord Canary / Discord Development / Vencord / OpenAsar / GitHub / Windows
- 文件名与路径：patcher.js、renderer.js、app.asar、C:\Users\... 等
- CLI flag 名：-install、-repair 等（**仅翻译其帮助描述**）
- popup/modal id：#patched、#update-prompt 等（imgui 内部标识）
- URL：https://vencord.dev/support 等
- GPL 版权声明与 -version 输出的法律文本
- 纯技术诊断信息：Content-Length 不匹配、HTTP 状态码细节等（错误弹窗采用"中文说明 + 原始技术详情(err.Error())"的结构，排错信息不丢失）
- self_updater.go 整文件未动（I4 范围；其下载 URL 与文本保持原样）

## GUI CJK 字体方案（系统字体，零嵌入）

**根因记录**（I2 实测发现的 giu v0.6.2 行为）：

- giu 的动态字体图集（FontAtlasProsessor）本身支持 CJK：Windows 上默认注册
  Calibri + **MSYH(微软雅黑)** + MSGOTHIC + MALGUNSL，按帧根据实际渲染的字符增量重建。
- 但 `StyleSetter.SetFontSize()` 在未显式 SetFont 时取 `defaultFonts[0]`（= Calibri，
  纯拉丁字体）为每个字号注册独立字体——本应用全部文本都包在 `SetFontSize(...)` 里，
  导致所有文本的 CJK 字形缺失，渲染为 `?`。

**修复**（gui.go，公开 API，3 行）：

```go
if runtime.GOOS == "windows" {
    if _, err := findfont.Find("MSYH"); err == nil {
        g.SetDefaultFont("MSYH", 16)  // 16 = giu 默认字号 14 + 2，与其 Windows 默认注册值一致
    }
}
```

`SetDefaultFont` 将 MSYH 前插到 `defaultFonts[0]`，此后每个 SetFontSize 派生字体均来自
微软雅黑，CJK 正常渲染；拉丁文本也由雅黑的拉丁字形渲染（中文 UI 的预期外观）。

- **不嵌入任何字体文件**（无 20MB 字体、无 License 负担），直接依赖 Windows 系统自带
  的 msyh.ttc（Win8.1+ 均内置）。findfont 找不到时静默跳过（保持 giu 原行为，不会 FATAL）。
- 非 Windows 平台行为完全不变（guard 为 GOOS==windows）。

## CLI 中文控制台

- `cli_windows.go` 新增 `enableUTF8Console()`：`SetConsoleOutputCP/SetConsoleCP(65001)`，
  防止中文输出在旧 conhost（代码页 936）下乱码；`cli_other.go` 提供空实现。
- promptui 菜单、flag 帮助、错误信息全部中文；交互行为（方向键/回车/Esc）不变，
  菜单项字符串同时是选择比较键（全部经 L 字段引用，无散落字面量）。

## 已验证（Windows 真机）

- GUI：中文主界面（警告卡/标题/列表/四按钮/tooltip）、断网错误卡（中文 + 按钮正确禁用）、
  `VENCORD_INSTALLER_LANG=en` 英文回退、640px 窄窗口折行正常、截图确认无 □□□/乱码/缺字。
- CLI：--help 全中文、错误路径全中文、交互菜单全中文、-install/-uninstall 行为不变。
- 测试：11 项单元测试全过（含 6 项本地化完整性测试）；`go vet` 无新增告警
  （仅存官方 `self_updater.go:102` 旧告警，I3 已记录）。
- I3 回归：中文版 CLI/GUI 仍只查询 yepyepos/Vencord，PoC Release 选择、4 asset 下载、
  SHA256、patch 链路全部不变。

## 已知限制

- DPI：本测试机为 100% 缩放（与官方相同 manifest：system dpi-aware）；125%/150% 需真机
  人工验证（giu 通过 GetContentScale + ScaleAllSizes 缩放，理论上正常）。
- 安装器 GUI 日志（stderr）在 GUI 模式不可见，属官方设计，未翻译 DEBUG 级日志。
- 上游已知行为（非 I2 引入）：若某次下载中途失败留下部分 dist 文件（patcher.js 已存在
  且 hash 匹配），下次运行会视为"已是最新"跳过下载。已记录，建议 I5/I6 增强（校验
  4 个文件齐全）。

---

## Self Updater（Phase I4 增补）

### 更新链架构

```text
VencordInstaller(-zh-CN).exe 启动（release 构建，buildinfo.InstallerTag ≠ "Unknown"）
  → GET https://api.github.com/repos/yepyepos/Installer/releases?per_page=20
  → 按版本号选最高 tag（含 pre-release；列表顺序不可依赖，已实测）
  → compareVersions(latest_tag, 本地注入 tag) > 0 才提示更新
  → 用户确认（GUI 弹窗 / CLI --update-self 或菜单）
  → 从该 Release 的 assets 中按名解析 VencordInstaller.exe / VencordInstallerCli.exe
  → 下载 + Content-Length 校验 → 覆盖自身 → 重启 → 仍是中文安装器
```

### 更新策略与版本命名

```text
官方 Installer:      v1.4.2
中文 Fork 基线:      v1.4.2-zh.1 → v1.4.2-zh.2 → …（zh.N 为 fork 自身迭代）
排序规则（version.go compareVersions）:
  v1.4.2 < v1.4.2-zh.1 < v1.4.2-zh.2 < v1.4.10-zh.1 < v1.5.0-zh.1
  - 数字段逐位比较（1.4.10 > 1.4.2，非字符串比较）
  - 带 -zh.N 后缀 > 无后缀（同核心版本时）
  - 无法解析的 tag（devbuild/Unknown/其它后缀）永不参与比较、永不触发更新
```

- 本地版本来自构建时 ldflags 注入的 `buildinfo.InstallerTag`（与官方机制相同）。
- 中文版本号与 Vencord 主项目（v1.15.7-zh.4）相互独立、不混用。

### Asset 命名

沿用官方 asset 名（`VencordInstaller.exe` / `VencordInstallerCli.exe` / …），保证
drop-in 兼容；zh-CN 标识写在 Release 标题/说明中。macOS/Linux 资产暂不发布（I5 处理）。

### 失败回滚

- 下载不完整（Content-Length 不匹配）→ 拒绝替换，临时文件删除，旧 exe 完好。
- 覆盖自身失败（rename 失败）→ 自动把 `.old` 备份恢复原名并提示，旧版继续可用。
- 更新检查/下载网络失败 → 明确报错退出；fallback 与主源同为 yepyepos/Installer，
  **架构上不可能回退到官方 Vencord/Installer 或 vencord.dev**。
- 最新 Release 缺少对应 asset → 明确报错，不触碰现有 exe。

### 与官方实现的差异（仅此四点）

1. `InstallerReleaseUrl/Fallback`：官方 `/releases/latest`（Vencord/Installer）→ fork
   列表端点（yepyepos/Installer）。
2. 版本判断：官方 tag 字符串不等比较 → fork 语义化比较（支持 zh.N 递增与无死循环）。
3. 下载：官方固定 `releases/latest/download/<name>` → fork 从选定 Release 的 assets
   解析（使 pre-release 可选、缺失 asset 有明确错误）。
4. 新增替换失败时的 `.old` 恢复逻辑；官方是删除失败即报错（可能留下无 exe 状态）。

替换、临时文件、重启（RelaunchSelf）、`.old` 清理机制与官方实现完全一致。

### 真机验证记录（v1.4.2-zh.1 → v1.4.2-zh.2）

- CLI `-update-self`：版本由 zh.1 变为 zh.2，`✔ 操作成功！`，重启后 `-version` 确认。
- zh.2 再次 `-update-self` → "无法自更新：已是最新版本"（无更新死循环）。
- GUI zh.1 启动即弹出中文更新提示（"安装器已有新版本！"截图确认）。
- 断网（假代理）更新检查 → 中文明确失败，旧 exe 完好，日志中只有 yepyepos 域名。
- 删除 zh.2 的 CLI asset 后更新 → "最新 Release 中未找到更新文件"，旧版完好。
- 更新后的安装器 I3 回归：仍从 yepyepos/Vencord 下载 4 个桌面文件，SHA256 全部一致。
