# etcd

## 启动

```
sudo docker run \
  -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-v3.3.20 \
  gcr.io/etcd-development/etcd:v3.3.20 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
```

## 基本操作

```
docker exec etcd-gcr-v3.3.20 /bin/sh -c "/usr/local/bin/etcd --version" # 查看etcd版本
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl version" # 查看etcdctl版本
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl endpoint health" #
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl put foo bar" # 存值
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get foo" # 取值
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl del foo" # 删除值
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get "/cron/jobs/" --prefix" # 按值前缀匹配对应的KV
```

```
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl watch 'cron/jobs/' --prefix" # 监听前缀变化

# 以下etcd操作会被etcdctl监听到
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl put 'cron/jobs/job1' '{key:value}'" # 更新一个key
docker exec etcd-gcr-v3.3.20 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl del 'cron/jobs/job1'" # 删除一个key
```