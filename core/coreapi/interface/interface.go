// Package iface defines IPFS Core API which is a set of interfaces used to
// interact with IPFS nodes.
package iface

import (
	"context"
	"errors"
	"io"
	"time"

	options "github.com/ipfs/go-ipfs/core/coreapi/interface/options"

	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	ipld "gx/ipfs/QmPN7cwmpcc4DWXb4KTB9dNAJgjuPY69h3npsMfhRrQL9c/go-ipld-format"
)

// Path is a generic wrapper for paths used in the API. A path can be resolved
// to a CID using one of Resolve functions in the API.
type Path interface {
	String() string
	Cid() *cid.Cid
	Root() *cid.Cid
	Resolved() bool
}

// TODO: should we really copy these?
//       if we didn't, godoc would generate nice links straight to go-ipld-format
type Node ipld.Node
type Link ipld.Link

type IpnsEntry struct {
	Name  string
	Value Path
}

type Reader interface {
	io.ReadSeeker
	io.Closer
}

// CoreAPI defines an unified interface to IPFS for Go programs.
type CoreAPI interface {
	// Unixfs returns an implementation of Unixfs API
	Unixfs() UnixfsAPI
	Name() NameAPI
	Key() KeyAPI

	// ResolvePath resolves the path using Unixfs resolver
	ResolvePath(context.Context, Path) (Path, error)

	// ResolveNode resolves the path (if not resolved already) using Unixfs
	// resolver, gets and returns the resolved Node
	ResolveNode(context.Context, Path) (Node, error)
}

// UnixfsAPI is the basic interface to immutable files in IPFS
type UnixfsAPI interface {
	// Add imports the data from the reader into merkledag file
	Add(context.Context, io.Reader) (Path, error)

	// Cat returns a reader for the file
	Cat(context.Context, Path) (Reader, error)

	// Ls returns the list of links in a directory
	Ls(context.Context, Path) ([]*Link, error)
}

// NameAPI specifies the interface to IPNS.
//
// IPNS is a PKI namespace, where names are the hashes of public keys, and the
// private key enables publishing new (signed) values. In both publish and
// resolve, the default name used is the node's own PeerID, which is the hash of
// its public key.
//
// You can use .Key API to list and generate more names and their respective keys.
type NameAPI interface {
	// Publish announces new IPNS name
	Publish(ctx context.Context, path Path, opts ...options.NamePublishOption) (*IpnsEntry, error)

	// WithValidTime is an option for Publish which specifies for how long the
	// entry will remain valid. Default value is 24h
	WithValidTime(validTime time.Duration) options.NamePublishOption

	// WithKey is an option for Publish which specifies the key to use for
	// publishing. Default value is "self" which is the node's own PeerID.
	//
	// You can use .Key API to list and generate more names and their respective keys.
	WithKey(key string) options.NamePublishOption

	// Resolve attempts to resolve the newest version of the specified name
	Resolve(ctx context.Context, name string, opts ...options.NameResolveOption) (Path, error)

	// WithRecursive is an option for Resolve which specifies whether to perform a
	// recursive lookup. Default value is false
	WithRecursive(recursive bool) options.NameResolveOption

	// WithLocal is an option for Resolve which specifies if the lookup should be
	// offline. Default value is false
	WithLocal(local bool) options.NameResolveOption

	// WithNoCache is an option for Resolve which specifies when set to true
	// disables the use of local name cache. Default value is false
	WithNoCache(nocache bool) options.NameResolveOption
}

// KeyAPI specifies the interface to Keystore
type KeyAPI interface {
	// Generate generates new key, stores it in the keystore under the specified
	// name and returns a base58 encoded multihash of it's public key
	Generate(ctx context.Context, name string, opts ...options.KeyGenerateOption) (string, error)

	// WithAlgorithm is an option for Generate which specifies which algorithm
	// should be used for the key. Default is "rsa"
	//
	// Supported algorithms:
	// * rsa
	// * ed25519
	WithAlgorithm(algorithm string) options.KeyGenerateOption

	// WithSize is an option for Generate which specifies the size of the key to
	// generated. Default is 0
	WithSize(size int) options.KeyGenerateOption

	// Rename renames oldName key to newName.
	Rename(ctx context.Context, oldName string, newName string, opts ...options.KeyRenameOption) (string, bool, error)

	// WithForce is an option for Rename which specifies whether to allow to
	// replace existing keys.
	WithForce(force bool) options.KeyRenameOption

	// List lists keys stored in keystore
	List(ctx context.Context) (map[string]string, error) //TODO: better key type?

	// Remove removes keys from keystore
	Remove(ctx context.Context, name string) (string, error)
}

// type ObjectAPI interface {
// 	New() (cid.Cid, Object)
// 	Get(string) (Object, error)
// 	Links(string) ([]*Link, error)
// 	Data(string) (Reader, error)
// 	Stat(string) (ObjectStat, error)
// 	Put(Object) (cid.Cid, error)
// 	SetData(string, Reader) (cid.Cid, error)
// 	AppendData(string, Data) (cid.Cid, error)
// 	AddLink(string, string, string) (cid.Cid, error)
// 	RmLink(string, string) (cid.Cid, error)
// }

// type ObjectStat struct {
// 	Cid            cid.Cid
// 	NumLinks       int
// 	BlockSize      int
// 	LinksSize      int
// 	DataSize       int
// 	CumulativeSize int
// }

var ErrIsDir = errors.New("object is a directory")
var ErrOffline = errors.New("can't resolve, ipfs node is offline")
