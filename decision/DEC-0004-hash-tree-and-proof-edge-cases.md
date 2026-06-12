# DEC-0004: Hash Tree and Proof Edge Cases（哈希树与证明边界）

Status: Accepted

## Context（背景）

Conception 使用多类哈希树：区块交易树、交易输出树、附件片组树、UTXO/UTCO 状态树。其算法用途已明确，但空树、单叶、奇数节点、叶子序号和验证路径编码尚未统一。

## Decision（决策）

### 通用二叉哈希树规则

适用范围：区块交易树、交易输出哈希树、UTXO/UTCO 中间层等链上共识用途；交易输入根按下方专用规则计算，不套通用二叉树构造。
**例外**：附件片组哈希树按 DEC-0002 的"附件片组哈希树例外"条款执行（不带域标签）。

- 分支节点前像为 `DomainTag("tree.branch") || leftHash || rightHash`。
- 叶子节点前像为 `DomainTag("tree.leaf") || payload`；其中 `payload` 由各具体结构定义内部排布（**序号若存在，作为 payload 内部前缀**，详见下方专用规则）。
- 奇数层最后一个节点直接提升到下一层，不复制自身。
- 单叶树根使用 `tree.branch` profile 归一化为 32 字节：`BLAKE3-256(DomainTag("tree.branch") || leafHash)`；不复制叶子，也不构造 `leafHash || leafHash` 双子分支。
- 空根由各结构自定义；UTXO/UTCO 见 DEC-0201。
- 通用证明路径编码为：
    - `leafHash`（已含序号等 payload）
    - `siblings[]`，每项包含 `direction`（`0`=左 / `1`=右）和 `hash`
    - `rootHash`
- 通用证明路径**不单独携带 `leafIndex` 字段**；若验证方需要按序号定位，应由外部提供含序号的 payload 并重算 `leafHash` 后比对。

### 专用规则

- **区块交易树**：叶子 payload = `seq(3 bytes, BE) || TxID`；叶子哈希 = `SHA3-384(DomainTag("tree.leaf") || seq || TxID)`。3 字节序号是叶子哈希前像的一部分，不作为通用 `leafIndex` 字段。若区块仅含一笔交易，交易树根仍按单叶根规则归一化为 32 字节。
- **附件片组哈希树**：按 DEC-0002 例外条款，叶子哈希 = `BLAKE3-256(2-byte seq || BLAKE3-256(piece_data))`，**不带域标签**；序号同样是叶子前像内部前缀。
- **交易输入根**：继续使用 `BLAKE3-256(ListHash || LeadPKHash)`，不套通用哈希树构造。
- **UTXO/UTCO 指纹**：宽成员状态树，分层结构由 DEC-0201 定义，不直接套通用二叉证明格式；其分支节点仍按 `tree.branch` 域标签编码。

## Rationale（理由）

不复制奇数叶可避免人为重复数据，且兼容 conception 中"宽成员/类 Merkle"两种结构。由于通用树叶使用 `SHA3-384`（48 字节）而分支使用 `BLAKE3-256`（32 字节），若单叶树根直接等于叶哈希，会造成单叶根 48 字节、多叶根 32 字节的不稳定宽度。因此单叶根必须通过 `tree.branch` profile 做一次一元归一化，保持 `Tree<...>` 根与 `Hash256:Tree<...>` 语义下的根宽度一致；这里的 `Hash256:Tree<Outputs>` 表示哈希树根采用 256 位树根 profile。

**空根由各结构自定义**，UTXO/UTCO 见 DEC-0201，避免引入多套空根派生规则。

**序号作为 payload 内部前缀**而非通用 `leafIndex` 字段，保留 conception `blockchain.md#交易约束`、`附.组队校验.md#铸造者验证` 中"叶子节点前置3字节序号"的原文语义，且让该序号实际进入叶子哈希值，便于参与"含序"验证路径。

**附件片组哈希树免域标签**遵循 DEC-0002 的例外条款，保持 `5.信用结构.md` "无需前置域标识（便于通用文件分享）"的设计取舍。

## Consequences（影响）

- 验证路径必须携带方向信息（`direction`），且不再需要单独的 `leafIndex` 字段。
- 通用树根宽度固定为 32 字节；单叶证明的兄弟路径为空，但验证时仍需按单叶根规则将 `leafHash` 归一化后与 `rootHash` 比对。
- 区块交易树与附件片组树虽都使用 2/3 字节序号作为叶子前缀，但因域标签差异、哈希算法不同（SHA3-384 vs BLAKE3-256），属于不同名字空间，无混淆风险。
- 交易输入根、UTXO/UTCO 指纹由 DEC-0101、DEC-0201 各自细化。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#交易约束`
- `docs/conception/blockchain.md#哈希树结构`
- `docs/conception/附.交易.md#哈希校验树`
- `docs/conception/5.信用结构.md#附件id的结构`
- `docs/conception/附.组队校验.md#utxoutco-指纹`

## Open Questions（开放问题）

（无）
