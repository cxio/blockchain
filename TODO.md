## 脚本指令修改项

- `SOURCE` 指令扩展&修改了标识，检查修订。
- `DICT` 可否泛型化，以支持其它可比较类型（如 `int`）？
- `LUA{}` 指令的处理，添加外部扩展性。

- `RETURN` 指令，以及 `MAP` 指令行为有调整。
- 较多指令附参已变更为变长整数，主要为长度和数量类。

- `GOTO` 和 `JUMP` 指令的附参修改为与交易输入格式相似：`n+32+n`。

请检查原实现并修改。
