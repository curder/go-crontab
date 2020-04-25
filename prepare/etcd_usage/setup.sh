#/bin/sh

# update yum mirror
sudo mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
sudo curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

# install some tools
sudo yum install -y git vim gcc glibc-static telnet bridge-utils net-tools wget

sudo yum update

# install docker
curl -fsSL get.docker.com -o get-docker.sh
sh get-docker.sh

# start docker service
sudo systemctl start docker

rm -rf get-docker.sh
