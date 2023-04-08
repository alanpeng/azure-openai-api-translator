# Azure-OpenAI-API-Translator
Web service to translate OpenAI request into Azure OpenAI request

## Usage: (命令用法，使用--help参数获得详细参数说明)

./azure-openai-api-translator --help

命令返回信息：

Usage of ./azure-openai-api-translator:

  -address string
  
    	The address for the proxy server to listen on (default "0.0.0.0:8080") 本服务侦听地址与端口，默认0.0.0.0:8080
      
  -azureApiVersion string
  
    	Azure OpenAI API version (default "2023-03-15-preview") 默认的OpenAI API版本
      
  -azureEndpoint string
  
    	Azure OpenAI Endpoint (default "https://xxxxxxxx.openai.azure.com/") 您的Azure OpenAI Serivce服务的访问入口
      
  -azureModelMapper string
  
    	Azure OpenAI Model Mapper (default "gpt-3.5-turbo=gpt-35-turbo,gpt-3.5-turbo-0301=gpt-35-turbo-0301") 模型名称映射，等于号右侧的是Azure平台上的模型部署名称
      
  -azureToken string
  
    	Azure OpenAI Token 访问密钥 （可选，也可以由客户端提供）
      
  -certFile string
  
    	Path to the SSL certificate file （如果需要使用https证书，可选用此服务证书文件，可以使用本项目的self-signed-certificate目录下的命令创建自签名证书）
      
  -keyFile string
  
    	Path to the SSL private key file （如果需要使用https证书，可选用此服务私钥文件，可以使用本项目的self-signed-certificate目录下的命令创建自签名证书）
      
  -mode string
  
    	The proxy mode to use (azure or openai) (default "azure")  翻译或代理的模式，默认是将OpenAI的API请求翻译为Azure OpenAI Service的API请求。

## Command sample to start https service:
```
./azure-openai-api-translator -certFile ./kubechatgpt.com.crt -keyFile ./kubechatgpt.com.key.pem -address "0.0.0.0:8080" -azureEndpoint "https://myazureservicename.openai.azure.com/" -azureModelMapper "gpt-3.5-turbo=gpt35"
```

## Command sample to start http service:
```
./azure-openai-api-translator -address "0.0.0.0:8080" -azureEndpoint "https://myazureservicename.openai.azure.com/" -azureModelMapper "gpt-3.5-turbo=gpt35"
```

以上命令假设：

Microsoft Azure OpenAI Service Endpoint: "https://myazureservicename.openai.azure.com/" 比如这是您的Azure OpenAI Service的服务实例名

Microsoft Azure OpenAI Service Model deployment name: gpt35  比如您的Azure服务部署的模型名称是gpt35

本程序支持Windows、MacOS、Linux平台。
