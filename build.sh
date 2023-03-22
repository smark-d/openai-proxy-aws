
rm -rf openai-proxy-aws openai-proxy-aws.zip
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
zip openai-proxy-aws.zip openai-proxy-aws