package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/evanphx/json-patch"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var d2 = `{
  "apiVersion": "extensions/v1beta1",
  "kind": "Deployment",
  "metadata": {
    "annotations": {
      "deployment.kubernetes.io/revision": "2",
      "docker.io-nginx-git-branch": "master",
      "docker.io-nginx-maintainer": "NGINX Docker Maintainers <docker-maint@nginx.com>",
      "docker.io-nginx-repo-url": "github.com/kubepack/pack-server",
      "docker.io-nginx-version": "1.11",
      "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"annotations\":{},\"labels\":{\"run\":\"nginx\"},\"name\":\"nginx\",\"namespace\":\"default\"},\"spec\":{\"replicas\":1,\"selector\":{\"matchLabels\":{\"run\":\"nginx\"}},\"template\":{\"metadata\":{\"labels\":{\"run\":\"nginx\"}},\"spec\":{\"containers\":[{\"image\":\"tigerworks/nginx:1.11\",\"imagePullPolicy\":\"IfNotPresent\",\"name\":\"nginx\"}],\"imagePullSecrets\":[{\"name\":\"regcred\"}]}}}}\n"
    },
    "creationTimestamp": "2018-04-27 15:13:47 UTC",
    "generation": 2,
    "labels": {
      "run": "nginx"
    },
    "name": "nginx",
    "namespace": "default",
    "resourceVersion": "2974",
    "selfLink": "/apis/extensions/v1beta1/namespaces/default/deployments/nginx",
    "uid": "955836b6-4a2d-11e8-b76d-080027dd1212"
  },
  "spec": {
    "progressDeadlineSeconds": 600,
    "replicas": 1,
    "revisionHistoryLimit": 10,
    "selector": {
      "matchLabels": {
        "run": "nginx"
      }
    },
    "strategy": {
      "rollingUpdate": {
        "maxSurge": "25%",
        "maxUnavailable": "25%"
      },
      "type": "RollingUpdate"
    },
    "template": {
      "metadata": {
        "creationTimestamp": null,
        "labels": {
          "run": "nginx"
        }
      },
      "spec": {
        "containers": [
          {
            "image": "tigerworks/nginx:1.11",
            "imagePullPolicy": "IfNotPresent",
            "name": "nginx",
            "resources": {
            },
            "terminationMessagePath": "/dev/termination-log",
            "terminationMessagePolicy": "File"
          }
        ],
        "dnsPolicy": "ClusterFirst",
        "imagePullSecrets": [
          {
            "name": "regcred"
          }
        ],
        "restartPolicy": "Always",
        "schedulerName": "default-scheduler",
        "securityContext": {
        },
        "terminationGracePeriodSeconds": 30
      }
    }
  },
  "status": {
    "availableReplicas": 1,
    "conditions": [
      {
        "lastTransitionTime": "2018-04-27 15:13:50 UTC",
        "lastUpdateTime": "2018-04-27 15:13:50 UTC",
        "message": "Deployment has minimum availability.",
        "reason": "MinimumReplicasAvailable",
        "status": "True",
        "type": "Available"
      },
      {
        "lastTransitionTime": "2018-04-27 15:13:47 UTC",
        "lastUpdateTime": "2018-04-27 15:14:30 UTC",
        "message": "ReplicaSet \"nginx-866c59697d\" is progressing.",
        "reason": "ReplicaSetUpdated",
        "status": "True",
        "type": "Progressing"
      }
    ],
    "observedGeneration": 2,
    "readyReplicas": 1,
    "replicas": 2,
    "unavailableReplicas": 1,
    "updatedReplicas": 1
  }
}`

type Deployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

func main() {
	var d Deployment
	err := json.Unmarshal([]byte(d2), &d)
	if err != nil {
		log.Fatalln(err)
	}

	data, ok := d.Annotations[core.LastAppliedConfigAnnotation]
	if ok {
		original := make(map[string]interface{})
		err := json.Unmarshal([]byte(data), &original)
		if err == nil {
			var m metav1.ObjectMeta
			if v, ok := original["metadata"]; ok {
				m = v.(metav1.ObjectMeta)
			}
			a := map[string]string{}
			m.Annotations = a
			original["metadata"] = m
		}
		if out, err := json.Marshal(original); err == nil {
			d.Annotations[core.LastAppliedConfigAnnotation] = string(out)
		}
	}
}

func main2() {
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
