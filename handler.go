package swarmprom

// swarmprom - prometheus http wrapper for swarm services
// Copyright (C) 2018 Maximilian Pachl

// MIT License
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"errors"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// ---------------------------------------------------------------------------------------
//  global variables
// ---------------------------------------------------------------------------------------

var (
	logger        = logrus.NewEntry(logrus.StandardLogger())
	rejectHandler = DefaultRejectHandler
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// SetLogger sets the new logger with should be used with.
func SetLogger(newLogger *logrus.Entry) {
	logger = newLogger
}

// SetRejectHandler sets the HTTP handler which is
// executed when the access is restricted.
func SetRejectHandler(fn http.HandlerFunc) {
	if fn == nil {
		return
	}

	rejectHandler = fn
}

// Handler returns a prometheus http hander which can be accessed
// by containers of the given service only.
// If the supplied name is empty the access control is bypassed
// an a warning is printed through the logger.
func Handler(service string) http.Handler {
	resolver := swarmResolver{}
	handler := promhttp.Handler()

	// if an empty service name is supplied
	// bypass the swarm service access control.
	if service == "" {
		if logger != nil {
			logger.Warnln("PROMETHEUS METRIC ENDPOINT ACCESS CONTROL IS BYPASSED")
		}

		return handler
	}

	// construct the middleware handler function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we only want the ip address of the requesting client
		remote := r.RemoteAddr
		if host, _, err := net.SplitHostPort(remote); err == nil {
			remote = host
		}

		// setup the logger
		var log *logrus.Entry
		if logger != nil {
			log = logger.WithField("addr", remote)
		}

		// resolve all container ips of the given service
		validationError := errors.New("not on swarm service list")
		allowedIps, err := resolver.GetServiceIps(service)
		if err != nil {
			validationError = err
		} else {
			// the calling client must be on the white list
			for _, ip := range allowedIps {
				if ip == remote {
					validationError = nil
					break
				}
			}
		}

		// reject access to the endpoint
		if validationError != nil {
			rejectHandler(w, r)
			if log != nil {
				log.Warnln("rejecting access to promhttp handler:", validationError.Error())
			}
			return
		}

		// execute the actual http handler
		handler.ServeHTTP(w, r)
	})
}

// DefaultRejectHandler is the standard reject handler.
// It returns a "forbidden" message with code 403.
func DefaultRejectHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "forbidden", http.StatusForbidden)
}
