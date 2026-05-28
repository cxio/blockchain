# DEC-0502: Script Float Profile（脚本浮点配置）

Status: Accepted

## Context（背景）

Conception 允许脚本使用 `Float`，并说明字面量不支持 NaN/Inf，但运算中可能产生异常浮点，可用 `ISEFV` 检测。浮点跨实现确定性、字节编码、取整、比较和异常传播需要冻结。

## Decision（决策）

建议 profile：

- `Float` 使用 IEEE 754 binary64。
- 字节编码使用 8 字节大端 bit pattern。
- 输入字面量不得表达 NaN、+Inf、-Inf。
- `POW`、除零、溢出等运算产生 NaN 或 Inf 时不立即崩溃，保留 IEEE 754 异常值继续执行，由 `ISEFV` 检测。
- `-0.0` 在数值比较中等于 `+0.0`，但字节编码保持原 bit pattern。
- `Float -> Int` 默认向零截断。
- `Float -> String` 默认使用最短 round-trip 十进制格式，不依赖本地 locale；等价于 Go `strconv.FormatFloat(f, 'g', -1, 64)`。
- 若脚本显式使用 `STRING{e/E/f/g/G/x/X}` 格式标识，则按该格式执行。
- `BYTES` 与 `PACK` 对 `Float` 输出 IEEE 754 binary64 的 8 字节大端 bit pattern。
- 异常浮点也允许被 `BYTES` 或 `PACK` 输出为其 IEEE bit pattern；转换后值类型为 `Bytes`，不再触发最终异常 `Float` 残留检查。

建议比较规则：

- 任一操作数为 NaN 时，除 `ISEFV` 外的比较返回 `false`。
- `EQUAL(+0.0, -0.0)` 返回 `true`。
- 排序类比较中遇到 NaN 会导致脚本执行失败，验证不通过。

## Rationale（理由）

完全禁止浮点会削弱脚本表达力；保留 binary64 并冻结异常语义，可在可用性和确定性之间折中。异常浮点不应静默通过公共验证，但脚本可以显式检测并转换为其它类型。最短 round-trip 字符串格式能保证默认文本表示可无损解析回同一 binary64，同时避免固定小数或固定科学计数法的冗长和边界问题。

## Consequences（影响）

- 需要跨平台测试向量覆盖 NaN、Inf、-0.0、舍入和字符串格式。
- 共识实现不得使用会受 CPU 或编译器 fast-math 影响的非标准行为。
- 高精度整数逻辑应使用 `Int` 或 `BigInt`，不要用 `Float` 表达金额。

## Conception References（构想层依据）

- `docs/conception/6.脚本系统.md`
- `docs/conception/Instruction/8.转换指令.md`
- `docs/conception/Instruction/9.运算指令.md`
- `docs/conception/Instruction/10.比较指令.md`

## Open Questions（开放问题）

None.
