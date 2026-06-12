# DEC-0601: Block Proof Package（区块证明包）

Status: Accepted

## Context（背景）

Conception 允许新区块先广播最小证明，再同步区块概要和交易数据。初始节点验证也需要末端区块的 Coinbase、验证路径、UTXO/UTCO 指纹和铸造者签名。本决策冻结区块证明包字段和快速预验证流程。

## Decision（决策）

区块证明包包含：

1. `BlockHeader`
2. `CoinbaseTx`
3. `CoinbaseTxIndex`，必须为 0
4. `CoinbaseMerklePath`
5. `TreeRoot`
6. `UTXORoot`
7. `UTCORoot`
8. `MinterCheckRootSignature`

字段规则：

- `TreeRoot` 为完整交易树根。验证者也必须用 `CoinbaseTxID`、`CoinbaseTxIndex` 和 `CoinbaseMerklePath` 重算并对比。
- 不单独携带 `CheckRoot`；以 `BlockHeader.CheckRoot` 为准。
- 不单独携带 `MintProof`；从 `CoinbaseTx.Minter` 读取。
- 不单独携带铸造者公钥；从 `CoinbaseTx.Minter.MintPubKey` 读取。
- 不携带 `ChainScope`；链作用域隐含在交易签名和本地链上下文中。
- 不携带 UTXO/UTCO 状态证明路径；`UTXORoot` 和 `UTCORoot` 只用于与本地当前状态快速比较并重算 `CheckRoot`。

快速预验证流程：

- 验证 `BlockHeader.PrevBlock` 与本地末端区块 ID 衔接。
- 验证 `CoinbaseTx.Minter` 存在，且铸造者是本地当前择优池成员。
- 对比 `UTXORoot`、`UTCORoot` 是否与本地当前 UTXO/UTCO 指纹一致。二者是上一区块完成后的结果，新区块广播时应当已知且全网统一。
- 验证 `CoinbaseTxIndex == 0`。
- 计算 `CoinbaseTxID`。
- 用 `CoinbaseTxID`、`CoinbaseTxIndex` 和 `CoinbaseMerklePath` 重算交易树根，并对比 `TreeRoot`。
- 用 `TreeRoot || UTXORoot || UTCORoot` 重算 `CheckRoot`，并对比 `BlockHeader.CheckRoot`。
- 验证 `CoinbaseTx.Minter` 的铸凭证明：`Nonce >= BlockHeader.Height`、`Solution` 可解析且索引升序并无重复、Equi-X 校验通过，并且重算 `MintHash` 与签名匹配（此步应当在择优池成员进入时已验证，可略过重验）。
- 用 `CoinbaseTx.Minter.MintPubKey` 验证 `MinterCheckRootSignature`。

## Rationale（理由）

证明包应足够小，使节点可先转播区块证明；同时必须包含独立验证铸造资格、Coinbase 入树关系和区块头合法性的最小材料。新区块广播发生在校验组管理层之间，`UTXORoot` 和 `UTCORoot` 是上一区块完成后的本地已知状态，证明包只需携带它们用于快速一致性比较和 `CheckRoot` 重算，不需要包含状态证明。

铸造者是否属于本地择优池是廉价且关键的检查，应先于哈希树路径和签名等较高成本验证。`MintProof` 已属于 Coinbase 交易体，重复携带会制造一致性风险；`ChainScope` 已由交易签名和本地链上下文约束，不作为区块证明包字段。

## Consequences（影响）

- 证明包不能替代完整区块验证，只能支持快速预验证和转播。
- 证明包不证明 UTXO/UTCO 状态本身，只要求与本地当前状态一致。
- 初始同步至少需要最近 31 块证明包以覆盖分叉安全窗口。
- 若节点缺少本地 UTXO/UTCO 指纹或需要验证状态真实性，应通过完整区块、Blockqs 或校验组获取额外状态数据；这不属于本证明包。
- 测试需要覆盖择优池成员检查、前一区块衔接、UTXO/UTCO 本地不一致、Coinbase 非 0 序位、TreeRoot 不匹配、MintProof `Nonce` 非法、`Solution` 非升序或有重复、MintProof 签名错误和 CheckRoot 签名错误。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#初始主链验证`
- `docs/conception/2.共识-端点约定.md#区块发布`
- `docs/conception/附.组队校验.md#区块发布`

## Open Questions（开放问题）

- 无。
