package mystate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"gopkg.in/yaml.v3"
)

// State implements state.CoreState.
type State struct {
	mu   sync.Mutex
	path string
}

func NewState() *State {

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	statePath := path + "/test"
	// os.RemoveAll(statePath)
	return &State{
		mu:   sync.Mutex{},
		path: statePath,
	}
}

// Create a resource.
//
// If a resource already exists, Create returns an error.
func (st *State) Create(ctx context.Context, r resource.Resource, ops ...state.CreateOption) error {
	var options state.CreateOptions
	for _, opt := range ops {
		opt(&options)
	}
	resCopy := r.DeepCopy()
	if err := resCopy.Metadata().SetOwner(options.Owner); err != nil {
		return err
	}
	st.mu.Lock()
	defer st.mu.Unlock()

	if err := os.MkdirAll(st.path+"/"+resCopy.Metadata().Type(), 0777); err != nil {
		return err
	}
	f, err := os.OpenFile(st.getResourcePath(r.Metadata()), os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	marshaled, _ := resource.MarshalYAML(r)
	data, err := yaml.Marshal(marshaled)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

// Get a resource by type and ID.
//
// If a resource is not found, error is returned.
func (st *State) Get(ctx context.Context, resourcePointer resource.Pointer, ops ...state.GetOption) (resource.Resource, error) {
	data, err := os.ReadFile(st.getResourcePath(resourcePointer))
	if err != nil {
		return nil, err
	}

	var r resource.Resource
	err = yaml.Unmarshal(data, r)
	return r, err
}

// List resources by type.
func (st *State) List(ctx context.Context, resourceKind resource.Kind, ops ...state.ListOption) (resource.List, error) {
	var options state.ListOptions
	for _, opt := range ops {
		opt(&options)
	}
	st.mu.Lock()
	defer st.mu.Unlock()

	result := resource.List{
		Items: []resource.Resource{},
	}

	dir := st.getResourceKindDirPath(resourceKind)
	files, err := os.ReadDir(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return result, nil
		}

		return result, err
	}

	for _, f := range files {
		data, err := os.ReadFile(dir + "/" + f.Name())
		if err != nil {
			return result, err
		}

		var r resource.Any = resource.Any{}
		err = yaml.Unmarshal(data, r)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		result.Items = append(result.Items, &r)
	}

	return result, nil
}

// Update a resource.
//
// If a resource doesn't exist, error is returned.
// On update current version of resource `new` in the state should match
// the version on the backend, otherwise conflict error is returned.
func (st *State) Update(ctx context.Context, newResource resource.Resource, opts ...state.UpdateOption) error {
	fmt.Println("state Update() call, unimplemented")
	return nil
}

// Destroy a resource.
//
// If a resource doesn't exist, error is returned.
// If a resource has pending finalizers, error is returned.
func (st *State) Destroy(ctx context.Context, resourcePointer resource.Pointer, ops ...state.DestroyOption) error {
	fmt.Println("state Destroy() call, unimplemented")
	return nil
}

// Watch state of a resource by type.
//
// It's fine to watch for a resource which doesn't exist yet.
// Watch is canceled when context gets canceled.
// Watch sends initial resource state as the very first event on the channel,
// and then sends any updates to the resource as events.
func (st *State) Watch(ctx context.Context, resourcePointer resource.Pointer, ch chan<- state.Event, osp ...state.WatchOption) error {
	fmt.Println("state Watch() call, unimplemented")
	return nil
}

// WatchKind watches resources of specific kind (namespace and type).
func (st *State) WatchKind(ctx context.Context, resourceKind resource.Kind, ch chan<- state.Event, ops ...state.WatchKindOption) error {
	fmt.Println("state WatchKind() call, unimplemented")
	return nil
}

// WatchKindAggregated watches resources of specific kind (namespace and type), updates are sent aggregated.
func (st *State) WatchKindAggregated(ctx context.Context, resourceKind resource.Kind, ch chan<- []state.Event, ops ...state.WatchKindOption) error {
	fmt.Println("state WatchKindAggregated() call, unimplemented")
	return nil
}

func (st *State) getResourcePath(r resource.Pointer) string {
	return st.getResourceKindDirPath(r) + "/" + r.ID() + ".yaml"
}

func (st *State) getResourceKindDirPath(r resource.Kind) string {
	return st.path + "/" + r.Type()
}
