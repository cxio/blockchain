# DEC-0301: PoH Mint Hash and Mint Proof（PoH 铸凭哈希与择优凭证）

Status: Accepted

## Context（背景）

Conception 已明确铸凭交易区间、评参区块、`Stakes` 来源、择优池大小和铸凭哈希因子。本决策冻结铸凭哈希前像字段编码、`X` 的数值规则、MintProof 字段顺序和初段 Coinbase 参与资格。

## Decision（决策）

铸凭哈希前像为：

```text
DomainTag("mint.hash") || MintPubKey || MintTxID || Stakes || RefMintHash || X
```

其中：

- `MintPubKey` 为铸造者公钥的 canonical bytes。
- `MintTxID` 为铸凭交易完整 48 字节 TxID。
- `Stakes` 为链末端 `-32` 区块头中的币权销毁值，uint64 按大端序编码。
- `RefMintHash` 为评参区块 Coinbase 中记录的铸凭哈希；创世块或初段无 `Minter` 时取全零 32 字节。
- `X = BE(minimal_unsigned(BlockHeight * Mix))`。
- `BlockHeight` 为当前待铸区块高度。
- `Mix = 0x517cc1b727220a95`。

`MintProof` 字段顺序为：

1. `TxHeight uint32`
2. `TxID [48]byte`
3. `MintPubKey bytes`
4. `MintHash [32]byte`
5. `Signature bytes`

排序规则：

- 铸凭哈希按 32 字节字典序升序，值小者优。
- 哈希相同则按完整 TxID 升序，再按公钥字节升序。

铸造者身份规则：

- 当铸凭交易头包含 `MintPKHash` 时，`MintPubKey` 的公钥哈希必须等于 `MintPKHash`，但不要求该公钥参与输入根。
- 当铸凭交易头不包含 `MintPKHash` 时，`MintPubKey` 的公钥哈希作为 `LeadPKHash`，必须参与输入根验证。
- Coinbase 交易没有输入项，因此没有 `LeadPKHash`；第 240 块之前，只要 Coinbase 设置了 `MintPKHash`，即可作为已确认铸凭交易参与铸造竞争。

## Rationale（理由）

前像字段均为公开且可验证数据。`X` 使用无损大整数编码，可避免固定宽度溢出。`MintProof` 把 `MintHash` 放在签名前，可便于检索和预筛选，但签名验证仍以重新计算值为准。

## Consequences（影响）

- `MintPubKey` 必须能证明其哈希等于 `MintPKHash`，或在无 `MintPKHash` 时作为 `LeadPKHash` 参与输入根验证。
- 若 `Stakes=0`，`X` 编码为单字节 `0x00`。
- 初段 Coinbase 作为铸凭交易时必须设置 `MintPKHash`，并且必须已经确认。

## Conception References（构想层依据）

- `docs/conception/1.共识-历史证明（PoH）.md#规则`
- `docs/conception/1.共识-历史证明（PoH）.md#铸凭哈希`
- `docs/conception/1.共识-历史证明（PoH）.md#择优凭证`
- `docs/conception/blockchain.md#启动`

## Open Questions（开放问题）

（无）
