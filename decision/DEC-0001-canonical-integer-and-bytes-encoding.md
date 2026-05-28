# DEC-0001: Canonical Integer and Bytes Encoding（规范整数与字节编码）

Status: Accepted

## Context（背景）

Conception 多处使用变长整数、固定宽度整数和长度受限字节序列，但尚未冻结统一字节序、变长整数格式、负数编码和非规范编码的拒绝规则。

## Decision（决策）

- 固定宽度整数统一使用大端序。该规则适用于 `uint16`、`uint32`、`uint64`、`int64` 和固定宽度脚本附参。
- 无符号变长整数使用 ULEB128（等同 Protocol Buffers wire format 中的 unsigned varint）：每字节低 7 位承载数据，高位为续位，低位组先出现。
- 变长整数必须使用最短编码。`0` 编码为 `0x00`，任何带冗余高位零组的编码非法。
- 协议字段不使用有符号变长整数。需要有符号语义时采用固定宽度（如 `int64`），按大端序编码。
- 字节序列默认编码为 `varint(length) || bytes`；固定长度字段不带长度前缀。
- 受限长度字段先按编码长度检查，再按语义长度检查。超过 conception 限制的字段非法。
- **默认编码原则**：协议字段中未显式声明定宽的整数，统一按无符号 varint（ULEB128）编码。当前已明确的定宽字段白名单：
    - 区块头 `Version` (uint32)、`Height` (uint32)、`Stakes` (uint64)
    - 交易头 `Version` (uint16)、`Timestamp` (int64)
    - 哈希字段（`PrevBlock`、`CheckRoot`、`YearBlock`、`TxID` 等）按各自哈希算法的固定字节数
    - 脚本附参中按 opcode 文档声明的定宽位
- **年度字段使用 UTC 自然年份数值**（如 `2025`），不使用相对创世的偏移量。此规则适用于交易输入项短引用、UTXO/UTCO 状态指纹分层等所有按年度划分的协议字段。
- **`BigInt` 二进制序列化**：脚本类型 `BigInt` 在哈希和签名输入中按下列字节布局序列化：

    ```text
    BigIntBytes = 0x01 || sign || magnitude
    ```

    - 第 1 字节为格式版本号，当前固定为 `0x01`。
    - 第 2 字节为符号位：`0x00` 表示零或正数，`0x01` 表示负数。
    - 其后为大端无符号绝对值字节序列。
    - 零值的 `magnitude` 长度为 0；非零值的 `magnitude` 首字节不得为 `0x00`（即必须为最短大端编码）。
    - 此布局与 Go 标准库 `math/big.Int` 的 `GobEncode`/`MarshalBinary` 输出在字节级一致。跨语言实现必须产出相同字节序列。

## Rationale（理由）

大端固定宽度编码便于人工审查和跨语言实现。ULEB128 简单、成熟，且能高效表达输出序位、年度、附件大小等小整数。最短编码规则可避免同一数据拥有多个 TxID 或 BlockID。

默认 varint、定宽字段白名单和"年度采用 UTC 自然年份"源自 conception 修订。`BigInt` 选择 Go 标准库格式可直接复用成熟实现，同时通过显式字节布局让其它语言无需读 Go 源码即可对齐。

## Consequences（影响）

- 所有参与哈希或签名的结构必须使用规范编码。
- 解码器必须拒绝非最短 varint，而不是容忍后再重编码。
- 字段宽度白名单是闭集；后续若需要新增定宽字段，必须在引入该字段的 DEC 中显式声明，并同步更新本 DEC 的白名单。
- `BigInt` 跨语言实现必须提供针对零值、正数、负数、单字节绝对值等边界的字节级测试向量。

## Conception References（构想层依据）

- `docs/conception/附.交易.md`
- `docs/conception/5.信用结构.md`
- `docs/conception/6.脚本系统.md`
- `docs/conception/Instruction/0.基本约束.md`

## Open Questions（开放问题）

（无）
