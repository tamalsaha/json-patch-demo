package main

import (
	"github.com/evanphx/json-patch"
	"fmt"
	"log"
)

func main() {
	d1 := `{
  "apiVersion": "apps/v1",
  "kind": "Deployment",
  "metadata": {
    "annotations": {
      "deployment.kubernetes.io/revision": "1",
      "docker.io/nginx.git-branch": "master",
      "docker.io/nginx.maintainer": "NGINX Docker Maintainers <docker-maint@nginx.com>",
      "docker.io/nginx.repo-url": "github.com/kubepack/pack-server",
      "docker.io/nginx.version": "1.12"
    },
    "labels": {
      "run": "nginx"
    },
    "name": "nginx",
    "namespace": "default"
  },
  "spec": {
    "replicas": 1,
    "selector": {
      "matchLabels": {
        "run": "nginx"
      }
    },
    "template": {
      "metadata": {
        "labels": {
          "run": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "image": "tigerworks/nginx:1.11",
            "imagePullPolicy": "IfNotPresent",
            "name": "nginx"
          }
        ],
        "imagePullSecrets": [
          {
            "name": "regcred"
          }
        ]
      }
    }
  }
}`
	p1 := `[
  {
    "op": "remove",
    "path": "/metadata/annotations/docker.io~1nginx.maintainer",
    "value": null
  },
  {
    "op": "remove",
    "path": "/metadata/annotations/docker.io~1nginx.repo-url",
    "value": null
  },
  {
    "op": "add",
    "path": "/kind",
    "value": "ReplicaSet"
  },
  {
    "op": "add",
    "path": "/apiVersion",
    "value": "extensions/v1beta1"
  }
]`

	patch, err := jsonpatch.DecodePatch([]byte(p1))
	if err != nil {
		log.Fatal(err)
	}

	o1, err := patch.Apply([]byte(d1))
	fmt.Println(err)
	fmt.Println(string(o1))
}
