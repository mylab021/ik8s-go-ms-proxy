docker build -f Dockerfile.user-service -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/user-service:latest .
docker build -f Dockerfile.order-service -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/order-service:latest .
docker build -f Dockerfile.gateway -t harbor-srv01.mylab021.com/ik8s-go-ms-proxy/gateway:latest .


docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/user-service:latest
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/order-service:latest
docker push harbor-srv01.mylab021.com/ik8s-go-ms-proxy/gateway:latest