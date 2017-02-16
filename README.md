`$ docker build -t rithium/gateway .`

`$ docker run -it -e ZK_HOSTS=172.17.0.2:2181 -p 11000:80 rithium/gateway`

