# DEC-0301: PoH Mint Hash and Mint Proof（PoH 铸凭哈希与择优凭证）

Status: Accepted

## Context（背景）

Conception 已明确铸凭交易区间、评参区块、`Stakes` 来源、择优池大小和铸凭哈希因子。本决策冻结 challenge seed 组成、Equi-X 解答约束、最终 `MintHash` 计算、MintProof 字段顺序和初段 Coinbase 参与资格。

## Decision（决策）

铸凭哈希采用两段式计算：先生成 challenge seed，再通过 Equi-X 生成 `hashList`，最后计算 `MintHash`。

固定参数：

- `Mix = 0x517cc1b727220a95`
- `BlockHeight` 为当前待铸区块高度。
- `X = BE(minimal_unsigned(BlockHeight * Mix))`。

challenge seed 为：

```text
ChallengeSeed = BLAKE3-256(MintPubKey || MintTxID || Stakes || RefMintHash || X)
```

其中：

- `MintPubKey` 为铸造者公钥的 canonical bytes。
- `MintTxID` 为铸凭交易完整 48 字节 TxID。
- `Stakes` 为链末端 `-32` 区块头中的币权销毁值，uint64 按大端序编码。
- `RefMintHash` 为评参区块 Coinbase 中记录的铸凭哈希；创世块或初段无 `Minter` 时取全零 32 字节。

Equi-X 规则：

- 求解输入固定为 `ChallengeSeed`。
- `Nonce` 必须满足 `Nonce >= BlockHeight`。
- `Solution` 在解析为索引序列后必须严格升序且无重复。
- 验证端必须对 `ChallengeSeed`、`Nonce`、`Solution` 运行 Equi-X 验证，且仅在验证通过时接受返回的 `HashList`。

最终铸凭哈希为：

```text
MintHash = BLAKE3-256(DomainTag("mint.hash") || HashList[0] || ... || HashList[n-1])
```

`HashList` 的拼接顺序必须与 `Solution` 索引顺序一致，不允许重排。

`MintProof` 字段顺序为：

1. `TxHeight uint32`
2. `TxID [48]byte`
3. `Nonce uint64`
4. `Solution bytes`（Equi-X 解答的 canonical bytes）
5. `MintPubKey bytes`
6. `MintHash [32]byte`
7. `Signature bytes`

排序规则：

- 先按 `Nonce` 升序，值小者优。
- 在 `Nonce` 相同条件下按 `MintHash` 的 32 字节字典序升序，值小者优。
- 若 `Nonce` 与 `MintHash` 仍相同，则按完整 `TxID` 升序，再按公钥字节升序。

铸造者身份规则：

- 当铸凭交易头包含 `MintPKHash` 时，`MintPubKey` 的公钥哈希必须等于 `MintPKHash`，但不要求该公钥参与输入根。
- 当铸凭交易头不包含 `MintPKHash` 时，`MintPubKey` 的公钥哈希作为 `LeadPKHash`，必须参与输入根验证。
- Coinbase 交易没有输入项，因此没有 `LeadPKHash`；第 240 块之前，只要 Coinbase 设置了 `MintPKHash`，即可作为已确认铸凭交易参与铸造竞争。

## Rationale（理由）

两段式设计把身份和历史因子先压缩为 `ChallengeSeed`，再引入 Equi-X 工作量以抑制实时发掘与算力偏置；最终 `MintHash` 只对有效 `HashList` 聚合并加域标签，保证用途隔离。`X` 使用无损大整数编码，可避免固定宽度溢出。排序先比较 `Nonce` 可抑制“反复试 nonce”带来的竞争偏移。

## Consequences（影响）

- `MintPubKey` 必须能证明其哈希等于 `MintPKHash`，或在无 `MintPKHash` 时作为 `LeadPKHash` 参与输入根验证。
- `MintProof` 验证包含 `Nonce >= BlockHeight`、`Solution` 升序且无重复和 Equi-X 有效性检查。
- 初段 Coinbase 作为铸凭交易时必须设置 `MintPKHash`，并且必须已经确认。

## Conception References（构想层依据）

- `docs/conception/1.共识-历史证明（PoH）.md#规则`
- `docs/conception/1.共识-历史证明（PoH）.md#铸凭哈希`
- `docs/conception/1.共识-历史证明（PoH）.md#择优凭证`
- `docs/conception/blockchain.md#启动`

## Open Questions（开放问题）

（无）
