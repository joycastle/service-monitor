cp -r /home/ec2-user/cert ./

/usr/local/lib/go/bin/go build

/usr/bin/nohup ./service-monitor
