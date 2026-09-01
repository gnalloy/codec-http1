# 配置说明

[English](configuration.md) | [文档索引](README.zh-CN.md)

Gnalloy 模块的配置必须显式。优先使用构造参数、option struct 和应用自有配置文件，不使用包级可变状态。

## 主要配置点
- frame、header、body 与 decoded-content 上限必须由服务的可信边界决定。
- 大 payload 应使用 streaming 或 chunked 模式，避免无界聚合正文。
- 压缩模块必须配置解码后大小限制，防御膨胀攻击。
- ByteBuf 所有权遵守 Gnalloy 消息规则：只有当前组件真正消费消息后才释放。

## 推荐默认值

- 从有界大小和较短集成测试超时开始。
- 只有在测量真实 payload 与对端行为后才提高限制。
- 安全相关默认值保持关闭或保守。
- 每个生产覆盖项都应记录在拥有该配置的服务中，而不是写进 library module。

## 环境变量

普通单元测试不要求该 library module 提供仓库专属环境变量。应用、examples、benchmarks 和 CI job 可以围绕它定义自己的变量。

## 本地 Workspace 开发

本地 `go.work` 只作为开发便利。发布用 module metadata 必须保持可移植，不能包含机器相关路径。
