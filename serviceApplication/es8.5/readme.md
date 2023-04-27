```editorconfig
--------------------------- Security autoconfiguration information ------------------------------

Authentication and authorization are enabled.
TLS for the transport and HTTP layers is enabled and configured.

The generated password for the elastic built-in superuser is : tvuU6HE3x3a7*N4o5ZZU    # 可能是这个 F0sr9auyTeG8AlLPIWwGBA

If this node should join an existing cluster, you can reconfigure this with
'/usr/share/elasticsearch/bin/elasticsearch-reconfigure-node --enrollment-token <token-here>'
after creating an enrollment token on your existing cluster.

You can complete the following actions at any time:

Reset the password of the elastic built-in superuser with
'/usr/share/elasticsearch/bin/elasticsearch-reset-password -u elastic'.

Generate an enrollment token for Kibana instances with
 '/usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s kibana'.

Generate an enrollment token for Elasticsearch nodes with
'/usr/share/elasticsearch/bin/elasticsearch-create-enrollment-token -s node'.

-------------------------------------------------------------------------------------------------
### NOT starting on installation, please execute the following statements to configure elasticsearch service to start automatically using systemd
 sudo systemctl daemon-reload
 sudo systemctl enable elasticsearch.service
### You can start elasticsearch service by executing
 sudo systemctl start elasticsearch.service

```



/etc/systemd/system/multi-user.target.wants/elasticsearch.service → /lib/systemd/system/elasticsearch.service.

/usr/share/elasticsearch/bin/systemd-entrypoint -p /var/run/elasticsearch/elasticsearch.pid --quiet



```editorconfig

✅ Elasticsearch security features have been automatically configured!
✅ Authentication is enabled and cluster connections are encrypted.



ℹ️  HTTP CA certificate SHA-256 fingerprint:
  659d8bde033e09cdedc8d4ea8e88490c419d95d64d52bb98c549fac8278ba797

ℹ️  Configure Kibana to use this cluster:
• Run Kibana and click the configuration link in the terminal when Kibana starts.
• Copy the following enrollment token and paste it into Kibana in your browser (valid for the next 30 minutes):
  eyJ2ZXIiOiI4LjUuMyIsImFkciI6WyIxNzIuMTguMTgyLjU2OjkyMDEiXSwiZmdyIjoiNjU5ZDhiZGUwMzNlMDljZGVkYzhkNGVhOGU4ODQ5MGM0MTlkOTVkNjRkNTJiYjk4YzU0OWZhYzgyNzhiYTc5NyIsImtleSI6IkN5WTRQb1VCLW5vYVU3REtqM19mOnFObmNNOW9ZUnQtWDBrcUFRd3Q0YUEifQ==

ℹ️  Configure other nodes to join this cluster:
• On this node:
  ⁃ Create an enrollment token with `bin/elasticsearch-create-enrollment-token -s node`.
  ⁃ Uncomment the transport.host setting at the end of config/elasticsearch.yml.
  ⁃ Restart Elasticsearch.
• On other nodes:
  ⁃ Start Elasticsearch with `bin/elasticsearch --enrollment-token <token>`, using the enrollment token that you generated.

```



在 sudo vim /etc/sysctl.conf文件最后添加一行
vm.max_map_count=262144
sudo /sbin/sysctl -p

sudo vim /etc/sysctl.conf

zr soft nofile 65535
zr hard nofile 65537




bootstrap check failure [1] of [2]: max file descriptors [4096] for elasticsearch process is too low, increase to at least [65535]
bootstrap check failure [2] of [2]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]





bootstrap check failure [1] of [2]: max number of threads [2048] for user [zr] is too low, increase to at least [4096]
bootstrap check failure [2] of [2]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
ERROR: Elasticsearch did not exit normally - check the logs at /var/log/elasticsearch/es8.5-zr.log

ERROR: [2] bootstrap checks failed. You must address the points described in the following [2] lines before starting Elasticsearch.

sudo vim /etc/security/limits.d/90-nproc.conf

* soft nproc 4096

/usr/share/elasticsearch/bin/systemd-entrypoint -p /home/zr/sdk/elasticsearch.pid --quiet
