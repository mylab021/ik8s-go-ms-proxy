iK8S_PROXY_VERSION=v1.0.2

docker build -f docker/Dockerfile.user-service -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/user-service:${iK8S_PROXY_VERSION} \
                                               -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/user-service:latest \
                                               .
docker build -f docker/Dockerfile.order-service -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/order-service:${iK8S_PROXY_VERSION} \
                                                -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/order-service:latest \
                                                .
docker build -f docker/Dockerfile.product-service -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/product-service:${iK8S_PROXY_VERSION} \
                                                -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/product-service:latest \
                                                .

docker build -f docker/Dockerfile.gateway -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/gateway:${iK8S_PROXY_VERSION} \
                                                -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/gateway:latest \
                                                .

docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/user-service:${iK8S_PROXY_VERSION}
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/order-service:${iK8S_PROXY_VERSION}
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/product-service:${iK8S_PROXY_VERSION}
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/gateway:${iK8S_PROXY_VERSION}

docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/user-service:latest
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/order-service:latest
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/product-service:latest
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/gateway:latest