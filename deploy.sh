cp -r /home/ec2-user/cert ./

/usr/local/lib/go/bin/go build

ps -ef|grep '\./service-monitor'|grep -v grep|awk  '{print $2}'|xargs kill -9

/usr/bin/nohup ./service-monitor
