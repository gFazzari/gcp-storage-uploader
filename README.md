# Simple Google Storage Uploader written in Go

## Compile instruction

```
# Golang 1.19
go build .
```

## Help

```
./storage_uploader -h
```

## Goal
This Golang program has been created to upload files into GCS.  
It's actually used to upload Argo CD database backup and Bitnami Sealed Secrets encryption keys.  
You can pass project ID, bucket name, filename to upload and service account, as arguments.