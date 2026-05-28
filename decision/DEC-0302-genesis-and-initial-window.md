# DEC-0302: Genesis and Initial Window（创世与初段窗口）

Status: Accepted

## Context（背景）

Conception 已明确创世区块、创世 Coinbase、#1/#2 启动逻辑和百日扩张阶段。本决策冻结创世工件的发布格式、初段评参区块、初段铸凭交易高度放宽规则，以及 Coinbase 作为铸凭交易的边界。

## Decision（决策）

创世工件必须包含：

- 创世区块头完整编码。
- 创世 Coinbase 完整交易体。
- 创世铸造者对 Coinbase 的签名。
- 创世铸造者对 `CheckRoot` 的签名。
- 创世声明 `FreeData`。

创世工件同时发布两套形式：

- `genesis.bin`：canonical 二进制编码，作为链上共识和客户端硬编码的权威形式。
- `genesis.json`：可读描述形式，仅用于人工审阅、工具生成和交叉校验，不参与共识。

创世 Coinbase 的 `Minter` 省略。创世块的 `MintHash` 定义为 32 字节全零值，用于初段评参区块铸凭哈希引用。

初段规则：

- `currentHeight < 8` 时评参区块高度为 `0`。
- `currentHeight >= 8` 时评参区块高度为 `currentHeight - 8`。
- `currentHeight < 480` 时铸凭交易高度检查放宽，但仍必须引用已确认交易。
- `currentHeight >= 480` 时使用 `h > 239 && h <= 80000`。

其中铸凭交易高度检查应按如下逻辑处理：

```go
if currentHeight < 480 {
    return txHeight < currentHeight
}
h := currentHeight - txHeight
return h > 239 && h <= 80000
```

#1/#2 规则：

- #1 由创世 `MintPKHash` 铸造。
- #1 可以包含花费创世输出的普通交易。
- #2 起允许基于已确认初段交易形成竞争。

Coinbase 可作为铸凭交易：

- 除初期特殊规则外，任何高度只要满足铸凭交易高度范围约束，且交易头设置了 `MintPKHash`，没有首领输入的 Coinbase 也可以成为正常铸凭交易。
- 初段从 #2 起，已确认 Coinbase 与普通交易都可以作为铸凭交易。
- 当前待铸区块内的交易不得作为当前区块的铸凭交易。

## Rationale（理由）

把创世工件作为客户端硬编码资料的一部分，可避免不同实现自行拼接创世块。二进制与 JSON 双形式兼顾共识确定性和人工审阅便利性。创世 `MintHash` 使用全零值，可在不伪造创世 `MintProof` 的前提下，为 #1 到 #7 的评参区块铸凭哈希引用提供确定值。

初段高度检查放宽到 `currentHeight < 480`，可避免进入 `currentHeight >= 240` 后只有第 240 号区块交易满足正常范围要求，给初期竞争留出额外一天缓冲。Coinbase 若显式设置 `MintPKHash`，其铸造者身份可直接验证，因此在满足范围约束时不应被排除在正常铸凭交易之外。

## Consequences（影响）

- 创世工件一旦发布即不可变。
- 创世 `MintHash` 全零值只用于创世块缺省 `Minter` 的引用语义，不表示有效择优凭证。
- 初段测试需要覆盖高度 0、1、2、7、8、239、240、479、480。
- #1 和 #2 的特殊处理不得泄漏到正常高度逻辑。
- Coinbase 铸凭资格验证需要覆盖有 `MintPKHash`、无首领输入、满足/不满足高度范围等场景。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#区块链启动`
- `docs/conception/1.共识-历史证明（PoH）.md#初段规则`

## Open Questions（开放问题）

- 无。
