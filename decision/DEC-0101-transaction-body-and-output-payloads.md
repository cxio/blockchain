# DEC-0101: Transaction Body and Output Payloads（交易体与输出载荷）

Status: Accepted

## Context（背景）

Conception 明确交易体由输入项和输出项构成，输出项可承载 Coin、Credit、Proof，并通过输出配置字节的高位标记控制接收者、内容、脚本是否以摘要参与输出项哈希。但三类输出 payload 的字段顺序、空值编码、附件字段、销毁输出和输入集合规范序列化尚未冻结。

## Decision（决策）

### 交易体编码

- `Inputs = varint(count) || Input*`
- `Outputs = varint(count) || Output*`
- 输入和输出按交易创建者给定顺序编码，不自动排序。
- `HashInputs` 与 `HashOutputs` 只由规范编码计算，不包含见证信息（解锁脚本另由 DEC-0103 处理）。

### 输入项编码

1. `Year varint`（年度，UTC 自然年份，见 DEC-0001）
2. `TxIDPart bytes`，长度必须 `>=16`
3. `OutIndex varint`
4. `UnlockScript bytes`

### 输出项公共头

1. `Config byte`：高 4 位为位标记（`bit7`=账户摘要 / `bit6`=内容摘要 / `bit5`=脚本摘要 / `bit4` 未用），低 4 位为类型（`1`=币金 / `2`=凭信 / `3`=存证 / `0` 预留）。**不包含销毁位**——普通交易不可销毁币金，销毁仅由 Coinbase 的 `BurnCoin` 字段表达（详见 DEC-0401）。
2. `Payload bytes`
3. `LockScript bytes`

### 三类 payload 编码

- **Coin**：`Receiver(bytes, length<256) || Amount(varint) || Memo(bytes, length<256)`
- **Credit**：`Receiver(bytes, length<256) || Creator(bytes, length<256) || Title(bytes, length<256) || Description(bytes) || AttachmentID(bytes)`
- **Proof**：`Creator(bytes, length<256) || Title(bytes, length<256) || Content(bytes) || AttachmentID(bytes)`

### 可选字段表示

- `Memo`、`AttachmentID` 等所有可选字节序列字段，统一以 `varint(length)==0` 表示缺省。
- 不引入位图标记；解析器遇到 `length==0` 即视为该可选字段未提供。
- 缺省字段在哈希前像中仍参与编码（即包含一个 `0x00` 长度字节），确保规范编码唯一性。

### 摘要标记与输出项哈希

- `bit7` 账户摘要：接收者使用哈希摘要参与计算输出项哈希；该标记通常无用，但作为配置位保留。
- `bit6` 内容摘要：内容部分使用哈希摘要参与计算输出项哈希，主要用于大负载。
- `bit5` 脚本摘要：锁定脚本使用哈希摘要参与计算输出项哈希，主要用于长脚本。
- `bit4` 未用，必须保持未置位。
- 这里的“内容”指输出项内容部分，即除锁定脚本和接收者之外的全部条目。
- 摘要标记只影响输出项哈希/证明的前像选择；完整交易仍以实际字段为源计算摘要，签名消息不读取该配置，并按 DEC-0102 的授权范围覆盖实际内容。

### Credit 过期边界

- Credit 失效条件为 `age > 31 * 87661`（即超过 31 个完整出块年后失效）。
- 当 `age == 31 * 87661` 时，该 Credit 在该区块仍可被引用花销。

## Rationale（理由）

保留输入输出创建顺序可支持授权语义中的 `SIGOUT_SELF` 和脚本引用。公共头拆分可让输出类型扩展不影响基础解析。

`varint(length)==0` 表缺省与 ULEB128 最短编码规则天然契合，避免引入额外位图字节；规范编码唯一性由"缺省字段也参与一字节 0x00 编码"保证。

摘要标记为大负载或长脚本的单项验证提供更小的证明前像：验证方可只取得相应片段摘要与必要的实际片段，而无需总是携带完整负载或完整脚本。该机制不改变三类输出类型，也不改变状态归属。

Credit `age > 31*87661` 边界使 31 周年当区块仍可使用，对持有人更友好；严格"超过 31 年"的自然语义也便于实现 `if age > limit then reject`。

销毁币金从普通交易输出中移除，统一由 Coinbase 的 `BurnCoin` 字段表达（详见 DEC-0401），避免在两处引入销毁语义。

## Consequences（影响）

- 同一输入在一笔交易中重复引用必须被视为非法，避免双花歧义。
- 可选字段必须使用 `varint(length)==0`，解码器不得接受省略整个字段的形式（否则破坏规范编码唯一性）。
- 摘要标记不得改变输出类型或状态归属；`bit7` 只表示账户摘要，不能解释为其它类型扩展或长度计数。
- 当摘要标记置位时，节点验证完整交易必须能从实际字段计算对应摘要，并用摘要参与输出项哈希。
- Credit 的 `age` 计算以"该 Credit 创建区块到当前区块的高度差"为准，由 DEC-0201 配套定义 UTCO 中 Credit 的存储格式。
- 普通交易输出不得设置销毁语义；试图通过 `Receiver` 为空表达销毁的交易必须被拒绝。

## Conception References（构想层依据）

- `docs/conception/附.交易.md#交易体`
- `docs/conception/附.交易.md#输入项`
- `docs/conception/附.交易.md#输出项`
- `docs/conception/5.信用结构.md`

## Open Questions（开放问题）

（无）
