module dslab.sjtu/kube-dtn-sidecar

go 1.22.3

require (
	dslab.sjtu/kube-dtn/api v0.0.0-00010101000000-000000000000
	github.com/digitalocean/go-openvswitch v0.0.0-20240130171624-c0f7d42efe24
	github.com/krolaw/dhcp4 v0.0.0-20190909130307-a50d88189771
	github.com/spf13/pflag v1.0.5
	google.golang.org/grpc v1.65.0
	k8s.io/apimachinery v0.30.2
	kubevirt.io/api v1.2.0
	kubevirt.io/client-go v1.2.0
	kubevirt.io/kubevirt v1.2.0
)

replace dslab.sjtu/kube-dtn/api => github.com/dtn-dslab/kube-dtn-api v0.0.5

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f // indirect
	golang.org/x/crypto v0.23.0 // indirect
)

require (
	github.com/containernetworking/cni v1.1.2
	github.com/digitalocean/go-libvirt v0.0.0-20240308204700-df736b2945cf
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/openshift/api v0.0.0 // indirect
	github.com/openshift/custom-resource-status v1.1.2 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/api v0.30.2 // indirect
	k8s.io/apiextensions-apiserver v0.28.1 // indirect
	k8s.io/klog/v2 v2.120.1 // indirect
	k8s.io/utils v0.0.0-20230726121419-3b25d923346b // indirect
	kubevirt.io/containerized-data-importer-api v1.57.0-alpha1 // indirect
	kubevirt.io/controller-lifecycle-operator-sdk/api v0.0.0-20220329064328-f3cc58c6ed90 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/openshift/api => github.com/openshift/api v0.0.1
