package controllers

import "net/http"

func (c Controller) Protected() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c.logger.Println("protected endpoint invoked")
	}
}