### openssl 非对称加密
1. openssl genrsa -out rsa_private_key.pem 2048
2. openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem

### openssl 对称加密
1. 加密 openssl enc -aes-256-ecb -a -in aaa -out bb -k 1234
2. 解密 openssl enc -aes-256-ecb -d -a -in bbb -out ccc -k 1234

### openssl 证书
1. openssl genrsa -out ca.key 2048
2. openssl req -new -nodes -x509 -key ca.key -days 365 -out ca.crt -subj "/C=CN/ST=Hubei/L=Wuhan/O=k8s/OU=systemGroup/CN=kubernetesEA:admin@com" 



### cfssl 生成证书


1. 在https://github.com/cloudflare/cfssl/releases上进行相关版的下载

    - 下载cfssl-certinfo
    - 下载cfssljson
    - 下载cfssl
    
2. 创建ca-csf.json
    ```json
    {
        "CN": "key-center",
        "key": {
            "algo": "rsa",
            "size": 2048
        },
        "names": [
            {
                "C": "CN",
                "ST": "ShangHai",
                "L": "ShangHai",
                "O": "Transsion",
                "OU": "key-center"
            }
        ],
        "ca": {
            "expiry": "876000h"
        }
    }
```
    