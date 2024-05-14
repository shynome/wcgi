# 简介

提升 [go-wagi](https://github.com/shynome/go-wagi) 性能的库

具体做法就是在 `stdio` 上建立 `yamux` 双向连接复用进程实现的

# 如何使用?

调用的时候含有环境变量 `WAGI_WCGI=true` 时就会使用 `yamux` 复用 `stdio`
