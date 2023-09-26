// Package port includes primary and secondary ports according to the hexagonal
// architecture. The primary ports should expose the functionalities offered
// by the core of our app (normally through application services). Secondary
// ports should define interfaces to comunicate with the outside world (i.e. to
// comunicate with databases, message brokers, REST services...). DDD
// repository interfaces should be defined here instead of the domain package.
package port
