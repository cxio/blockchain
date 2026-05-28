# DEC-0002: Domain Tags and Hash Profiles（域标签与哈希配置）

Status: Accepted

## Context（背景）

Conception 已指定区块头、交易头、树枝、附件、公钥哈希和铸凭哈希的算法分配，但尚未冻结域隔离标签、标签编码方式、版本号和各类哈希前像的精确分隔。

## Decision（决策）

### 域标签编码

- 域标签编码为 ASCII 字符串：`"Evidcoin/v1/" || name || 0x00`。
- 域标签必须作为哈希前像首段，不参与用户可控字段的长度解释。
- 域标签作为 ASCII 字符串常量直接进入哈希前像，不压缩为整数编号。
- 同一算法不同用途必须使用不同标签，即使输入结构当前不可混淆。
- 版本号位于标签头部（`v1`），用于未来协议升级时整体替换；当前所有标签均使用 `v1`。

### 域标签清单

链上共识用途的域标签：

- `block.header`
- `tx.header`
- `tree.leaf`
- `tree.branch`
- `checkroot`
- `utxo.leaf`
- `utco.leaf`
- `mint.hash`
- `signature.message`
- `attachment.fingerprint`
- `address.single`
- `address.multi`

说明：交易输出哈希树（`HashOutputs = Hash256(Tree<Outputs>)`）、交易输入哈希树、区块交易哈希树、UTXO/UTCO 中间层均使用通用的 `tree.leaf` / `tree.branch` 域标签，不另设专属标签。

### 算法 profile

沿用 conception `blockchain.md#哈希策略`：

- 区块头、交易头、CheckRoot、UTXO/UTCO 叶：`SHA3-384`。
- 通用哈希树分支（含区块交易树、交易输入树、交易输出树、UTXO/UTCO 中间层）：`BLAKE3-256`。
- 通用哈希树叶（默认）：`SHA3-384`。
- 附件完整指纹：`SHA3-512`。
- 公钥哈希：`SHA3-256(BLAKE2b-512(...))`。
- 铸凭哈希：`BLAKE3-256`。
- BLAKE3 不使用 keyed mode，所有用途统一为普通 hash 加域标签前缀。

### 附件片组哈希树的例外

附件片组哈希树（leaf 与 branch 皆为 `BLAKE3-256`）**不前置任何域标签**，与本 DEC 其它链上哈希用途相区别。理由是该哈希树面向通用文件分享场景，便于与外部数据网络互操作。

- 叶子节点：`BLAKE3-256(2-byte seq || BLAKE3-256(piece_data))`，34 字节预像中 2 字节为顺序号、32 字节为分片数据的哈希。
- 分支节点：`BLAKE3-256(left || right)`。
- 这是全局"哈希前像首段必须为域标签"规则的**唯一例外**。

## Rationale（理由）

域隔离可防止不同结构的相同字节前像产生跨用途混淆。字符串标签比纯整数表更易审计；尾部 `0x00` 分隔符可避免标签名与后续字段拼接歧义。版本号位于头部便于未来整体迁移到 `Evidcoin/v2/...`。

附件片组哈希树例外是 conception 明确的设计取舍：附件数据由外部 P2P 网络承担，使用通用哈希格式可让链外节点直接复用普通文件分片验证工具，不必依赖本协议域标签规范。

## Consequences（影响）

- 任何冻结后的哈希前像都必须列明域标签（除附件片组哈希树例外）。
- 历史版本升级时，应通过 `Evidcoin/vN` 实现整体隔离。
- 附件片组哈希值与协议内其它树根虽算法相同，但因前像无域标签而属于不同名字空间。
- BLAKE3 不使用 keyed mode 意味着同一字节序列在任意上下文中产生相同摘要，跨用途隔离完全依赖域标签。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#哈希策略`
- `docs/conception/附.交易.md`
- `docs/conception/5.信用结构.md#关于附件`
- `docs/conception/1.共识-历史证明（PoH）.md#铸凭哈希`

## Open Questions（开放问题）

（无）
