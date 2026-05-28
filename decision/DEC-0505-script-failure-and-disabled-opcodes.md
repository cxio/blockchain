# DEC-0505: Script Failure and Disabled Opcodes（脚本失败语义与禁用指令）

Status: Accepted

## Context（背景）

Conception 中“失败”“正常失败”“正常退出”“公共验证隐式 END”等术语存在混用；同时 `Instruction/AGENTS.md` 声明若干指令前期禁用，覆盖其它地方的允许约定。Decision 补充执行状态机、禁用规则和公共/私有路径边界。

## Decision（决策）

执行状态分为：

- `Running`：继续执行。
- `PassStop`：正常结束，返回当前 PASS 状态。`PassStop(true)` 表示公共验证通过；`PassStop(false)` 表示验证失败，交易不合法。
- `VerifyFail`：验证失败，交易不合法。
- `ScriptError`：脚本错误，交易不合法。
- `CostFail`：成本超限，交易不合法。
- `PrivateStop`：私有路径停止，不影响公共验证结果。

通关状态语义：

- 初始通关状态为 `true`。
- `PASS` 和 `CHECK` 都写入当前通关状态，后写入者覆盖前值。
- `PASS false` 写入 `false` 后立即产生 `VerifyFail`。
- `PASS true` 写入 `true` 后继续执行。
- `CHECK true/false` 写入对应状态后继续执行。
- `END` 和公共验证路径中的 `INPUT` 以当前通关状态产生 `PassStop`。
- `PassStop(false)` 不等同于 `PASS false` 的立即失败，但最终共识后果同为交易不合法。

错误与停止语义：

- 类型错误、栈下溢、非法 opcode 产生 `ScriptError`。
- 成本超限产生 `CostFail`。
- `INPUT` 在公共验证路径中产生 `PassStop`，保留既有 PASS 状态。
- `END` 在公共验证节点产生 `PassStop`；私有节点可忽略并继续私有路径。
- 公共验证路径遇到 `SYS_TIME` 或 `EXT_PRIV` 产生 `ScriptError`。
- 公共验证路径遇到 `SHELL` 时忽略该指令，但正常消费实参。
- `SHELL` 的“忽略”仅指不执行本地程序、不产生公共副作用；公共实际执行路径触达时仍需按指令规则消费实参、执行栈/类型检查，并按成本模型计入公共成本。
- 公共验证在 `END` 或公共 `INPUT` 后停止；停止后的私有路径产生的 `OUTPUT`、`BUFDUMP`、`PRINT` 或其它私有副作用不记录、不比较、不参与共识。

前期禁用指令清单已由 `docs/conception/Instruction/AGENTS.md` 明确，当前主网清单为：

- `SCRIPT`
- `VALUE`
- `EVAL`
- `INOUT`

执行边界：

- 禁用不是“未实现则忽略”；当前协议有效交易不得依赖禁用指令。
- 禁用指令只在公共验证实际执行路径触达时产生 `ScriptError`。
- 公共 `END` 或公共 `INPUT` 停止后的私有路径可以保留禁用指令，但不得影响公共验证结果。
- 私有工具可解析和显示禁用指令，但必须标注当前禁用状态。
- 若未来解除禁用，应由 conception 或新 DEC 指定协议版本或激活高度。

## Rationale（理由）

失败语义直接决定交易是否有效，必须比指令说明更高层地统一。禁用指令若被不同实现忽略或执行，会导致共识分裂。`SHELL` 只允许产生私有副作用；公共验证忽略并消费实参可保持与 conception 工具指令说明一致，同时避免本地程序执行进入共识。

## Consequences（影响）

- 标准脚本模板必须避开禁用指令。
- 禁用解除需要新的协议版本或明确激活高度。
- 示例文档中使用禁用指令的内容只能作为未来能力说明，不能作为当前有效脚本。
- 实现需要区分 `VerifyFail`、`ScriptError` 和 `CostFail` 的错误码或诊断信息，但三者的共识后果均为交易不合法。

## Conception References（构想层依据）

- `docs/conception/6.脚本系统.md#缓存区和外部监听`
- `docs/conception/6.脚本系统.md#3个特例指令`
- `docs/conception/Instruction/AGENTS.md`
- `docs/conception/Instruction/0.基本约束.md`
- `docs/conception/Instruction/5.交互指令.md`

## Resolution Notes（确认记录）

- 2026-05-28：作者确认主网前期禁用清单按 `Instruction/AGENTS.md` 的 4 项冻结：`SCRIPT`、`VALUE`、`EVAL`、`INOUT`。
- 2026-05-28：作者确认 `CALL` 和 `SHELL` 不属于前期禁用指令。
- 2026-05-28：作者确认 `SHELL` 在公共验证路径中忽略并正常消费实参。
- 2026-05-28：作者确认禁用指令只在公共验证实际执行路径触达时拒绝；公共 `END` 或公共 `INPUT` 后的私有路径可保留禁用指令。
- 2026-05-28：作者确认 `PASS` 和 `CHECK` 均写入通关状态，最后写入者决定最终通关状态；`PASS false` 立即失败。
- 2026-05-28：作者确认私有路径输出不进入公共验证记录。
- 2026-05-28：作者确认成本超限独立为 `CostFail`，交易同样不合法。

## Open Questions（开放问题）

- 禁用解除未来采用逐项 opcode 激活，还是统一脚本版本或协议版本激活。
