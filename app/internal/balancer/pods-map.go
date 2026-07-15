package balancer

import "net/netip"

type PodsMap map[string]netip.Addr

// Add an IP address for the Pod by name
func (podsMap *PodsMap) AddAddress(name string, ip string) {

}

// Update a Pod name, but keep the associated IP address
func (podsMap *PodsMap) UpdateName(oldName string, newName string) {

}

// Update an IP address associated with a Pod
func (podsMap *PodsMap) UpdateAddress(name string, ip string) {

}

// Delete a Pod by name
func (podsMap *PodsMap) DropByName(name string) {

}

// Delete a Pod by an associated IP address
func (podsMap *PodsMap) DropByAddress(ip string) {

}
