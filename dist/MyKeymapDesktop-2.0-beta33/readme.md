[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/xianyukang/MyKeymap/blob/master/readme.en.md)

# MyKeymap-re

MyKeymap-re 是基于 [MyKeymap](https://github.com/xianyukang/MyKeymap) 的桌面化分支，继续保留原有 AutoHotkey 键位映射能力，同时补上更适合长期使用的 Windows 桌面端宿主、后台驻留和配置热重载体验。

## MyKeymap-re 说明

- 保留原版 MyKeymap 的核心映射逻辑、配置格式和 AutoHotkey 脚本链路
- 新增桌面宿主，关闭窗口后可继续常驻后台
- 修改配置后支持保存 / 应用，并可直接热重载生效
- 配置界面升级为更紧凑的桌面端布局，减少浏览器式配置页的割裂感
- 补充进程管理、启动项控制和更一致的图标与品牌资源

## Features

- **程序启动切换**: 用快捷键启动程序和切换窗口，对比搜索型的启动器，效率更高
- **键盘控制鼠标**: 用键盘控制鼠标，能减少键鼠切换，不必为了点一下而大幅移动手掌
- **按键重新映射**:
  - 把常用按键重映射到主键区，提升输入速度和编辑效率
  - 内置多套键位模式，例如光标控制、数字输入、符号输入
- **桌面端管理**:
  - 支持后台驻留、窗口隐藏、托盘恢复
  - 支持配置保存、应用和引擎热重载

## Usage

- [快速入门](https://xianyukang.com/MyKeymap.html#mykeymap-%E7%AE%80%E4%BB%8B) & [视频介绍](https://www.bilibili.com/video/BV1Sf4y1c7p8)
- 桌面端入口: `MyKeymapDesktop.exe`
- 映射引擎入口: `MyKeymap.exe`
- 默认配置文件: `data/config.json`

## Fork Notes

本仓库当前优先维护 `MyKeymap-re` 的 Windows 桌面版体验，不追求改写原始 AHK 核心，而是围绕以下方向持续迭代：

- 保持对原有配置的兼容
- 优化后台驻留和系统进程管理
- 提升保存 / 应用 / 热重载链路的稳定性
- 逐步升级 UI，但避免大面积破坏原有功能

## Screenshots

| ![features](./doc/features.png) | ![settings](./doc/settings.png) |
| ------------------------------- | -------------------------------- |
