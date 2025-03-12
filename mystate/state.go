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
	mu   sync.Mutex
	path string

	watcher        *fsnotify.Watcher
	watchersByID   map[resourceUrl]chan<- state.Event
	watchersByType map[resource.Type]chan<- state.Event
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
	// os.RemoveAll(statePath)
	os.Mkdir(statePath, os.ModePerm)
	return &State{
		mu:             sync.Mutex{},
		path:           statePath,
		watcher:        watcher,
		watchersByID:   make(map[resourceUrl]chan<- state.Event),
		watchersByType: make(map[resource.Type]chan<- state.Event),
	}
}

type resourceUrl string

func (url resourceUrl) parts() (resourceType resource.Type, id resource.ID) {
	parts := strings.Split(string(url), "/")
	return resource.Type(parts[0]), resource.ID(parts[1])
}

func getResourceUrl(resourceType resource.Type, id resource.ID) resourceUrl {
	return resourceUrl(resourceType + "/" + id)
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
	path := st.getResourcePath(resourcePointer)
	data, err := os.ReadFile(path)
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

func (st *State) CloseFileWatcher() {
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
	r, err := st.Get(ctx, resourcePointer)
	if err != nil {
		return fmt.Errorf("failed to get item to destroy: %w", err)
	}
	finalizers := r.Metadata().Finalizers()
	if !finalizers.Empty() {
		return fmt.Errorf("failed delete resource %s#%s, has pending finalizers", resourcePointer.Type(), resourcePointer.ID())
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	return os.Remove(st.getResourcePath(resourcePointer))
}

// Watch state of a resource by type.
//
// It's fine to watch for a resource which doesn't exist yet.
// Watch is canceled when context gets canceled.
// Watch sends initial resource state as the very first event on the channel,
// and then sends any updates to the resource as events.
func (st *State) Watch(ctx context.Context, resourcePointer resource.Pointer, ch chan<- state.Event, osp ...state.WatchOption) error {
	id := resourcePointer.ID()
	fmt.Println("Watch called for "+resourcePointer.Type(), "#", id)
	return nil
}

// WatchKind watches resources of specific kind (namespace and type).
func (st *State) WatchKind(ctx context.Context, resourceKind resource.Kind, ch chan<- state.Event, ops ...state.WatchKindOption) error {
	path := st.getResourceKindDirPath(resourceKind)
	os.Mkdir(path, os.ModePerm)
	err := st.watcher.Add(path)
	if err != nil {
		return err
	}
	st.watchersByType[resourceKind.Type()] = ch
	go func() {
		<-ctx.Done()
		fmt.Println("removing watcher for ", path)
		st.watcher.Remove(path)
		delete(st.watchersByType, resourceKind.Type())
	}()
	return nil
}

func (st *State) StartFileWatcher(ctx context.Context) {
	for {
		select {
		case err := <-st.watcher.Errors:
			log.Println("watch error:", err)
		case <-ctx.Done():
			return
		case event := <-st.watcher.Events:
			if event.Has(fsnotify.Chmod) {
				continue
			}

			resourceType, id := getResourceDataFromEvent(event)
			listeners := st.findListeners(resourceType, id)
			if len(listeners) == 0 {
				continue
			}

			e := state.Event{Resource: myresource.Cast(resource.NewMetadata(constants.NS, resourceType, id, resource.VersionUndefined), nil)}
			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				e.Type = state.Destroyed
				sendToListeners(listeners, e)
				continue
			}
			r, err := st.Get(ctx, resource.NewMetadata(constants.NS, resourceType, id, resource.VersionUndefined))
			if state.IsNotFoundError(err) {
				continue
			}
			e.Error = err
			e.Resource = r
			if event.Has(fsnotify.Create) {
				e.Type = state.Created
				sendToListeners(listeners, e)
				continue
			}
			if event.Has(fsnotify.Write) {
				e.Type = state.Updated
				sendToListeners(listeners, e)
				continue
			}
		}
	}
}

func sendToListeners(listeners []chan<- state.Event, e state.Event) {
	for _, listener := range listeners {
		listener <- e
	}
}

func getResourceDataFromEvent(event fsnotify.Event) (string, string) {
	pathParts := strings.Split(event.Name, "/")
	filename, _ := last(pathParts)
	resourceType := pathParts[len(pathParts)-2]
	id := strings.Split(filename, ".")[0]
	return resourceType, id
}

func (st *State) findListeners(resourceType string, id string) []chan<- state.Event {
	var listeners []chan<- state.Event
	for kind, ch := range st.watchersByType {
		if kind == resourceType {
			listeners = append(listeners, ch)
		}
	}
	for resourcePointer, ch := range st.watchersByID {
		rType, ID := resourcePointer.parts()
		if rType == resourceType && ID == id {
			listeners = append(listeners, ch)
		}
	}
	return listeners
}

// WatchKindAggregated watches resources of specific kind (namespace and type), updates are sent aggregated.
func (st *State) WatchKindAggregated(ctx context.Context, resourceKind resource.Kind, ch chan<- []state.Event, ops ...state.WatchKindOption) error {
	var options state.WatchKindOptions
	for _, opt := range ops {
		opt(&options)
	}
	if options.BootstrapContents {
		resources, err := st.List(ctx, resourceKind)
		if err != nil {
			return err
		}
		var events []state.Event
		for _, r := range resources.Items {
			events = append(events, state.Event{Resource: r, Type: state.Created})
		}
		events = append(events, state.Event{Type: state.Bootstrapped})
		ch <- events
	}
	singleElChan := make(chan state.Event)
	go func() {
		for {
			select {
			case event := <-singleElChan:
				fmt.Printf("Event: %s %s#%s\n", event.Type, event.Resource.Metadata().Type(), event.Resource.Metadata().ID())
				ch <- []state.Event{event}
			case <-ctx.Done():
				return
			}
		}
	}()
	return st.WatchKind(ctx, resourceKind, singleElChan, nil)
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
