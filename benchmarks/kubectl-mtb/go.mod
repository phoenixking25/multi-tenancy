module sigs.k8s.io/multi-tenancy/benchmarks/kubectl-mtb

go 1.13

require (
	github.com/creasty/defaults v1.3.0
	github.com/gobuffalo/packr/v2 v2.5.1 // indirect
	github.com/google/logger v1.1.0 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kubernetes-sigs/multi-tenancy v0.0.0-20200630205007-c76bb22a12b8
	github.com/onsi/gomega v1.10.1
	github.com/phoenixking25/kubectl-mtb v0.0.0-20200607051717-90ca51647eb6 // indirect
	github.com/spf13/cobra v1.0.0
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/cli-runtime v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/kubectl v0.0.0
	k8s.io/kubernetes v1.18.3
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
	sigs.k8s.io/kind v0.8.1
)

replace (
	github.com/kubernetes-sigs/multi-tenancy => github.com/Divya063/multi-tenancy v0.0.0-20200703122503-2307447baef6
	k8s.io/api => k8s.io/api v0.18.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.3
	k8s.io/apiserver => k8s.io/apiserver v0.18.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.3
	k8s.io/client-go => k8s.io/client-go v0.18.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.3
	k8s.io/code-generator => k8s.io/code-generator v0.18.3
	k8s.io/component-base => k8s.io/component-base v0.18.3
	k8s.io/cri-api => k8s.io/cri-api v0.18.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.3
	k8s.io/gengo => k8s.io/gengo v0.0.0-20200114144118-36b2048a9120
	k8s.io/heapster => k8s.io/heapster v1.2.0-beta.1
	k8s.io/klog => k8s.io/klog v1.0.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.3
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200410145947-61e04a5be9a6 // release-1.18
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.3
	k8s.io/kubectl => k8s.io/kubectl v0.18.3
	k8s.io/kubelet => k8s.io/kubelet v0.18.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.3
	k8s.io/metrics => k8s.io/metrics v0.18.3
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.1-alpha.1
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.3
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.18.3
	k8s.io/sample-controller => k8s.io/sample-controller v0.18.3
	k8s.io/system-validators => k8s.io/system-validators v1.0.4
	k8s.io/utils => k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
)
