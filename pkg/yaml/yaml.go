package k8syaml

import (
	"bytes"
	"os"

	"github.com/fimreal/goutils/ezap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	yaml "sigs.k8s.io/yaml"
)

//	type workload[T k8s.Deployment | k8s.DaemonSet | k8s.StatefulSet] interface {
//		ToStruct()
//	}
// type workload interface {
// 	v1.Deployment | v1.DaemonSet | v1.StatefulSet

// 	ToYaml()
// }

type K8sYaml struct {
	Kind     string
	ByteData []byte
}

type PodSpec corev1.PodSpec

func SplitYamlFile(filename string) ([]K8sYaml, error) {
	dataByte, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dataArr := bytes.Split(dataByte, []byte("---\n"))

	var yamlData []K8sYaml
	for _, v := range dataArr {
		kind, err := WhatKindOf(v)
		if err != nil {
			return nil, err
		}
		k8syaml := K8sYaml{Kind: kind, ByteData: v}
		yamlData = append(yamlData, k8syaml)
	}
	return yamlData, nil
}

func WhatKindOf(yaml []byte) (kind string, err error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, k, err := decode(yaml, nil, nil)
	if err != nil {
		ezap.Fatal(err)
	}
	ezap.Debug(obj)

	return k.Kind, err
}

func (y *K8sYaml) UpdateImage(image string, containerName string) error {
	switch y.Kind {
	case "Deployment":
		workload := &appsv1.Deployment{}
		err := yaml.Unmarshal(y.ByteData, workload)
		if err != nil {
			return err
		}
		ezap.Debug("解析 yaml 成功")

		podspec := PodSpec(workload.Spec.Template.Spec)
		err = podspec.UpdateImage(image, containerName)
		if err != nil {
			return err
		}

		// 修改后内容写回 yaml
		y.ByteData, err = yaml.Marshal(workload)
		if err != nil {
			return err
		}
		return nil
	case "StatefulSet":
		workload := &appsv1.StatefulSet{}
		err := yaml.Unmarshal(y.ByteData, workload)
		if err != nil {
			return err
		}
		ezap.Debug("解析 yaml 成功")
		podspec := PodSpec(workload.Spec.Template.Spec)
		err = podspec.UpdateImage(image, containerName)
		if err != nil {
			return err
		}

		// 修改后内容写回 yaml
		y.ByteData, err = yaml.Marshal(workload)
		if err != nil {
			return err
		}
		return nil
	case "DaemonSet":
		var workload appsv1.DaemonSet
		err := yaml.Unmarshal(y.ByteData, workload)
		if err != nil {
			return err
		}

		podspec := PodSpec(workload.Spec.Template.Spec)
		err = podspec.UpdateImage(image, containerName)
		if err != nil {
			return err
		}

		// 修改后内容写回 yaml
		y.ByteData, err = yaml.Marshal(workload)
		if err != nil {
			return err
		}
		return nil
	default:
		ezap.Debugf("Skip not support kind: %s", y.Kind)
		return nil
	}
}

func (p *PodSpec) UpdateImage(image string, containerName string) error {
	c := p.Containers
	if len(c) < 1 {
		ezap.Fatal("传入文件格式有误，请查证后重试")
	} else if len(c) > 1 {
		if containerName == "" {
			ezap.Fatal("传入文件包含多个容器配置，请指定需要修改的容器名称")
		}
		for i := range c {
			if c[i].Name == containerName {
				ezap.Infof("旧容器[%s]镜像为: %s", c[i].Name, c[i].Image)
				c[i].Image = image
				ezap.Infof("将容器[%s]镜像更新为: %s", c[i].Name, c[i].Image)
				return nil
			}
		}
	}

	ezap.Infof("旧容器[%s]镜像为: %s", c[0].Name, c[0].Image)
	c[0].Image = image
	ezap.Infof("将容器[%s]镜像更新为: %s", c[0].Name, c[0].Image)
	return nil
}