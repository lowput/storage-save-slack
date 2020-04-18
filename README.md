# deploy

```sh
gcloud functions deploy storage-save-slack --env-vars-file .env.yaml --entry-point StorageImageSend --runtime go111 --trigger-event=google.storage.object.finalize --trigger-resource $BAKET_NAME
```
