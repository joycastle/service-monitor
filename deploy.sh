cp -r /home/ec2-user/cert ./
docker build -f Dockerfile -t service-monitor:latest .
docker run -d --restart=always \
    --name service-monitor \
    -v /home/ec2-user/var/service-monitor/logs:/app/var/logs \
    service-monitor:latest
