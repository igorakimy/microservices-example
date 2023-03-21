Init docker swarm
```bash
docker swarm init
```

Deploy docker swarm
```bash
docker stack deploy -c swarm.yml micro
```

Watch list of services
```bash
docker service ls
```

Pull new version of service
```bash
docker pull igorakimov/logger-service:1.0.1
```

Run several replicas of service
```bash
docker service scale micro_listener=2
```

Update service in docker swarm stack to a new version
```bash
docker service update --image igorakimov/logger-service:1.0.1 micro_logger
```
