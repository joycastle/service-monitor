cp -r /home/ec2-user/cert ./

go build

nohub ./service-monitor
