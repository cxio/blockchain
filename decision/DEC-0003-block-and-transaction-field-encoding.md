# DEC-0003: Block and Transaction Field Encoding（区块与交易字段编码）

Status: Accepted

## Context（背景）

Conception 已给出 `BlockHeader` 与 `TxHeader` 的伪结构，但字段存在可选项、Coinbase 特例和若干宽度待冻结问题。字段编码会直接影响 BlockID、TxID、签名消息和验证路径。

## Decision（决策）

### 区块头编码顺序

1. `Version uint32`
2. `Height uint32`
3. `PrevBlock [48]byte`（SHA3-384 区块 ID）
4. `CheckRoot [48]byte`（SHA3-384）
5. `Stakes uint64`（币权销毁累计值；溢出截断）
6. `YearBlock [48]byte`（仅年块存在，省略规则见下）

`YearBlock` 字段当且仅当 `Height % 87661 == 0` 时存在，否则完全省略（不编码全零占位）。

### 普通交易头编码顺序

1. `Version uint16`
2. `HashInputs [32]byte`（BLAKE3-256）
3. `HashOutputs [32]byte`（BLAKE3-256）
4. `Timestamp int64`（Unix 毫秒；固定 `int64`，不限制负值）
5. `MintPKHash` 可选：编码为 `varint(length) || bytes`

`MintPKHash` 可选字段规则：

- `length == 0`：表示未设置，无后续字节。
- `length == 32`：紧随 32 字节铸凭公钥哈希。
- 其它 `length` 值为非法。

### Coinbase 交易头编码顺序

1. `Version uint16`
2. `HashOutputs [32]byte`（BLAKE3-256）
3. `Timestamp int64`（Unix 毫秒）
4. `MintPKHash [32]byte`（Coinbase 中必填，不使用变长封装）
5. `BlockHeight uint32`
6. `Minter MintProof`（创世 Coinbase 省略，省略规则见下）
7. `FreeData bytes<256>`（`varint(length) || bytes`，length ≤ 255）
8. `BurnCoin int64`（交易费销毁量，单位 chx；语义见 conception `附.交易.md#收益分成`）
9. `AwardSlots [18]byte`（公共服务兑奖槽；编码与确认语义见 DEC-0401）

Coinbase 不含 `HashInputs` 字段（无输入项）。`Minter` 的省略规则：当且仅当 `BlockHeight == 0`（创世）时省略，否则必须存在。无需额外 presence 标识。

## Rationale（理由）

- 字段顺序沿用 conception 伪代码，降低实现与文档偏差。
- 非年块省略 `YearBlock` 可保持区块头常规尺寸较小；年块判定锚定 conception "每年 87661 块"的全局常量。
- `MintPKHash` 在普通交易中采用变长字节序列封装，与 conception "空值时省略"语义直接对应，避免为单一可选字段引入整张 presence 位图。
- Coinbase 中 `MintPKHash` 固定 32 字节，因 conception 已明确"Coinbase 通常会设置此值"，并且 Coinbase 头解析 profile 独立，不与普通交易共用编解码器。
- Coinbase 取消 `HashInputs` 已在 conception 伪代码注释中明示。
- `Minter` 仅在创世区块缺席，使用 `BlockHeight == 0` 作为唯一识别条件可避免额外标志位。
- `BurnCoin` 字段在 conception 伪代码中明列，DEC 原稿遗漏，本次补回。
- `Timestamp` 类型固定 `int64`，覆盖所有合法历史范围；不对负值额外约束，留给上层共识规则裁定。

## Consequences（影响）

- 普通交易头与 Coinbase 交易头必须使用不同解析 profile，避免字段错位。
- `MintPKHash` 普通交易使用变长封装、Coinbase 使用定长 32 字节，二者解析器不共用。
- 区块头常规尺寸：`4+4+48+48+8 = 112` 字节；年块多 48 字节 = 160 字节。
- Coinbase 头不出现 `HashInputs`，对应 TxID 计算时前像中也不出现该字段。
- 冻结后任何字段顺序或 presence 规则的改动均需新增 DEC 并触发协议版本号变更。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#区块头`
- `docs/conception/blockchain.md#创世交易coinbase`
- `docs/conception/附.交易.md#交易头`

## Open Questions（开放问题）

（无）
