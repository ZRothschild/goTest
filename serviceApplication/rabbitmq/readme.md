
docker  run -p 15671:15671 -p 15672:15672 -p 15691:15691 -p 15692:15692 -p 25672:25672 -p 4369:4369 -p 5671:5671 -p 5672:5672 -d --hostname zr-rabbit --name rabbit -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest -e RABBITMQ_DEFAULT_VHOST=my_vhost rabbitmq:management
