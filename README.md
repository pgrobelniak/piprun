# piprun 

Synchronous buddy pipeline runner.

This app is meant to be deployed to sandbox. 

When deployed it allows to run pipelines synchronously using get request, such as

```
curl https://SANDBOX-URL.buddy.cloud/?token=API-TOKEN&workspace=WORKSPACE-NAMEproject=PROJECT-NAME&pipeline=PIPELINE-ID
```

Response will contain execution status.
