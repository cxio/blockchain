# DEC-0602: Network Summary TxID Profile（网络概要交易 ID 配置）

Status: Accepted

## Context（背景）

Conception 建议区块概要中每个交易 ID 仅截取前 16 字节，以降低同步数据量；碰撞时发布者补充更长信息。短摘要格式、碰撞回退和与完整交易树的关系需要冻结。

## Decision（决策）

区块概要基础格式：

```text
Summary = BlockID || TxCount || TxIDPrefix*
```

规则：

- `TxIDPrefix` 固定为完整 TxID 的前 16 字节，不设置 `TxIDPrefixLen` 字段，也不协商其它长度。
- `TxIDPrefix` 按区块交易序列顺序排列，包含 Coinbase。
- 接收方发现本地候选交易中有多个匹配时，按交易序位请求碰撞回退信息。
- 碰撞回退响应不属于基础 `Summary` 本体；发布方对指定交易序位返回完整 48 字节 TxID。
- 区块概要不需要发布方单独签名，只作为网络同步优化信息。
- 最终验证必须使用完整 TxID 序列计算交易树根。

## Rationale（理由）

16 字节前缀足以覆盖常规同步场景，并显著降低区块概要大小。固定长度可以避免协商复杂度；碰撞极少，直接返回完整 TxID 比逐步延长前缀更简单。区块概要不参与共识，真实性最终由区块证明包和完整交易树验证保证，因此无需额外签名。

## Consequences（影响）

- 区块概要只是网络优化，不是共识数据。
- 节点不得因短前缀无法解析就接受不完整区块。
- 恶意发布方提供错误摘要会在完整交易树验证阶段失败。
- 错误或缺失的碰撞回退响应只能导致同步失败或重试，不能使节点接受无效区块。

## Conception References（构想层依据）

- `docs/conception/附.组队校验.md#同步优化`
- `docs/conception/附.交易.md#输入项`

## Open Questions（开放问题）

None.
