### openssl 生成证书
1. openssl genrsa -out rsa_private_key.pem 1024
2. openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem


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
    