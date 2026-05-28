# DEC-0401: Coinbase Serialization, Rewards and Award Slots（Coinbase 序列化、奖励与兑奖槽）

Status: Accepted

## Context（背景）

Conception 已明确 Coinbase 无输入、收益来源、奖励比例、公共服务 48 块兑奖窗口、百日前不启用公共服务奖励，以及长期发行规则。本决策冻结输出顺序、取整、兑奖槽位、截留回收和百日前 profile。

## Decision（决策）

Coinbase 输出配置值和百日后输出顺序：

- `1`：铸凭者，10%。
- `2`：校验组，40%。
- `3`：Blockqs，20%。
- `4`：Depots，20%。
- `5`：STUN，10%。

百日后 Coinbase 必须包含 5 笔输出，并按配置值升序排列：铸凭者、校验组、Blockqs、Depots、STUN。

百日前 Coinbase profile：

- `height <= 24000` 时无公共服务奖励。
- `height <= 24000` 时 Coinbase 只包含 2 笔输出：铸凭者、校验组。
- 百日前比例重标定为铸凭者 20%、校验组 80%。
- 百日前不包含 Blockqs、Depots、STUN 输出，也不使用 `SYS_AWARD`。
- `height >= 24001` 起启用五输出 profile，比例为 10/40/20/20/10。

金额规则：

```text
burned_tx_fee = total_tx_fee / 2
unburned_tx_fee = total_tx_fee - burned_tx_fee
RewardBase = issuance + unburned_tx_fee + reclaimed_award
```

- 整数除法向下取整。
- 若交易费为奇数 `chx`，多出来的 1 `chx` 归未销毁部分。
- Coinbase 头字段 `BurnCoin` 记录非负的 `burned_tx_fee`。
- 不再使用负值 `BurnCoin` 表达奇偶或余数。
- 奖励金额按输出顺序逐项计算。
- 前 N-1 项按 `RewardBase * percent / 100` 向下取整。
- 最后一项获得剩余全部金额，承接所有余数。
- 百日前最后一项是校验组；百日后最后一项是 STUN。

兑奖槽：

- Blockqs、Depots、STUN 各自 6 字节，共 18 字节。
- 每个服务槽覆盖前 48 个区块，每块 1 bit。
- Coinbase 头部携带 `AwardSlots [18]byte`，不作为输出项。
- 槽位顺序为 Blockqs 6 字节、Depots 6 字节、STUN 6 字节。
- bit0 对应 `H-1`，bit47 对应 `H-48`。
- 某区块 `K` 的公共服务奖励，在 `K+1..K+48` 被后续 Coinbase 的对应服务槽确认。
- 达到 1 次确认可兑 50%，达到 2 次确认可兑 100%。
- 花费公共服务奖励输出时，必须至少在 `K+31` 后，避免分叉影响。
- 到 `K+49` 时，未被确认可兑的剩余部分进入该块 Coinbase 的 `reclaimed_award`。

`reclaimed_award` 表达：

- 不新增单独输出项。
- 不作为 Coinbase 头字段。
- 作为 `RewardBase` 的隐含输入项，由验证器根据 `H-49` 区块公共服务输出金额和后续 48 块兑奖槽计算得出。
- 当前块 Coinbase 金额校验时必须重算 `reclaimed_award`，并纳入 `RewardBase`。

## Rationale（理由）

Coinbase 输出顺序会直接影响 TxID，因此固定为配置值升序，与 `附.交易.md` 中的 5 笔输出顺序一致。百日前无公共服务奖励，按铸凭者与校验组原始 10:40 比例重标定为 20:80，可与创世示例保持一致。把交易费销毁和奖励分配放在同一 DEC，可避免 Coinbase 金额校验分散。兑奖槽分服务独立，符合 conception 对后期调整的预留。

`reclaimed_award` 采用隐含计算，可避免 Coinbase 自报回收额，也不增加额外头字段或输出项。`BurnCoin` 统一为非负销毁额，可通过 `total_tx_fee - burned_tx_fee` 无歧义推出未销毁部分，负值奇偶编码不再需要。

## Consequences（影响）

- Coinbase 序列化必须包含 `AwardSlots [18]byte` 字段；百日前该字段应为全零或按编码规则表达为空槽。
- `SYS_AWARD` 只适用于公共服务奖励输出；百日前 Coinbase 不包含此类输出。
- 第 49 块回收额不直接编码，但会影响 Coinbase 输出金额和 TxID。
- 测试需要覆盖高度 0、24000、24001，奇数交易费，输出余数归属，兑奖 0/50/100%，以及第 49 块回收。

## Conception References（构想层依据）

- `docs/conception/blockchain.md#百日扩张`
- `docs/conception/4.激励机制.md`
- `docs/conception/附.交易.md#铸币交易coinbase`

## Open Questions（开放问题）

- 无。
