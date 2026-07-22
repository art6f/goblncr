package balancer

import (
	"errors"
	"fmt"
	"net/netip"

	corev1 "k8s.io/api/core/v1"
)

type PodInfo struct {
	Ip    netip.Addr
	Ready bool
}

type PodsMap map[string]*PodInfo

func (podsMap PodsMap) Add(pod *corev1.Pod) error {

	if len(pod.Name) == 0 {
		return errors.New("Unable to add pod: pod has no name")
	}

	if len(pod.Status.PodIP) == 0 {
		return errors.New("Unable to add pod: pod has no IP address")
	}

	ip, err := netip.ParseAddr(pod.Status.PodIP)
	if err != nil {
		return fmt.Errorf("Unable to add pod - IP parsing failed: %s", err.Error())
	}

	podsMap[pod.Name] = &PodInfo{
		Ip:    ip,
		Ready: podsMap.isReadyPhase(pod.Status.Phase),
	}

	return nil
}

func (podsMap PodsMap) Update(pod *corev1.Pod) error {
	if len(pod.Name) == 0 {
		return errors.New("Unable to update pod: pod has no name")
	}

	entry, ok := podsMap[pod.Name]
	if ok != true {
		return errors.New("Unable to update pod: unknown")
	}

	entry.Ready = podsMap.isReadyPhase(pod.Status.Phase)

	return nil
}

func (podsMap PodsMap) Delete(pod *corev1.Pod) error {
	if len(pod.Name) == 0 {
		return errors.New("Unable to delete pod: pod has no name")
	}

	if _, ok := podsMap[pod.Name]; ok != true {
		return errors.New("Unable to delete pod: unknown")
	}

	delete(podsMap, pod.Name)

	return nil
}

func (podsMap PodsMap) isReadyPhase(phase corev1.PodPhase) bool {
	return phase == corev1.PodRunning
}

//// Delete a Pod by an associated IP address
//func (podsMap PodsMap) DropByName(ip name) error {
//	return nil
//}
//
//// Delete a Pod by an associated IP address
//func (podsMap PodsMap) DropByAddress(ip string) error {
//	return nil
//}
//
//func (podsMap PodsMap) hasPod(name string) bool {
//	_, ok := podsMap[name]
//
//	return ok
//}
//
//func (podsMap PodsMap) hasIp(ip string) bool {
//	return false
//}
