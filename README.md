## 简介
一个Golang编写的OpenAI API 代理程序，部署在AWS Lambda，用于解决GFW的问题

## 使用方法

### 编译
``` bash
bash build.sh
```
### 部署

编译之后，会把静态资源打包成zip文件

1. 登录AWS Lambda，创建一个新的函数，选择`Go 1.x`作为运行时，开启URL访问
2. 上传zip文件，设置Handler为`openai-proxy-aws`
3. 设置环境变量`OPENAI_API_KEY`为你的OpenAI API Key
4. 保存并测试