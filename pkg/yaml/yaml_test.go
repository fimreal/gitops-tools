package yaml

import (
	"bytes"
	"os"
	"testing"

	"github.com/fimreal/goutils/ezap"
	mfile "github.com/fimreal/goutils/file"
)

const clash = `apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: clash-proxy
  name: clash-proxy
  namespace: infra
spec:
  replicas: 1
  selector:
    matchLabels:
      app: clash-proxy
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: clash-proxy
        tier: srv
    spec:
      # affinity:
      #   nodeAffinity:
      #     preferredDuringSchedulingIgnoredDuringExecution:
      #       - preference:
      #           matchExpressions:
      #           - key: kubernetes.io/hostname
      #             operator: In
      #             values:
      #               - "home"
      #         weight: 10
      #       - preference:
      #           matchExpressions:
      #           - key: kubernetes.io/hostname
      #             operator: In
      #             values:
      #               - "gw"
      #         weight: 5
      #   podAntiAffinity:
      #     requiredDuringSchedulingIgnoredDuringExecution:
      #     - labelSelector:
      #         matchExpressions:
      #         - key: app
      #           operator: In
      #           values:
      #           - clash-proxy
      #       topologyKey: "kubernetes.io/hostname"
      # tolerations:
      #   - key: node-role.kubernetes.io/master
      #     effect: NoSchedule
      initContainers:
        - name: create-white-ip-file
          image: epurs/awk
          args:
            - sh
            - -c
            - awk '@include "inplace";!a[$0]++' /etc/nginx/whiteip.txt
          volumeMounts:
            - name: clash-config
              mountPath: /etc/nginx
              subPath: clash-config
      containers:
        - name: clash-proxy
          image: epurs/openresty
          imagePullPolicy: Always
          ports:
            - containerPort: 80
              name: http
            - containerPort: 1080
              name: http-proxy
            - containerPort: 1081
              name: socks5-proxy
          resources:
            limits:
              cpu: 10m
              memory: 50Mi
            requests:
              cpu: 10m
              memory: 30Mi
          readinessProbe:
            tcpSocket:
              port: 80
            initialDelaySeconds: 10
            periodSeconds: 30
            timeoutSeconds: 5
            successThreshold: 1
            failureThreshold: 3
          volumeMounts:
            - name: http
              mountPath: "/etc/nginx/conf.d"
            - name: user-root
              mountPath: "/etc/nginx/global"
            - name: clash-stream-proxy
              mountPath: "/etc/nginx/stream"
            - name: clash-config
              mountPath: /etc/nginx
              subPath: clash-config
      volumes:
        - name: clash-config
          persistentVolumeClaim:
            claimName: clash-config
        - name: http
          configMap:
            name: ngx-conf
            items:
              - key: dyna-ip.conf
                path: dyna-ip.conf
        - name: user-root
          configMap:
            name: ngx-conf
            items:
              - key: user-root.conf
                path: user-root.conf
        - name: clash-stream-proxy
          configMap:
            name: ngx-conf
            items:
              - key: clash-stream-proxy.conf
                path: clash-stream-proxy.conf
---
# ---------------  svc -------------- ##
apiVersion: v1
kind: Service
metadata:
  name: clash-proxy
  namespace: netio
spec:
  type: ClusterIP
  ports:
    - name: http
      port: 80
      targetPort: 80
  selector:
    app: clash-proxy
---
# ---------------  svc -------------- ##
apiVersion: v1
kind: Service
metadata:
  name: hproxy
  namespace: netio
spec:
  type: NodePort
  externalTrafficPolicy: Local
  ports:
    - name: hproxy
      port: 1080
      targetPort: 1080
      nodePort: 6789
    - name: socks5
      port: 1081
      targetPort: 1081
      nodePort: 6780
  selector:
    app: clash-proxy
---
# ---------------  ing -------------- ##
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
  name: clash-proxy
  namespace: netio
spec:
  rules:
    - host: "hproxy.epurs.com"
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: clash-proxy
                port:
                  number: 80
  tls:
    - secretName: epurs-com
`
const pgsql = `
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
  namespace: infra
spec:
  replicas: 1
  serviceName: "postgres"
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:14.3-alpine
          env:
          # - name: POSTGRES_DB
          #   valueFrom:
          #     secretKeyRef:
          #       name: postgres-config
          #       key: POSTGRES_DB
          #       optional: true
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: postgres-config
                  key: POSTGRES_USER
                  optional: false
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: postgres-config
                  key: POSTGRES_PASSWORD
                  optional: false
          ports:
            - containerPort: 5432
              name: postgredb
          resources:
            limits:
              cpu: 500m
              memory: 500Mi
            requests:
              cpu: 300m
              memory: 100Mi
          volumeMounts:
            - name: postgres-data
              mountPath: /var/lib/postgresql/data
              subPath: postgres
  volumeClaimTemplates:
    - metadata:
        name: postgres-data
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: longhorn
        resources:
          requests:
            storage: 5Gi
---

apiVersion: v1
kind: Service
metadata:
  name: postgres
  namespace: infra
  labels:
    app: postgres
spec:
  ports:
    - port: 5432
      targetPort: 5432
      name: postgres
  type: ClusterIP
  selector:
    app: postgres`

func init() {
	os.WriteFile("clash.yaml", []byte(clash), 0644)
	os.WriteFile("postgres.yaml", []byte(pgsql), 0644)
}

var testFiles = []string{"clash.yaml", "postgres.yaml"}

func TestWhatKindOf(t *testing.T) {
	for _, f := range testFiles {
		b, err := os.ReadFile(f)
		if err != nil {
			ezap.Fatal(err)
		}

		dataArr := bytes.Split(b, []byte("---\n"))
		for _, desc := range dataArr {
			t.Log(WhatKindOf(desc))
		}
	}
}

func TestSplitYamlFile(t *testing.T) {
	for _, f := range testFiles {

		k8syaml, err := SplitYamlFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, y := range k8syaml {
			t.Log(y.Kind)
			t.Log(string(y.ByteData))
		}
	}
}

func TestUpdateImage(t *testing.T) {
	// ezap.SetLevel("debug")
	for _, f := range testFiles {

		k8syaml, err := SplitYamlFile(f)
		if err != nil {
			t.Fatal(err)
		}
		mfile.WriteToFile(f, nil)
		for _, y := range k8syaml {
			err := y.UpdateImage("docker.io/epurs/openresty:latest", "")
			if err != nil {
				t.Fatal(err)
			}
			mfile.AppendToFile(f, y.ByteData)
			mfile.AppendToFile(f, []byte("---\n"))
		}
	}
}
