/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package resource

import (
	"io"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/klog"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// Parse - 지정한 정보를 기준으로 k8s 리소스 정보 반환
// conditions:
func Parse(reader io.Reader, namespace, releaseName string, local bool) ([]*resource.Info, error) {
	// klog 2 Level의 로그 출력 및 종료시 성능을 포함한 로그 출력
	if klog.V(2) {
		klog.Infof("parse resources, namespace: %s, release: %s", namespace, releaseName)
		start := time.Now()
		defer func() {
			klog.Infof("parse resource end, namespace: %s, release: %s, cost: %v", namespace, releaseName, time.Now().Sub(start))
		}()
	}

	// 기본 k8s cli 플래그 설정
	kubeConfigFlags := genericclioptions.NewConfigFlags(true)
	// 버전에 맞는 플래그 설정
	matchVersionKubeConfigFlags := cmdutil.NewMatchVersionFlags(kubeConfigFlags)
	// Flag 기준의 Factory 생성
	f := cmdutil.NewFactory(matchVersionKubeConfigFlags)
	// Factory를 이용해서 Namespace, ReleaseName 기준으로 Builder 생성
	builder := f.NewBuilder().Unstructured().NamespaceParam(namespace).ContinueOnError().Stream(reader, releaseName).Flatten()

	// local 처리인 경우는 local builder 구성
	if local {
		builder = builder.Local()
	}

	// Builder 실행
	r := builder.Do()
	// 처리 결과에서 Resource 들에 대한 Infos 추출
	infos, err := r.Infos()
	if err != nil {
		return nil, err
	}

	if !local {
		for i := range infos {
			infos[i].Namespace = namespace
			// infos 내부 오류 여부 검증
			err := infos[i].Get()
			if err != nil {
				return nil, err
			}
		}
	}

	return infos, err
}
