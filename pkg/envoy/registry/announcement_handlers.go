package registry

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/k8s/events"
)

// ReleaseCertificateHandler releases certificates based on podDelete events
func (pr *ProxyRegistry) ReleaseCertificateHandler(certManager certificateReleaser, stop <-chan struct{}) {
	podDeleteChan, unsub := pr.msgBroker.SubscribeKubeEvents(events.Pod.Deleted())
	defer unsub()

	for {
		select {
		case <-stop:
			return

		case podDeletedMsg := <-podDeleteChan:
			psubMessage, castOk := podDeletedMsg.(events.PubSubMessage)
			if !castOk {
				log.Error().Msgf("Error casting to events.PubSubMessage, got type %T", psubMessage)
				continue
			}

			// guaranteed can only be a PodDeleted event
			deletedPodObj, castOk := psubMessage.OldObj.(*corev1.Pod)
			if !castOk {
				log.Error().Msgf("Error casting to *corev1.Pod, got type %T", deletedPodObj)
				continue
			}

			proxyUUID := deletedPodObj.Labels[constants.EnvoyUniqueIDLabelName]
			if proxy := pr.GetConnectedProxy(proxyUUID); proxy != nil {
				log.Warn().Msgf("Pod with label %s: %s found in proxy registry; releasing certificate for proxy %s", constants.EnvoyUniqueIDLabelName, proxyUUID, proxy.Identity)
				certManager.ReleaseCertificate(envoy.NewXDSCertCNPrefix(proxy.UUID, proxy.Kind(), proxy.Identity))
			} else {
				log.Warn().Msgf("Pod with label %s: %s not found in proxy registry", constants.EnvoyUniqueIDLabelName, proxyUUID)
			}
		}
	}
}
