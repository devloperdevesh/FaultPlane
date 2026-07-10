// Package storage implements checkpoint persistence.
//
// Current backend:
//
//   - In-memory storage
//
// Future backends:
//
//   - Redis
//   - PostgreSQL
//   - S3
//   - BadgerDB
//
// The control plane depends only on the Store interface.
package storage