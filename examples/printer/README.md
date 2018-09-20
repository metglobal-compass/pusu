## Summary

This example is a simulation of pub/sub workflow for Google Cloud in local environment.

## Building

Use docker-compose for building the example


```docker-compose up --build```



Send following curl request to Google Cloud Pub/Sub Emulator

```
curl -X POST \
  http://localhost:8085/v1/projects/my-project-id/topics/printer:publish \
  -H 'content-type: application/json' \
  -d '{"messages": [
    {
    "data":"YXNkYWY=",
    "messageId":"32141234234234234"
    }
  ]
}'
```
