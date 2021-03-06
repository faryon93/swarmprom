package swarmprom

// swarmprom - prometheus http wrapper for swarm services
// Copyright (C) 2018 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"net"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

// Resolver for swarm service tasks.
type swarmResolver struct {
	cache map[string][]string
}

// A list of IP addresses.
type ipList []string

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// Returns the IP addresses of all containers beloning to the given service.
func (s *swarmResolver) GetServiceIps(service string) (ipList, error) {
	if s.cache == nil {
		s.cache = make(map[string][]string)
	}

	list, err := net.LookupHost("tasks." + service)
	if err != nil {
		list = s.cache[service]
		if list == nil {
			return nil, err
		}
	} else {
		s.cache[service] = list
	}

	return list, nil
}

// Contains returns true if the given ip is on this list.
func (l ipList) Contains(ip string) bool {
	for _, haystack := range l {
		if ip == haystack {
			return true
		}
	}

	return false
}