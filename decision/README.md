# Decision Index（架构决策索引）

`docs/decision/` 只记录 `docs/conception/` 尚未直接冻结、但会影响实现一致性的补充决策。`docs/conception/` 是上游正式依据；若 Decision 与 conception 冲突，以 conception 为准并修订 Decision。

本目录已按 2026-05-21 的 conception 修订精简重写。旧编号不保留兼容；已被 conception 明确的规则不再单独成篇，历史占位和废弃设计已移除。

## Status Definition（状态定义）

| Status | 含义 |
|--------|------|
| Proposed | 建议性决策。可用于实现草案，但冻结区块/交易 ID 前需要作者确认。 |
| Accepted | 已确认可作为实现依据。只有 conception 明确吸收或作者确认后才能使用。 |
| Deprecated | 已废弃，仅保留历史说明。当前目录不保留独立废弃 DEC。 |
| Absorbed | 已被 conception 直接明确，不再保留独立 DEC。 |

> **本轮重写策略：** 旧状态清零。除非 conception 已直接裁决，否则新 DEC 初始均为 `Proposed`。

## Current Index（当前索引）

| DEC | Title | Status | 主题 |
|-----|-------|--------|------|
| [DEC-0001](DEC-0001-canonical-integer-and-bytes-encoding.md) | Canonical Integer and Bytes Encoding | Accepted | 基础编码 |
| [DEC-0002](DEC-0002-domain-tags-and-hash-profiles.md) | Domain Tags and Hash Profiles | Accepted | 哈希 |
| [DEC-0003](DEC-0003-block-and-transaction-field-encoding.md) | Block and Transaction Field Encoding | Accepted | 字段编码 |
| [DEC-0004](DEC-0004-hash-tree-and-proof-edge-cases.md) | Hash Tree and Proof Edge Cases | Accepted | 哈希树 |
| [DEC-0101](DEC-0101-transaction-body-and-output-payloads.md) | Transaction Body and Output Payloads | Accepted | 交易 |
| [DEC-0102](DEC-0102-signature-message-profile.md) | Signature Message Profile | Accepted | 签名 |
| [DEC-0103](DEC-0103-witness-container-and-pruning.md) | Witness Container and Pruning | Accepted | 见证 |
| [DEC-0104](DEC-0104-address-and-ml-dsa-profile.md) | Address and ML-DSA Profile | Accepted | 密码学 |
| [DEC-0201](DEC-0201-utxo-utco-state-fingerprint.md) | UTXO/UTCO State Fingerprint | Accepted | 状态 |
| [DEC-0301](DEC-0301-poh-mint-hash-and-mint-proof.md) | PoH Mint Hash and Mint Proof | Accepted | PoH |
| [DEC-0302](DEC-0302-genesis-and-initial-window.md) | Genesis and Initial Window | Accepted | 创世初段 |
| [DEC-0303](DEC-0303-fork-choice-and-randomx-tiebreaker.md) | Fork Choice and RandomX Tiebreaker | Accepted | 分叉 |
| [DEC-0401](DEC-0401-coinbase-serialization-rewards-and-award-slots.md) | Coinbase Serialization, Rewards and Award Slots | Accepted | Coinbase |
| [DEC-0501](DEC-0501-script-bytecode-encoding.md) | Script Bytecode Encoding | Accepted | 脚本 |
| [DEC-0502](DEC-0502-script-float-profile.md) | Script Float Profile | Accepted | 脚本 |
| [DEC-0503](DEC-0503-script-registry-and-environment-boundary.md) | Script Registry and Environment Boundary | Accepted | 脚本 |
| [DEC-0504](DEC-0504-script-cost-budget.md) | Script Cost Budget | Accepted | 脚本 |
| [DEC-0505](DEC-0505-script-failure-and-disabled-opcodes.md) | Script Failure and Disabled Opcodes | Accepted | 脚本 |
| [DEC-0601](DEC-0601-block-proof-package.md) | Block Proof Package | Accepted | 证明 |
| [DEC-0602](DEC-0602-network-summary-txid-profile.md) | Network Summary TxID Profile | Accepted | 网络概要 |
| [DEC-0603](DEC-0603-blockqs-verification-data-profile.md) | Blockqs Verification Data Profile | Accepted | Blockqs |

## Topic Index（主题索引）

| 主题 | 决策 |
|------|------|
| 基础编码与哈希 | [DEC-0001](DEC-0001-canonical-integer-and-bytes-encoding.md), [DEC-0002](DEC-0002-domain-tags-and-hash-profiles.md), [DEC-0003](DEC-0003-block-and-transaction-field-encoding.md), [DEC-0004](DEC-0004-hash-tree-and-proof-edge-cases.md) |
| 交易、签名与密码学 | [DEC-0101](DEC-0101-transaction-body-and-output-payloads.md), [DEC-0102](DEC-0102-signature-message-profile.md), [DEC-0103](DEC-0103-witness-container-and-pruning.md), [DEC-0104](DEC-0104-address-and-ml-dsa-profile.md) |
| 状态 | [DEC-0201](DEC-0201-utxo-utco-state-fingerprint.md) |
| 共识 | [DEC-0301](DEC-0301-poh-mint-hash-and-mint-proof.md), [DEC-0302](DEC-0302-genesis-and-initial-window.md), [DEC-0303](DEC-0303-fork-choice-and-randomx-tiebreaker.md) |
| 激励 | [DEC-0401](DEC-0401-coinbase-serialization-rewards-and-award-slots.md) |
| 脚本 | [DEC-0501](DEC-0501-script-bytecode-encoding.md), [DEC-0502](DEC-0502-script-float-profile.md), [DEC-0503](DEC-0503-script-registry-and-environment-boundary.md), [DEC-0504](DEC-0504-script-cost-budget.md), [DEC-0505](DEC-0505-script-failure-and-disabled-opcodes.md) |
| 证明与服务 | [DEC-0601](DEC-0601-block-proof-package.md), [DEC-0602](DEC-0602-network-summary-txid-profile.md), [DEC-0603](DEC-0603-blockqs-verification-data-profile.md) |

## Status Summary（状态统计）

| Status | Count |
|--------|-------|
| Proposed | 0 |
| Accepted | 21 |
| Deprecated | 0 |
| Absorbed | 0（历史清单见下文） |

## Absorbed Or Removed（已吸收或移除）

| 旧内容 | 处理 |
|--------|------|
| 旧 `DEC-0001` 至 `DEC-0004` | 重写为 `DEC-0001` 至 `DEC-0004`，只保留规范编码和哈希树边界问题。 |
| 旧 `DEC-0005`、`DEC-0006` | 合并为 `DEC-0104`；地址校验流程已由 `docs/conception/附.交易.md` 明确。 |
| 旧 `DEC-0007` 至 `DEC-0009` | 重写为 `DEC-0101` 至 `DEC-0103`。 |
| 旧 `DEC-0010` | 重写为 `DEC-0201`。 |
| 旧 `DEC-0011` 至 `DEC-0015` | 重写为 `DEC-0301` 至 `DEC-0303`；PoH 主要参数已由 conception 明确。 |
| 旧 `DEC-0016`、`DEC-0018` 至 `DEC-0020` | 合并为 `DEC-0401`；发行与激活规则已由 conception 明确，Decision 只补编码和边界。 |
| 旧 `DEC-0017` | 删除。Coinbase 省略 `HashInputs` 已由 conception 明确，不再需要废弃占位。 |
| 旧 `DEC-0021` | 删除。全网通告设计已取消，不再保留废弃占位。 |
| 旧 `DEC-0022` 至 `DEC-0026` | 重写为 `DEC-0501` 至 `DEC-0505`。 |
| 旧 `DEC-0027` 至 `DEC-0029` | 重写为 `DEC-0601` 至 `DEC-0603`。 |

## Open Question Groups（开放问题分组）

- 区块、交易和 Coinbase 编码冻结：见 `DEC-0003`、`DEC-0101`、`DEC-0401`。
- PoH 铸凭前像和分叉竞争边界：见 `DEC-0301`、`DEC-0303`。
- 脚本执行、成本和禁用指令：见 `DEC-0501` 至 `DEC-0505`。
- 服务证明与查询格式：见 `DEC-0601` 至 `DEC-0603`。

## Maintenance Rules（维护规则）

- 新增 Decision 前必须先检查 `docs/conception/` 是否已经明确该规则。
- 若 conception 已明确，只在 README 的吸收清单中记录，不新增 DEC。
- 若无法从 conception 直接裁定，状态必须为 `Proposed`，并列出待裁决参数。
- 若后续 conception 吸收某 DEC，应删除该 DEC，或迁移为 README 中的吸收记录。
