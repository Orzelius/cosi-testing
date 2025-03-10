package mystate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Orzelius/cosi-testing/constants"
	"github.com/Orzelius/cosi-testing/myresource"
	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/cosi-project/runtime/pkg/state"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// State implements state.CoreState.
type State struct {
	mu      sync.Mutex
	path    string
	watcher *fsnotify.Watcher
}

func NewState() *State {

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	statePath := path + "/test"
	os.RemoveAll(statePath)
	os.Mkdir(statePath, os.ModePerm)
	return &State{
		mu:      sync.Mutex{},
		path:    statePath,
		watcher: watcher,
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
	marshalable := myresource.Marshalable{Val: r.Spec(), Md: r.Metadata()}
	data, err := yaml.Marshal(&marshalable)
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
	st.mu.Lock()
	defer st.mu.Unlock()
	data, err := os.ReadFile(st.getResourcePath(resourcePointer))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrNotFound(resourcePointer)
		}

		return nil, err
	}

	var r = myresource.UnMarshalable{}
	err = yaml.Unmarshal(data, &r)
	return myresource.Cast(r.Md, r.Val.Val), err
}

// List resources by type.
func (st *State) List(ctx context.Context, resourceKind resource.Kind, ops ...state.ListOption) (resource.List, error) {
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

		var r = myresource.UnMarshalable{}
		err = yaml.Unmarshal(data, &r)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal data: %w", err)
		}
		result.Items = append(result.Items, myresource.Cast(r.Md, r.Val.Val))
	}

	return result, nil
}

func (st *State) Close() {
	err := st.watcher.Close()
	if err != nil {
		panic(err)
	}
}

// Update a resource.
//
// If a resource doesn't exist, error is returned.
// On update current version of resource `new` in the state should match
// the version on the backend, otherwise conflict error is returned.
func (st *State) Update(ctx context.Context, newResource resource.Resource, ops ...state.UpdateOption) error {
	var options state.UpdateOptions
	for _, opt := range ops {
		opt(&options)
	}

	st.mu.Lock()
	defer st.mu.Unlock()
	f, err := os.OpenFile(st.getResourcePath(newResource.Metadata()), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ErrNotFound(newResource.Metadata())
		}

		return err
	}
	defer f.Close()
	marshalable := myresource.Marshalable{Val: newResource.Spec(), Md: newResource.Metadata()}
	data, err := yaml.Marshal(&marshalable)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
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
	id := resourcePointer.ID()
	return st.watch(ctx, resourcePointer, ch, &id)
}

// WatchKind watches resources of specific kind (namespace and type).
func (st *State) WatchKind(ctx context.Context, resourceKind resource.Kind, ch chan<- state.Event, ops ...state.WatchKindOption) error {
	return st.watch(ctx, resourceKind, ch, nil)
}

// WatchKind watches resources of specific kind (namespace and type).
func (st *State) watch(ctx context.Context, resourceKind resource.Kind, ch chan<- state.Event, watchID *resource.ID, ops ...state.WatchKindOption) error {
	fmt.Printf("watch called for type: %s\n", resourceKind.Type())
	go func() {
		for {
			select {
			case event := <-st.watcher.Events:
				log.Println("event:", event)
				filename, _ := last(strings.Split(event.Name, "/"))
				id := strings.Split(filename, ".")[0]
				if watchID != nil && *watchID != id {
					continue
				}
				e := state.Event{Resource: myresource.Cast(resource.NewMetadata(constants.NS, resourceKind.Type(), id, resource.VersionUndefined), nil)}
				if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					e.Type = state.Destroyed
					ch <- e
					continue
				}
				r, err := st.Get(ctx, resource.NewMetadata(constants.NS, resourceKind.Type(), id, resource.VersionUndefined))
				if state.IsNotFoundError(err) {
					continue
				}
				e.Error = err
				e.Resource = r
				if event.Has(fsnotify.Create) {
					e.Type = state.Created
					ch <- e
					continue
				}
				if event.Has(fsnotify.Write) {
					e.Type = state.Updated
					ch <- e
					continue
				}
			case err := <-st.watcher.Errors:
				log.Println("watch error:", err)
			case <-ctx.Done():
				return
			}
		}
	}()
	path := st.getResourceKindDirPath(resourceKind)
	os.Mkdir(path, os.ModePerm)
	err := st.watcher.Add(path)
	if err != nil {
		return err
	}
	return nil
}

// WatchKindAggregated watches resources of specific kind (namespace and type), updates are sent aggregated.
func (st *State) WatchKindAggregated(ctx context.Context, resourceKind resource.Kind, ch chan<- []state.Event, ops ...state.WatchKindOption) error {
	singleElChan := make(chan state.Event)
	go func() {
		for {
			select {
			case event := <-singleElChan:
				ch <- []state.Event{event}
			case <-ctx.Done():
				return
			}
		}
	}()
	return st.watch(ctx, resourceKind, singleElChan, nil)
}

func (st *State) getResourcePath(r resource.Pointer) string {
	return st.getResourceKindDirPath(r) + "/" + r.ID() + ".yaml"
}

func (st *State) getResourceKindDirPath(r resource.Kind) string {
	return st.path + "/" + r.Type()
}

func last[E any](s []E) (E, bool) {
	if len(s) == 0 {
		var zero E
		return zero, false
	}
	return s[len(s)-1], true
}
