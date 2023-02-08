// inittestenv.sh 用于初始化单机虚拟网口环境
// starttestenv.sh 用于配置服务器环境变量，启动服务
// stoptestenv.sh 关闭所有服务
// putmapping.sh 元数据服务启动以后，在其上创建 metadata 索引以及 objects 类型的映射。

//计算内容sha-256的hash值
// ＄ echo -n "this is obiect test3"|openssl dgst -sha256 -binary|base64
// GYqgadpPt+CSenupao/Ceu3kwcwnoADKaxpiZtabgau-

// 存放test3对象，把散列值加入 PUT 请求的 Digest 头部.
// $ curl -v 10.29.2.2:12345/objects/test3 -XPUT -d"this is object test3”
//-H "Digest: SHA-256-GYqgadpPt+CSenupao/Ceu3kwcwnoADKaxpiZtabgau-"

//存放一个新test3对象
// cur1 -v 10.29.2.1:12345/objects/test3 -xPUT -a"this is objeet test3
// version 2" -H "Digest: SHA-256=CAPVSxZe1PR54ZIESQYOBaxC1PYJIVaHSF3qEOZYYIo="

//查看存放对象的服务器
// cur1 10.29.2.1:12345/locate/GYqqAdFPt+CScnUDc0g2FGcu3kwCWmOADKNYpi ZtdbgsM=
// "10.29.1.6:12345"
// cur1 10.29.2.1:12345/1ocate/cAPVSxZe1PR54ZIESQYOBaxC1pYJIvaHSF3gEOZYYIO=
// "10.29.1.3:12345"

//查看test3的所有版本
// $ cur1 10.29.2.1:12345 /versions /test3
// {"Name":"test3","Version":1,"Size";20, "Hash":"GYggAdrPt+CSenuDc0Gcu3kwCWmOADKNYpiztdlgsM="}
// {"Name":"test3","Version":2,"Size":30,"Hash":"cAPVSxZe1PR5AzIESQYOBaxC1PYJIvaHSF3qEOZYYIo="}

//删除对象test3
//＄ curl -v 10.29.2.1:12345/objects/test3 -XDELETE

//查看test3所有版本
// curl 10.29.2.1:12345/versions /test3
// "Name":"test3", wVersion":1,"Size":20,"Hash":"GYqqAdFPt+CScnUDcGcB3KwCWInOADKNYpiztabgsM-}
// {"Name":"test3","Version":2,"Size":30,"Eash":"cAPVSxZe1PR5AZIESQyOB aXC1pYJIVaHST3gEOZYYIo="}
// {"Name":"test3","Version":3,"Size":0, "Hash":""}
//仍可以获取以前的版本