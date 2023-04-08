package main

import (
	"github.com/alanpeng/azure-openai-api-translator/pkg/azure"
	"github.com/alanpeng/azure-openai-api-translator/pkg/openai"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"flag"
	"crypto/tls"
)

var (
	Address   = "0.0.0.0:8080"
	ProxyMode = "azure"
	CertFile = ""
    KeyFile  = ""
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	flag.StringVar(&CertFile, "certFile", "", "Path to the SSL certificate file")
    flag.StringVar(&KeyFile, "keyFile", "", "Path to the SSL private key file")
}

func main() {
	// 定义命令行参数
	addressFlag := flag.String("address", "0.0.0.0:8080", "The address for the proxy server to listen on")
	proxyModeFlag := flag.String("mode", "azure", "The proxy mode to use (azure or openai)")
	apiVersionFlag := flag.String("azureApiVersion", "2023-03-15-preview", "Azure OpenAI API version")
	endpointFlag := flag.String("azureEndpoint", "https://xxxxxxxx.openai.azure.com/", "Azure OpenAI Endpoint")
	modelMapperFlag := flag.String("azureModelMapper", "gpt-3.5-turbo=gpt-35-turbo,gpt-3.5-turbo-0301=gpt-35-turbo-0301", "Azure OpenAI Model Mapper")
	tokenFlag := flag.String("azureToken", "", "Azure OpenAI Token")
	

	// 解析命令行参数
	flag.Parse()

	// 设置从命令行参数获取的代理服务器地址和模式
	Address = *addressFlag
	ProxyMode = *proxyModeFlag

	log.Printf("loading azure openai proxy address: %s", Address)
	log.Printf("loading azure openai proxy mode: %s", ProxyMode)

	// 初始化 azure 包
	azure.InitAzureConfig(*apiVersionFlag, *endpointFlag, *modelMapperFlag, *tokenFlag)

	router := gin.Default()
	if ProxyMode == "azure" {
		router.GET("/v1/models", handleGetModels)
		router.OPTIONS("/v1/*path", handleOptions)

		router.POST("/v1/chat/completions", handleAzureProxy)
		router.POST("/v1/completions", handleAzureProxy)
		router.POST("/v1/embeddings", handleAzureProxy)
	} else {
		router.Any("*path", handleOpenAIProxy)
	}
    // 使用HTTPS
    if CertFile != "" && KeyFile != "" {
        log.Printf("Starting server with HTTPS on %s", Address)
        server := &http.Server{
            Addr:      Address,
            Handler:   router,
            TLSConfig: &tls.Config{},
        }
        log.Fatal(server.ListenAndServeTLS(CertFile, KeyFile))
    } else {
        log.Printf("Starting server with HTTP on %s", Address)
        router.Run(Address)
    }


	//router.Run(Address)
}

func handleGetModels(c *gin.Context) {
	models := []string{"gpt-4", "gpt-4-0314", "gpt-4-32k", "gpt-4-32k-0314", "gpt-3.5-turbo", "gpt-3.5-turbo-0301", "text-davinci-003", "text-embedding-ada-002"}
	result := azure.ListModelResponse{
		Object: "list",
	}
	for _, model := range models {
		result.Data = append(result.Data, azure.Model{
			ID:      model,
			Object:  "model",
			Created: 1677649963,
			OwnedBy: "openai",
			Permission: []azure.ModelPermission{
				{
					ID:                 "",
					Object:             "model",
					Created:            1679602087,
					AllowCreateEngine:  true,
					AllowSampling:      true,
					AllowLogprobs:      true,
					AllowSearchIndices: true,
					AllowView:          true,
					AllowFineTuning:    true,
					Organization:       "*",
					Group:              nil,
					IsBlocking:         false,
				},
			},
			Root:   model,
			Parent: nil,
		})
	}
	c.JSON(200, result)
}

func handleOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Status(200)
	return
}

func handleAzureProxy(c *gin.Context) {
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Status(200)
		return
	}

	server := azure.NewOpenAIReverseProxy()
	server.ServeHTTP(c.Writer, c.Request)
	//BUGFIX: try to fix the difference between azure and openai
	//Azure's response is missing a \n at the end of the stream
	//see https://github.com/Chanzhaoyu/chatgpt-web/issues/831
	if c.Writer.Header().Get("Content-Type") == "text/event-stream" {
		if _, err := c.Writer.Write([]byte("\n")); err != nil {
			log.Printf("rewrite azure response error: %v", err)
		}
	}
}

func handleOpenAIProxy(c *gin.Context) {
	server := openai.NewOpenAIReverseProxy()
	server.ServeHTTP(c.Writer, c.Request)
}
