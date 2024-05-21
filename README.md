# k8s

k8s operator 학습

## Description

https://youtu.be/ND4haK4pDF4?feature=shared 를 따라하며 작성했습니다.

## Getting Started

### Prerequisites

- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

---

### 1. 스켈레톤 프로젝트 생성

```sh
kubectll get no # 노드 상태 확인
```

```sh
# go 종속성 관리 파일 생성
go mod init markrluer.com
cat go mod
```

```sh
# 스켈레톤 프로젝트 생성
kubebuilder init --domain markruler.com
kubebuilder create api --group rmk --version v1alpha1 --kind Machine
# resource y, controller y 로 같이 생성
```

### 2. API Resource 정의

- types.go 수정
- controller.go 수정

```sh
make install # crd 배포
kubectl get crd # 상태 확인
```

```sh
kubectl get machine # 머신은 아직 없어야 함
```

- {group}_{versoin}_{domain}.yaml 수정

```sh
kubectl apply -f config/samples/rmk_v1alpha1_machine.yaml # 파일 수정 적용
kubectl get machine # 머신 조회 --> 존재해야 함
```

### 3. Reconcile Loop 개발

- controller.go 수정
- rmk_v1alpha1_machine.yaml 수정

```sh
make run
kubectl get machine
```

![image](https://github.com/vo0a/k8s_operator_study/assets/44438366/d39d909a-e3cc-49a3-ba67-2d30f6505ddf)

### 4. 컨트롤러 컨테이너 배포

- 사전 준비: 도커 계정 설정

```sh
sudo docker login
```

- docker hub 에 이미지 업로드

```sh
sudo make docker-build docker-push IMG=vo0a/operator-markruler:0.1.0
sudo docker images | grep vo0a
make deploy IMG=vo0a/operator-markruler:0.1.0 # 배포
```

```sh
# 상태 바꾸고, 변경이 감지되어 machine 상태가 바뀌는지 확인
kubectl get po -A
# 여기서 나온 컨트롤러 이름으로 log 확인
# k8s-controller-manager-bdb68d44c-887p4

# 여기서 role garbage 로 변경
watch kubectl get machine
{group}_{versoin}_{domain}.yaml 수정

# 상태 확인용
kubectl logs -f k8s-controller-manager-bdb68d44c-887p4 -n k8s-system -c manager

# garbage 상태로 재배포
sudo kubectl apply -f config/samples/rmk_v1alpha1_machine.yaml
```

pod list

![image](https://github.com/vo0a/k8s_operator_study/assets/44438366/16755177-0e21-48ce-ac52-69a973c64764)

재배포 결과 - 삭제됨

![image](https://github.com/vo0a/k8s_operator_study/assets/44438366/fb9725f2-fe70-4a50-a90c-0ad72cb9fc06)

영상
https://github.com/vo0a/k8s_operator_study/assets/44438366/d27a9c0d-88f2-4684-9c6d-4389c0d69138

> todo: custom controller 로 log -f 변경

---

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/k8s:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/k8s:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

> **NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/k8s:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/k8s/<tag or branch>/dist/install.yaml
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
