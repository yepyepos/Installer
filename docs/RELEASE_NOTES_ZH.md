# Release Notes（存档）— Vencord Installer zh-CN v1.4.2-zh.3 RC

> 本文为 GitHub Release notes 的仓库存档，随维护更新。
> Release 页面：https://github.com/yepyepos/Installer/releases/tag/v1.4.2-zh.3

## 身份声明

本项目是 [Vencord Installer](https://github.com/Vencord/Installer) 的**社区中文化 Fork**，
不是官方产品。官方 Installer 由 Vendicated 及 Vencord 贡献者维护；本 fork 由
[yepyepos](https://github.com/yepyepos) 维护，跟随上游更新（同步 → 保留汉化与
yepyepos 源 → 重新验证）。License：GPL-3.0（© 2023 Vendicated and Vencord contributors，
中文版修改部分同样以 GPL-3.0 发布）。

## 版本对应

```text
Installer（本仓库）:  v1.4.2-zh.3（基于官方 v1.4.2；zh.N 为 fork 自身迭代号）
Vencord（主项目）:    v1.15.7-zh.4（独立仓库、独立版本，勿混用）
```

## 本版本内容

- GUI / CLI 全量简体中文（I2）：文案层 localization.go，`VENCORD_INSTALLER_LANG=en`
  可切回英文；CJK 系统字体方案（微软雅黑，不嵌入字体）
- Vencord 下载链（I3）：只认 yepyepos/Vencord，Release 资产契约见
  `docs/VENCORD_RELEASE_CONTRACT.md`；四桌面文件齐全性校验 + Content-Length 校验
- 自更新链（I4）：只认 yepyepos/Installer，版本语义比较，失败自动回滚旧版
- 构建工程（I5）：Windows-only Release workflow（GUI x64 + CLI x86 + SHA256SUMS），
  tag 推送自动构建草稿 Release；winget workflow 已移除（官方渠道不适用 fork）

## 校验

见 `docs/RELEASE_CHECKSUMS_ZH.md` 与 Release 资产中的 SHA256SUMS.txt。

## 已知限制

- 未签名（官方 Windows 版亦未签名），可能触发 SmartScreen / 杀软提示
- Windows 版本信息经 FileVersionInfo API 读取为空——与官方 v1.4.2 一致（go-winres
  0.3.3 既有行为），资源内容实体正确
- macOS / Linux 构建未包含（fork 仅维护 Windows 目标；官方 workflow 中的 mac job
  依赖上游签名机密）
