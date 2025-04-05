Go application that orchestrates the communication in between MeshAPI-Messages and MeshAPI-SMS. 

Use the following commands to deploy with Docker:

```
sudo docker pull meshrelay0/meshrelay-meshsms-orchestrator-go
sudo docker run -d -p 3040:3040 --name meshsms-orchestrator-go \
--network network_name \
-e MESSAGES_API_ENDPOINT=messages_api_endpoint \
-e SMS_API_ENDPOINT=sms_api_endpoint \
-e PROCESS_INTERVAL=5000 \
--restart always \
meshrelay0/meshrelay-meshsms-orchestrator-go
```