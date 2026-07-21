package balancer

import "net/netip"

type PodsMap map[string]netip.Addr


// Add an IP address for the Pod by name
func (podsMap PodsMap) AddAddress(name string, ip string) error {
	return nil
}

// Update a Pod name, but keep the associated IP address
func (podsMap PodsMap) UpdateName(oldName string, newName string) error {
	return nil
}

// Update an IP address associated with a Pod
func (podsMap PodsMap) UpdateAddress(name string, ip string) error {
	return nil
}

// Delete a Pod by name
func (podsMap PodsMap) DropName(name string) error {
	return nil
}

// Delete a Pod by an associated IP address
func (podsMap PodsMap) DropAddress(ip string) error {
	return nil
}


func (podsMap PodsMap) hasPod(name string) bool {
	_, ok := podsMap[name]
	
	return ok
}

func (podsMap PodsMap) hasIp(ip string) bool {
	return false
}
