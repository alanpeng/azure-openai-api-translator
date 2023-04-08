# 自签泛域名证书
此工具用于颁发泛域名证书，方便开发环境调试。

## 使用
```bash
./gen.cert.sh <domain> [<domain2>] [<domain3>] [<domain4>] ...
```
把 `<domain>` 替换成你的域名，例如 `kubechatgpt.com`

生成的证书位于：
```text
out/<domain>/<domain>.crt
out/<domain>/<domain>.bundle.crt
```

证书有效期是 100 年，你可以修改 `ca.cnf` 来修改这个年限。

根证书位于：  
`out/root.crt`  
成功之后，把根证书导入到操作系统里面，信任这个证书。

根证书的有效期是 100 年，你可以修改 `gen.root.sh` 来修改这个年限。

证书私钥位于：  
`out/cert.key.pem`

其中 `<domain>.bundle.crt` 是已经拼接好 CA 的证书，可以添加到 `nginx` 配置里面。  
然后你就可以愉快地用 `https` 来访问你本地的开发网站了。

## 清空
你可以运行 `flush.sh` 来清空所有历史，包括根证书和网站证书。

## 配置
你可以修改 `ca.cnf` 来修改你的证书年限。
```ini
default_days    = 730
```

可以修改 `gen.root.sh` 来自定义你的根证书名称和组织。

也可以修改 `gen.cert.sh` 来自定义你的网站证书组织。
