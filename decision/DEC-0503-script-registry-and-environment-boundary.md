# DEC-0503: Script Registry and Environment Boundary（脚本注册表与环境边界）

Status: Accepted

## Context（背景）

Conception 定义了环境指令、系统指令、函数指令、模块指令和扩展指令。但环境值编号、函数/模块注册表、公共验证与私有执行边界、外部引用目标缺失语义仍需冻结。

## Decision（决策）

注册表规则：

- opcode 是一级注册表，固定 1 字节。
- `ENV`、`IN`、`OUT`、`VALUE`、`FN`、`MO`、`EXT` 等子空间必须分别维护编号，不共享数值语义。
- 当前 `docs/conception/Instruction/*.md` 已列出的 opcode、附参子空间和值名称是冻结基线。
- 禁用指令仍保留编号，但公共验证实际执行路径触达时拒绝执行。
- 新增注册项只能追加，不得复用已发布编号。

公共验证边界：

- 公共验证路径不得依赖本地时钟、外部输入、外部程序、私有扩展和网络查询的不确定结果。
- `SYS_TIME` 不得位于公共验证路径。
- `INPUT` 在公共验证节点视为隐式结束，不导入外部数据。
- `SHELL` 在公共验证路径中被忽略，但正常消费实参；不得执行本地程序或产生公共验证副作用。
- `EXT_PRIV` 在公共验证路径中非法。
- `GOTO`、`EMBED` 只能引用已确认且可验证的链上脚本。
- `GOTO`、`EMBED` 目标缺失或不可验证时，公共验证归为 `ScriptError`，交易验证失败；不得恢复为 `false` 或忽略。
- 已启用的外部引用指令优先读取用户提供的验证上下文数据，无果之后才会触发网络查询。
- 当前主网公共验证实际执行路径触达 `SCRIPT` 或 `INOUT` 时立即产生 `ScriptError`，不得进入上下文读取或网络查询；`GOTO`、`EMBED` 按验证上下文优先规则处理。

环境命名：

- 区块推导时间使用 `BlockTime`。
- 交易时间戳使用 `TxTime`。
- 当前输入、来源输出、当前输出等环境必须显式区分。

## Rationale（理由）

脚本系统需要同时服务公共验证和私有中间件，必须明确哪些信息会影响共识。注册表分空间可避免环境值和函数编号混淆。

## Consequences（影响）

- 前期实现可只实现公共验证子集。
- 私有功能可以存在于源码和工具层，但不能影响交易合法性。
- `VALUE` 等禁用项解除前，不应进入任何链上标准脚本模板。

## Conception References（构想层依据）

- `docs/conception/6.脚本系统.md#缓存区和外部监听`
- `docs/conception/6.脚本系统.md#3个特例指令`
- `docs/conception/Instruction/1.值指令.md`
- `docs/conception/Instruction/6.结果指令.md`
- `docs/conception/Instruction/13.环境指令.md`
- `docs/conception/Instruction/15.系统指令.md`
- `docs/conception/Instruction/16.函数指令.md`
- `docs/conception/Instruction/17.模块指令.md`
- `docs/conception/Instruction/18.扩展指令.md`

## Resolution Notes（确认记录）

- 各环境值、函数、模块和扩展的编号表以当前 conception 指令文档为冻结基线。
- `GOTO`、`EMBED` 目标缺失或不可验证时按 `ScriptError` 处理。
- 当前主网中，`SCRIPT`、`INOUT` 属于前期禁用指令；公共验证实际执行路径触达时立即失败，不进入上下文读取或网络查询。
- `GOTO`、`EMBED` 优先读取用户提供的验证上下文数据，无果之后才会触发网络查询。
- `SHELL` 在公共验证路径中忽略并消费实参；该裁决与 `DEC-0505` 保持一致。
