# 贡献指南

欢迎为项目贡献代码！在提交贡献之前，请先阅读以下指南。

## 代码风格
请遵循 Go 语言的官方代码风格规范，使用 `gofmt` 格式化代码。

## 提交规范
- 提交信息应简洁明了，描述所做的更改。
- 使用英文祈使句，如 "Add new feature" 或 "Fix bug in user authentication"。

## 测试
在提交代码之前，请确保所有测试都能通过。运行测试的命令如下：
```bash
go test ./...
```

## 提交贡献的步骤
1. Fork 本项目。
2. 创建新分支 ( git checkout -b feature/your-feature )。
3. 提交修改 ( git commit -am 'Add some feature' )。
4. 推送到分支 ( git push origin feature/your-feature )。
5. 创建 Pull Request。
