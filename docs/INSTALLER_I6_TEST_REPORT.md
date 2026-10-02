# Installer zh-CN — Phase I6 实机验收报告（v1.4.2-zh.3 → v1.4.2-zh.4）

> 测试日期：2026-10-02
> 测试机：Windows 10 22H2（10.0.22621），2560x1440 @ 100% DPI，标准用户权限（非提升）
> 测试对象：**GitHub Release 实际下载件**（非本地构建产物）
> 结论：**P0=0，P1=1（发现并已修复，v1.4.2-zh.4 重验通过）→ 允许进入 I7 Stable**

## 测试版本与资产

| 版本 | 状态 | 说明 |
|---|---|---|
| v1.4.2-zh.3 | 已发布 Pre-release | I6 初测发现 P1 |
| **v1.4.2-zh.4** | 已发布 Pre-release | 含 P1 修复，最终验收对象 |

下载件 SHA256 与 Release 附带 SHA256SUMS.txt **100% 一致**（两个 EXE 均 OK）。

## 发现的问题

### P1（已修复）：Discord 更新残留空壳目录导致安装器完全失效

- **现象**：Discord 自动更新在 `%LOCALAPPDATA%\Discord` 留下 `app-1.0.9260`（resources
  为空、无 build_info）空壳目录。`ParseDiscord` 取"最大 app-* 路径"且不校验 asar 存在，
  于是选定空壳 → 修补时报 "rename app.asar: 系统找不到指定的文件"，Install/Repair/
  Uninstall 全部失败，且 `-location` 指向根目录仍会选中空壳，无绕过手段。GUI 列表
  显示空壳路径且丢失"(已安装 Vencord)"标记。
- **影响**：Discord 更新中断是常见场景；官方 Installer 同样存在此缺陷（上游问题）。
- **修复**：`find_discord_windows.go` — resources 中不存在 `app.asar` 且不存在
  `_app.asar` 的 app-* 目录视为无效（skip）。3 个单测覆盖（空壳跳过/已 patch 优先/
  无 app 目录拒绝）。
- **实机重验**：zh.4 在同一真实环境（空壳仍在）下 Install/Repair/Uninstall 全部成功，
  GUI 列表正确显示唯一真实安装。
- 修复发布为 **v1.4.2-zh.4**（CI 构建，tag 推送自动出包）。

### P3（已修复，下次构建生效）：SHA256SUMS.txt 为 CRLF

Windows runner 上 pwsh `Out-File` 生成 CRLF，GNU `sha256sum -c` 报错。
release.yml 已改为显式 LF 写出（并排除自身重复哈希）。不影响内容正确性，
`tr -d '\r'` 可校验；修复在 I7 Stable 构建中生效。

### 观察项（非缺陷）

- SmartScreen：对带 MotW（模拟下载标记）的未签名 exe 实际拦截启动（"操作已被用户取消"）。
  用户点击"更多信息 → 仍要运行"即可；Release Notes 已说明。未做任何绕过。
- Discord 更新残留：Discord 自身更新失败留下空壳目录属 Discord 更新器行为，
  本修复只保证安装器正确跳过。
- FileVersionInfo API 读版本为空：官方 parity（I5 已对照），非本 fork 问题。

## 测试结果矩阵（对象：v1.4.2-zh.4 下载件，另有标注除外）

| 测试项 | 结果 | 说明 |
| --- | --- | --- |
| Release 资产 SHA256 | PASS | 两个 EXE 与 SHA256SUMS.txt 全部一致 |
| Defender 扫描 | PASS | 下载件 0 威胁（含 zh.3 下载件） |
| SmartScreen | PASS（预期行为） | MotW + 未签名 → 拦截提示，如实记录，未绕过 |
| GUI 启动/中文/CJK | PASS | 警示卡/标题/按钮/列表/tooltip 全中文，无 ? 无 □（截图） |
| GUI patched 状态 | PASS | "Discord（已安装 Vencord）"，空壳目录不再出现（截图） |
| GUI 英文 fallback | PASS | VENCORD_INSTALLER_LANG=en 全英文，**下载/更新源仍为 yepyepos**（DEBUG 日志证实） |
| GUI 窄窗口 620px | PASS | 按钮不重叠、列表完整、警告卡内滚条（与英文版一致行为） |
| GUI 100% DPI | PASS | 当前系统值，布局正常（截图） |
| GUI 125%/150% DPI | NOT TESTED | 需系统级注销生效，无法安全自动化；giu 有运行时缩放（GetContentScale） |
| 窗口最大化 | NOT RUN | 与窄窗同一布局路径（SingleWindow 自适应），低风险 |
| CLI --help 中文 | PASS | flag 名英文、描述中文 |
| CLI 交互菜单 | PASS | 中文/UTF-8 正常（选择器渲染） |
| CLI -version | PASS | `v1.4.2-zh.4 (c3f3bd7)` |
| Stable 检测 | PASS | 真实 Discord（app-1.0.9259）正确识别，空壳跳过 |
| Fresh Install | PASS | 卸载态 → 安装成功，redirector + _app.asar 备份 |
| Repair | PASS | 多次执行无重复 patch |
| Reinstall | PASS | 旧 patch → 新 patch，dist 内容正确 |
| Uninstall | PASS | 原始 app.asar 精确恢复（3614082 B），_app.asar 清除 |
| 卸载后再安装 | PASS | 完整循环 Install→Repair→Reinstall→Uninstall→Install 全绿 |
| 安装后 Discord | PASS | 启动正常，窗口标题"好友 - Discord"（zh Vencord 加载） |
| Partial dist | PASS | 仅 patcher.js（hash 匹配）→ 判定不完整 → 4 文件全量重下 |
| Hash mismatch | PASS | marker deadbee → 判 outdated + WARN"契约校验失败"，不误判最新 |
| Vencord Release 缺 asset | PASS | 缺 renderer.js（缓存过期后）→ 明确报错，0 下载 0 残留 |
| 断网（GitHub 不可达） | PASS | 中文明确失败；0 次官方 URL 请求；现有安装不受影响 |
| Self Update zh.3→zh.4 | PASS | 下载的正是 CI 发布产物，版本递增，.old 自动清理 |
| Up-to-date 无循环 | PASS | zh.4 -update-self → "已是最新版本" |
| 更新失败保旧版 | PASS | 断网 -update-self → 明确失败，旧 exe 完好，0 残留 |
| 标准（非提升）用户权限 | PASS | 本轮全部循环在非管理员 shell 完成 → 无需管理员 |
| 无效/受保护路径 | PASS | 明确中文错误，无静默失败 |
| 磁盘残留 | PASS | resources/VencordData 无 .old/临时文件堆积 |
| PTB | NA | Not tested — installation not present |
| Canary | NA | Not tested — installation not present |
| Development | NA | Not tested — installation not present |
| 多分支并存 | NA | 机器仅 Stable；检测逻辑（分支名→目录映射）经单测覆盖 |
| ServerInfo Owner 四场景 | NEEDS MANUAL | 属 Vencord 主项目功能（非 Installer）；需在 Discord 内人工进入大型/小型服务器右键验证；本阶段以"zh dist SHA256 正确 + 中文 UI 加载"为 Installer 侧验收界 |
| Windows 重启验证 | NOT REQUIRED | 安装/卸载均即时生效，无重启要求 |
| go test | PASS | 22 项全过（含 P1 修复 3 项新测试） |
| go vet | PASS | 0 告警 |
| 日志审查 | PASS | 全部运行仅请求 api.github.com/repos/yepyepos/*；"vencord.dev"仅出现在支持链接文本 |
| 性能 | PASS | 启动/检查/下载即时；断网场景无无限重试 |

## 门禁结论

```text
P0 = 0
P1 = 0（发现的 P1 已修复并以 zh.4 重验）
P2 = 0
P3 = 2（SHA256SUMS CRLF——已修复代码，I7 生效；SmartScreen 提示——未签名预期行为，已文档化）
```

**允许进入 I7 Stable Release。** Stable 建议直接以 v1.4.2-zh.4 的 commit 为基础重新
打 Stable tag 构建（或按 I7 方案），本报告作为发布前验收依据。

---

# I7 Stable 发布验证补充（2026-10-02）

对象：v1.4.2-zh.5（首发 Stable 候选）→ v1.4.2-zh.6（最终 Stable）

## I7 发现并修复的 P1

- **CLI 下载失败假成功**：网络中断导致 Vencord dist 下载失败时，`patch()` 沿用上游
  `return nil`（GUI 有弹窗所以合理），CLI 却打印"✔ 操作成功"——实际什么都没 patch。
  真实网络抖动下复现（同一会话内 3 次）。
- **处置**（按阻塞协议）：zh.5 立即降级 Pre-release → 最小修复（CLI 在 patch 后校验
  installed hash 是否更新，未更新则以中文原因失败；exitFailure 前打印静默错误）→
  发布 v1.4.2-zh.6 Stable。
- **修复实证**：断网/metadata 失败路径正确拒绝；真实网络抖动下 8 次连续下载失败
  全部正确 ❌、目标 asar 全程未损、0 假成功；网络恢复后 install 真实成功
  （redirector + _app.asar）。

## I7 Stable（v1.4.2-zh.6）发布后矩阵

| 项目 | 结果 |
| --- | --- |
| 发布后重新下载 SHA256（LF 直校） | PASS（sha256sum -c 直接 OK） |
| Defender（Stable 下载件） | PASS（0 威胁） |
| -version / --help（下载件） | PASS（v1.4.2-zh.6 (be81596)，中文帮助） |
| GUI Stable 启动 | PASS（中文、已安装标记、无更新弹窗、空壳不出现，截图） |
| 英文 fallback | PASS（I6 已验，代码未变） |
| Fresh install（Stable CLI，真实网络） | PASS（✔ + redirector + _app.asar） |
| Repair | PASS |
| Uninstall（字节级还原） | PASS（sha256 与原始一致） |
| 卸载后重装 | PASS |
| 发布后升级 zh.4 → zh.6 自更新 | PASS（✔ 操作成功，版本 be81596） |
| Already up to date（zh.6） | PASS（"已是最新版本"，无循环） |
| 网络失败保旧版 | PASS（zh.4 自更新 3 次失败均完好回滚） |
| 空壳 app-* 永久回归 | PASS（真实空壳 9260 存在下全部操作正确跳过） |
| 残留检查 | PASS（无 .old/临时文件堆积） |
| go test / go vet | PASS / 0 warnings |
| 125%/150% DPI | NOT TESTED（需系统注销，已在 Stable Release Notes 披露） |
| Administrator 提升 | NOT TESTED（UAC 交互无法自动化；全部功能路径已在标准用户下验证） |

**最终：P0=0，P1=0（I7 P1 已修复重验），P2=0，P3=已知项（SmartScreen/未签名、DPI 未实测）——Stable v1.4.2-zh.6 发布通过。**
