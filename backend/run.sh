sudo docker run --log-opt max-size=10m --log-opt max-file=5 -d --name sayhi  -p 9996:9996 --restart=always  -v /opt/www/sayhi/backend:/opt/www/sayhi/backend   sayhi:v1
