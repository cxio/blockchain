# DEC-0102: Signature Message Profile（签名消息配置）

Status: Accepted

## Context（背景）

Conception 定义了普通交易的授权标记和链标识前缀，但签名消息的精确字段顺序、缺省字段、输入/输出子集编码和 Coinbase 签名消息尚未冻结。

## Decision（决策）

### 普通交易签名消息

签名消息整体布局：

```text
SigMessage = DomainTag("signature.message")
           || ChainScope
           || SigScope
           || TxHeaderCore
           || CoveredInputs
           || CoveredOutputs
```

- **DomainTag**：见 DEC-0002，前缀 `"Evidcoin/v1/signature.message\x00"`。
- **ChainScope**：
    ```text
    ChainScope = ProtoLen || ProtocolID
              || ChainLen || ChainID
              || GenesisID(48 bytes)
              || BoundLen  || BoundID
    ```
    - `ProtocolID`、`ChainID` 按 `varint(length) || ASCII bytes` 编码（取自 conception `blockchain.md#主链和分叉标识`）。
    - `GenesisID` 为创世区块哈希，固定 48 字节（SHA3-384，见 DEC-0002）。
    - `BoundID` 可空；空时编码为 `varint(0)`（占位），不省略字段，保证主链/分叉链签名消息结构一致。
- **SigScope**：`chk_type(1B) || auth_flag(1B) || input_index(varint)`。
    - `chk_type ∈ {1=币金花费, 2=凭信转移}`，与脚本 `FN_CHECKSIG`/`FN_MCHECKSIG` 第一实参一致。
    - `auth_flag` 即 conception `附.交易.md#授权种类` 的 8 位标记字节。
    - `input_index` 为当前正在被验证的输入项序位（验证器注入）。
- **TxHeaderCore**：固定包含 `Version(uint16, BE)` 与 `Timestamp(int64, BE)`；当交易头存在铸凭公钥哈希时追加 `MintPKHash(varint(length) || bytes)`，否则该字段编码为 `varint(0)`。
- **CoveredInputs**、**CoveredOutputs**：依 `auth_flag` 选定的子集，规则见下文"覆盖范围"。

### 覆盖范围规则

- **未被 `auth_flag` 覆盖的字段完全不进入签名消息**（与 conception `附.交易.md` "未指定的部分不会包含" 一致）。
- **CoveredInputs**：
    - `SIGIN_ALL`：依 DEC-0101 输入项规范编码全部输入项，**含解锁脚本**，序列前置 `varint(count)`。
    - `SIGIN_SELF`：仅当前 `input_index` 对应的输入项，前置 `varint(input_index)` 后跟该输入项规范编码。
    - 二者皆未设置：`CoveredInputs` 为空（零字节）。
- **CoveredOutputs**：
    - 主项 `SIGOUT_ALL` / `SIGOUT_SELF` 决定**输出范围**：
        - `SIGOUT_ALL`：全部输出，前置 `varint(count)`；每输出前置自身序位（`varint(index)`）。
        - `SIGOUT_SELF`：仅与 `input_index` 同序位的输出；若 `input_index >= len(outputs)`，签名验证必须失败。
    - 辅项决定**每个输出内嵌入哪些字段**：
        - `SCRIPT`：仅 `LockScript`
        - `CONTENT`：除 `LockScript` 与接收者外的所有字段（即 Payload 中去除 `Receiver`/`Creator` 等接收者标识后剩余内容）
        - `RECEIVER`：仅接收者字段（Coin/Credit 的 `Receiver`，Proof 的 `Creator` 视作接收者）
        - `OUTPUT`：等价于 `SCRIPT|CONTENT|RECEIVER` 的语义并集
    - 每个嵌入字段段以 `varint(length) || bytes` 编码；未被辅项选中的字段不出现。

### 辅项冲突规则

- **`OUTPUT` 与 `SCRIPT|CONTENT|RECEIVER` 任一位互斥**：若 `auth_flag` 中 `OUTPUT` 位与其它三位中任意一位同时为 1，签名消息构造失败，交易必须被拒绝。
- 主项缺失但辅项存在 / 辅项缺失但主项存在 均按 conception `附.交易.md` "辅项需与主项配合" 处理，构造失败即拒绝。
- 该规范化保证签名消息的字节序列唯一对应一组 `auth_flag`，避免歧义解释。

### Coinbase 签名消息

- Coinbase 不使用授权种类标记（依 conception `附.交易.md:319`、`blockchain.md:200`）。
- Coinbase 签名消息为：
    ```text
    CoinbaseSigMessage = DomainTag("signature.message")
                      || ChainScope
                      || 0x00          // chk_type = 0，标记 Coinbase 域
                      || CoinbaseTxID  // 完整 Coinbase 交易的 TxID（48 字节）
    ```
- 铸造者对区块 `CheckRoot` 的签名独立于 Coinbase 交易签名，使用各自的域标签 / 消息构造。

### 多签签名集顺序

- 解锁脚本 `FN_MCHECKSIG` 实参中签名集（`sigs`）、公钥集（`pubKeys`）、补全公钥哈希集（`baseHashs`）按**见证提供顺序**入栈与编码，不强制按公钥哈希排序。
- conception `附.交易.md:280` "结果集与补全集里的初级哈希混合排序、串联" 仅用于**复合公钥哈希计算**，不约束签名集本身的顺序。
- 同一多签结果允许有多种合法字节排列；规范编码唯一性由交易内 `UnlockScript` 整体哈希保证，不强加签名集排序。

## Rationale（理由）

签名消息需要同时支持普通支付授权、局部授权和链重放隔离。链标识作为签名前缀与 conception `blockchain.md#主链和分叉标识` 一致；外层增加 `DomainTag("signature.message")` 防止与其它哈希用途碰撞（DEC-0002 域隔离原则）。

`BoundID` 空时仍占位编码（`varint(0)`），确保主链转分叉链时签名消息结构稳定，便于实现层用同一编解码器处理两种状态。

辅项互斥（OUTPUT 与 SCRIPT/CONTENT/RECEIVER 并存即拒绝）来自 conception "保持简化" 的隐含哲学，并与 DEC-0001 "最短编码、规范唯一" 原则一致。

多签签名集按见证提供顺序保留构造者意图，与 `FN_MCHECKSIG` 实参顺序天然一致；conception "混合排序" 仅用于复合公钥哈希派生，不约束签名集字节序列。

Coinbase 无输入，使用完整 TxID 作为签名消息可避免授权子集带来的歧义，与 `blockchain.md:200`、`附.交易.md:319` 一致。

## Consequences（影响）

- 签名验证必须知道当前 `input_index`；脚本引擎在执行 `FN_CHECKSIG`/`FN_MCHECKSIG` 时由运行时上下文注入。
- `SIGOUT_SELF` 在 `input_index >= len(outputs)` 时签名验证应失败。
- `auth_flag` 中 `OUTPUT` 与辅项三位之一同时置 1 的交易必须被拒绝；该规则在交易准入阶段即可静态检查。
- 多签签名集顺序由解锁脚本构造者决定；同一 (M-of-N) 配置允许多个合法签名见证组合。
- 剪枝见证后，长期安全依赖 TxID 和区块哈希链，不依赖签名重验（conception `blockchain.md:276`）。
- 节点 P2P 握手使用的链标识声明（conception `blockchain.md:279`）与签名消息中的 `ChainScope` 在语义上一致，但握手可使用字符串形式，签名消息按本 DEC 字节布局。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#主链和分叉标识`
- `docs/conception/附.交易.md#签名消息`
- `docs/conception/6.脚本系统.md#例币金支付验证`

## Open Questions（开放问题）

（无）
