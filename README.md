# What is iHosts
A tool for updating the hosts file

## How To Use ?
- `go mod tidy`
- `go run main.go bilibili/config/popular.json`

## FAQ
- 配置文件该怎么写？
  - 请参考 `bilibili/config/popular.json`
  - 注意的点: "html_content_xpath" 和 "root" 要用 Js ，而非 xpath,其必须具有 length 方法。
  - "max_request" 要根据机器性能做设定，默认设置的是 10，如果机器性能不够强，建议设置成 2-5。
  - 如果网速慢建议增加 "timeout" 的值。