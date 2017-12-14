package iface

import (
	"context"
	"errors"
	"io"

	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"
	ipld "gx/ipfs/QmPN7cwmpcc4DWXb4KTB9dNAJgjuPY69h3npsMfhRrQL9c/go-ipld-format"
)

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

type Reader interface {
	io.ReadSeeker
	io.Closer
}

type CoreAPI interface {
	Unixfs() UnixfsAPI
	ResolvePath(context.Context, Path) (Path, error)
	ResolveNode(context.Context, Path) (Node, error)
}

type UnixfsAPI interface {
	Add(context.Context, io.Reader) (Path, error)
	Cat(context.Context, Path) (Reader, error)
	Ls(context.Context, Path) ([]*Link, error)
}

//TODO: Should this use paths instead of cids?
type ObjectAPI interface {
	New(ctx context.Context) (Node, error)
	Put(context.Context, Node) error
	Get(context.Context, Path) (Node, error)
	Data(context.Context, Path) (io.Reader, error)
	Links(context.Context, Path) ([]*Link, error)
	Stat(context.Context, Path) (*ObjectStat, error)

	AddLink(ctx context.Context, base Path, name string, child Path, create bool) (Node, error) //TODO: make create optional
	RmLink(context.Context, Path, string) (Node, error)
	AppendData(context.Context, Path, io.Reader) (Node, error)
	SetData(context.Context, Path, io.Reader) (Node, error)
}

type ObjectStat struct {
	Cid            *cid.Cid
	NumLinks       int
	BlockSize      int
	LinksSize      int
	DataSize       int
	CumulativeSize int
}

var ErrIsDir = errors.New("object is a directory")
var ErrOffline = errors.New("can't resolve, ipfs node is offline")
