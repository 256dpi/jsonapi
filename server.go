package jsonapi

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// ServerConfig is used to configure a server.
type ServerConfig struct {
	Prefix string
	Types  []string
}

// Server implements a basic in-memory jsonapi resource server intended for
// testing purposes.
type Server struct {
	Config  ServerConfig
	Parser  *Parser
	Data    map[string]map[string]*Resource
	Counter int
	Mutex   sync.Mutex
}

// NewServer will create and return a new server.
func NewServer(config ServerConfig) *Server {
	// clean prefix
	config.Prefix = "/" + strings.Trim(config.Prefix, "/")

	// prepare parser
	parser := &Parser{
		Prefix: config.Prefix,
	}

	return &Server{
		Config: config,
		Data:   map[string]map[string]*Resource{},
		Parser: parser,
	}
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// acquire mutex
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	// parse request
	req, err := s.Parser.ParseRequest(r)
	if err != nil {
		_ = WriteError(w, err)
		return
	}

	// check resource type if list is given
	if len(s.Config.Types) > 0 {
		var ok bool
		for _, typ := range s.Config.Types {
			if req.ResourceType == typ {
				ok = true
			}
		}
		if !ok {
			_ = WriteError(w, BadRequest("unsupported resource type"))
			return
		}
	}

	// parse document
	var doc *Document
	if req.Intent.DocumentExpected() {
		doc, err = ParseDocument(r.Body)
		if err != nil {
			_ = WriteError(w, err)
			return
		}
	}

	// handle intent
	switch req.Intent {
	case ListResources:
		err = s.listResources(req, w)
	case FindResource:
		err = s.findResources(req, w)
	case CreateResource:
		err = s.createResource(req, doc, w)
	case UpdateResource:
		err = s.updateResource(req, doc, w)
	case DeleteResource:
		err = s.deleteResource(req, w)
	default:
		err = BadRequest("unsupported request method")
	}
	if err != nil {
		_ = WriteError(w, err)
	}
}

func (s *Server) listResources(req *Request, w http.ResponseWriter) error {
	// prepare list
	var list []*Resource

	// get resources from collection
	coll := s.Data[req.ResourceType]
	if coll != nil {
		list = make([]*Resource, 0, len(coll))
		for _, res := range coll {
			list = append(list, res)
		}
	} else {
		list = []*Resource{}
	}

	// sort list
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	// get offset and limit
	offset := int(req.PageOffset)
	limit := int(req.PageLimit)
	if offset == 0 && req.PageNumber > 0 {
		offset = int(req.PageNumber * req.PageSize)
		limit = offset + int(req.PageSize)
	}

	// check offset and limit
	if offset > 0 && (offset >= len(list) || limit > len(list)) {
		return BadRequest("invalid pagination parameters")
	}

	// apply pagination
	if offset > 0 {
		list = list[offset:limit]
	}

	return WriteResources(w, http.StatusOK, list, &DocumentLinks{
		Self: req.Self(),
	})
}

func (s *Server) findResources(req *Request, w http.ResponseWriter) error {
	// get collection
	coll := s.Data[req.ResourceType]
	if coll == nil {
		return NotFound("unknown resource")
	}

	// get resource
	res := coll[req.ResourceID]
	if res == nil {
		return NotFound("unknown resource")
	}

	return WriteResource(w, http.StatusOK, res, &DocumentLinks{
		Self: req.Self(),
	})
}

func (s *Server) createResource(req *Request, doc *Document, w http.ResponseWriter) error {
	// check document
	if doc.Data == nil || doc.Data.One == nil {
		return BadRequest("missing resource")
	}

	// get resource
	res := doc.Data.One

	// check type
	if res.Type != req.ResourceType {
		return BadRequest("resource type mismatch")
	}

	// check id
	if res.ID == "" {
		s.Counter++
		res.ID = "s-" + strconv.Itoa(s.Counter)
	}

	// link relationships
	s.linkRelationships(res)

	// get collection
	coll := s.Data[req.ResourceType]
	if coll == nil {
		coll = map[string]*Resource{}
		s.Data[req.ResourceType] = coll
	}

	// check existence
	if coll[res.ID] != nil {
		return BadRequest("conflicting resource")
	}

	// store resource
	coll[res.ID] = res

	// set id
	req.ResourceID = res.ID

	return WriteResource(w, http.StatusOK, res, &DocumentLinks{
		Self: req.Self(),
	})
}

func (s *Server) updateResource(req *Request, doc *Document, w http.ResponseWriter) error {
	// check document
	if doc.Data == nil || doc.Data.One == nil {
		return BadRequest("missing resource")
	}

	// get resource
	res := doc.Data.One

	// check type
	if res.Type != req.ResourceType {
		return BadRequest("resource type mismatch")
	}

	// check id
	if res.ID != req.ResourceID {
		return BadRequest("resource id mismatch")
	}

	// link relationships
	s.linkRelationships(res)

	// get collection
	coll := s.Data[req.ResourceType]
	if coll == nil {
		return NotFound("unknown resource")
	}

	// get resource
	if coll[req.ResourceID] == nil {
		return NotFound("unknown resource")
	}

	// update resource
	coll[req.ResourceID] = res

	return WriteResource(w, http.StatusOK, res, &DocumentLinks{
		Self: req.Self(),
	})
}

func (s *Server) deleteResource(req *Request, w http.ResponseWriter) error {
	// get collection
	coll := s.Data[req.ResourceType]
	if coll == nil {
		return NotFound("unknown resource")
	}

	// check resource
	if coll[req.ResourceID] == nil {
		return NotFound("unknown resource")
	}

	// delete resource
	delete(coll, req.ResourceID)

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (s *Server) linkRelationships(res *Resource) {
	for name, doc := range res.Relationships {
		doc.Links = &DocumentLinks{
			Self:    fmt.Sprintf("%s/%s/%s/relationships/%s", s.Config.Prefix, res.Type, res.ID, name),
			Related: fmt.Sprintf("%s/%s/%s/%s", s.Config.Prefix, res.Type, res.ID, name),
		}
	}
}
