package k8s

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

	ip, err := podsMap.getPodIp(pod)
	if err != nil {
		return err
	}

	podsMap[pod.Name] = &PodInfo{
		Ip:    ip,
		Ready: podsMap.isReadyPhase(pod.Status.Phase),
	}

	return nil
}

func (podsMap PodsMap) getPodIp(pod *corev1.Pod) (netip.Addr, error) {
	if len(pod.Status.PodIP) == 0 {
		return netip.Addr{}, nil
	}
	parsedIp, err := netip.ParseAddr(pod.Status.PodIP)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("Unable to add pod - IP parsing failed: %s", err.Error())
	}

	return parsedIp, nil
}

func (podsMap PodsMap) Update(pod *corev1.Pod) error {
	if len(pod.Name) == 0 {
		return errors.New("Unable to update pod: pod has no name")
	}

	entry, ok := podsMap[pod.Name]
	if ok != true {
		return errors.New("Unable to update pod: unknown")
	}

	ip, err := podsMap.getPodIp(pod)
	if err != nil {
		return err
	}

	entry.Ip = ip
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
