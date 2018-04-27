# echo-on-gae
LabStack Echo on Google App Engine Sample


```
// local
$ dev_appserver.py backend/app.yaml

// gcloudを使用する場合
$ GOPATH=`pwd`/gopath gcloud app deploy backend/app.yaml

// appcfg を使用する場合
$ appcfg.py update --application=YOUR_APP_ID --version=1 --oauth2_access_token=$(gcloud auth print-access-token 2> /dev/null) backend/app.yaml
```