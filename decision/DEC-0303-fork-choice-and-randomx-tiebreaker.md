# DEC-0303: Fork Choice and RandomX Tiebreaker（分叉选择与 RandomX 平局裁决）

Status: Accepted

## Context（背景）

Conception 已明确 31 块分叉竞争、16 块过半胜出、长度 20 临界裁决、同铸造者低收益原则和 RandomX 平局方案。本决策冻结完整比较算法、低收益定义、RandomX profile 和交易量约束流程（Stakes 严格超3倍则替换）。

## Decision（决策）

分叉链段比较：

- 只比较分叉点之后最多 31 个区块。
- 每个高度先按本决策的同铸造者多签归一化和交易量约束选出有效候选块。
- 逐高度比较两条链对应有效候选块的 `MintHash`。
- 单高度 `MintHash` 较小者得 1 分。
- 单高度 `MintHash` 完全相等时双方都不得分。
- 任一链先达到 16 分即胜出。
- 31 个高度比较完成仍平局时进入 RandomX 裁决。

同高度同铸造者多签归一化：

- 先按铸造者公钥哈希分组。
- 每组仅保留“铸造者个人可得收益最低”的区块。
- 若个人可得收益相同，保留交易费总额更低者。
- 仍相同则保留 BlockID 更小者。

其中“铸造者个人可得收益”指 Coinbase 中直接分配给该铸造者身份的金额，不包含校验组工作报酬、公共服务奖励或其它第三方收益。

RandomX 平局：

```text
seed = ForkPointBlockID
input = FirstForkBlockID
score = RandomX(seed, input)
```

- `score` 按字典序升序，较小者胜。
- `score` 相同则比较分叉首块 ID，较小者胜。

RandomX profile：

- 使用官方 RandomX 实现：`https://github.com/tevador/RandomX`。
- 版本固定为 `v2.0.1`，tag commit `aaafe71322df6602c21a5c72937ac284724ae561`。
- 输出长度固定为 32 字节，即官方 `RANDOMX_HASH_SIZE`。
- 输入 `seed` 为 48 字节 `ForkPointBlockID`。
- 输入 `input` 为 48 字节 `FirstForkBlockID`。
- 使用完整 VM 语义计算；实现可以通过 CGO 封装官方 C/C++ 库。
- 不得使用会改变哈希结果的非官方兼容实现或参数变体。

交易量约束确定性算法（Stakes `>3x`）：

- 只比较同一高度、同一前一区块上的冗余出块。
- 候选块按铸造者在择优池中的排名升序排列。
- 缺位候选者跳过，不生成空候选。
- 从当前最优候选 `winner` 开始，依次考察后位候选 `challenger`。
- `challenger.Stakes > winner.Stakes * 3` 时，替换 `winner` 并继续考察后位候选；否则停止。
- `TxCount` 仅作上层统计与展示，不参与候选归一化比较。
- 若 `winner.Stakes == 0`，仍按上述公式处理；因此后位候选只要 `Stakes > 0` 即可满足 `>3x`。
- 相等不算超越，必须严格 `> 3x`。

## Rationale（理由）

先对同铸造者多签做归一化，低收益定义采用铸造者个人可得收益，直接对应抑制多签收益动机。

RandomX 仅作为低概率平局裁决，不进入常规出块路径；但一旦触发，哈希正确性直接影响主链选择，因此冻结到官方实现和固定版本。

交易量不参与归一化比较，否则可能导致铸造者竟相构造大量微交易来追求数量优势，反而创造了攻击面。币权更能反应真实交易价值，也更难被操纵。

## Consequences（影响）

- 实现需要引入官方 RandomX 库，Go 实现可通过 CGO 简单封装。
- 低收益比较依赖 Coinbase 中铸造者个人收益金额的可验证计算。
- 交易量约束（Stakes `>3x` 或 TxCount `>2x`）的冗余出块规则需在区块接收阶段独立实现，并在分叉链段比较前完成候选块选择。
- 测试需要覆盖同铸造者多签、`MintHash` 相等、31 块平局、RandomX 裁决、连续后位超越、`Stakes=0`、`TxCount=0`、仅 Stakes 超越、仅 TxCount 超越和相等边界。

## Conception References（构想层依据）

- `docs/conception/2.共识-端点约定.md#分叉竞争`
- `docs/conception/2.共识-端点约定.md#平局的可能性及解决`
- `docs/conception/附.组队校验.md#低收益原则`
- `docs/conception/附.组队校验.md#交易量约束`

## Open Questions（开放问题）

- 无。
