# Jenkins node state exporter

Prometheus exporter for Jenkins nodes

This exporter listen on port `9723` and the endpoint is `/metrics`

```
# HELP jenkins_node_busy_status The busy status of a jenkins computer node 0:IDLE 1:BUSY
# TYPE jenkins_node_busy_status gauge
jenkins_node_busy_status{computerName="master"} 0
jenkins_node_busy_status{computerName="node1"} 0
# HELP jenkins_node_exporter_failure The number of faillure to get/parse api data
# TYPE jenkins_node_exporter_failure counter
jenkins_node_exporter_failure 0
# HELP jenkins_node_maintenance_status The maintenance status of a jenkins computer node 0:ONLINE 1:MAINTENANCE 2:OFFLINE
# TYPE jenkins_node_maintenance_status gauge
jenkins_node_maintenance_status{computerName="master"} 0
jenkins_node_maintenance_status{computerName="node1"} 2
```

Available options:
```
Usage:
  jenkins-node-state-exporter [flags]

Flags:
  -a, --address string            address of the jenkins server (default "http://localhost/jenkins")
      --disable-authentication    disable authentication
      --fetch-interval duration   fetch-interval in seconds (default 30s)
  -h, --help                      help for jenkins-node-state-exporter
  -p, --password string           password of the jenkins user account (default "admin")
      --port int                  port to listen on (default 9827)
  -u, --username string           username of the jenkins user account (default "admin")
```
