package stepman

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

const (
	// StepmanDirname ...
	StepmanDirname string = ".stepman"
	// RoutingFilename ...
	RoutingFilename string = "routing.json"
	// CollectionsDirname ...
	CollectionsDirname string = "step_collections"
)

var (
	// StepManDirPath ...
	StepManDirPath string
	// CollectionsDirPath ...
	CollectionsDirPath string

	routingFilePath string
)

// SteplibRoute ...
type SteplibRoute struct {
	SteplibURI  string
	FolderAlias string
}

// SteplibRoutes ...
type SteplibRoutes []SteplibRoute

// GetRoute ...
func (routes SteplibRoutes) GetRoute(URI string) (route SteplibRoute, found bool) {
	for _, route := range routes {
		if route.SteplibURI == URI {
			return route, true
		}
	}
	return SteplibRoute{}, false
}

// ReadRoute ...
func ReadRoute(uri string) (route SteplibRoute, found bool) {
	routes, err := readRouteMap()
	if err != nil {
		return SteplibRoute{}, false
	}

	return routes.GetRoute(uri)
}

func (routes SteplibRoutes) writeToFile() error {
	routeMap := map[string]string{}
	for _, route := range routes {
		routeMap[route.SteplibURI] = route.FolderAlias
	}
	bytes, err := json.MarshalIndent(routeMap, "", "\t")
	if err != nil {
		return err
	}
	return fileutil.WriteBytesToFile(routingFilePath, bytes)
}

// CleanupRoute ...
func CleanupRoute(route SteplibRoute) error {
	pth := CollectionsDirPath + "/" + route.FolderAlias
	if err := cmdex.RemoveDir(pth); err != nil {
		return err
	}
	if err := RemoveRoute(route); err != nil {
		return err
	}
	return nil
}

// RootExistForCollection ...
func RootExistForCollection(collectionURI string) (bool, error) {
	routes, err := readRouteMap()
	if err != nil {
		return false, err
	}

	_, found := routes.GetRoute(collectionURI)
	return found, nil
}

func getAlias(uri string) (string, error) {
	routes, err := readRouteMap()
	if err != nil {
		return "", err
	}

	route, found := routes.GetRoute(uri)
	if found == false {
		return "", errors.New("No routes exist for uri:" + uri)
	}
	return route.FolderAlias, nil
}

// RemoveRoute ...
func RemoveRoute(route SteplibRoute) error {
	routes, err := readRouteMap()
	if err != nil {
		return err
	}

	newRoutes := SteplibRoutes{}
	for _, aRoute := range routes {
		if aRoute.SteplibURI != route.SteplibURI {
			newRoutes = append(newRoutes, aRoute)
		}
	}
	if err := newRoutes.writeToFile(); err != nil {
		return err
	}
	return nil
}

// AddRoute ...
func AddRoute(route SteplibRoute) error {
	routes, err := readRouteMap()
	if err != nil {
		return err
	}

	routes = append(routes, route)
	if err := routes.writeToFile(); err != nil {
		return err
	}

	return nil
}

// GenerateFolderAlias ...
func GenerateFolderAlias() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

func readRouteMap() (SteplibRoutes, error) {
	exist, err := pathutil.IsPathExists(routingFilePath)
	if err != nil {
		return SteplibRoutes{}, err
	} else if !exist {
		return SteplibRoutes{}, nil
	}

	bytes, err := fileutil.ReadBytesFromFile(routingFilePath)
	if err != nil {
		return SteplibRoutes{}, err
	}
	var routeMap map[string]string
	if err := json.Unmarshal(bytes, &routeMap); err != nil {
		return SteplibRoutes{}, err
	}

	routes := []SteplibRoute{}
	for key, value := range routeMap {
		routes = append(routes, SteplibRoute{
			SteplibURI:  key,
			FolderAlias: value,
		})
	}

	return routes, nil
}

// CreateStepManDirIfNeeded ...
func CreateStepManDirIfNeeded() error {
	return os.MkdirAll(StepManDirPath, 0777)
}

// GetStepSpecPath ...
func GetStepSpecPath(route SteplibRoute) string {
	return CollectionsDirPath + "/" + route.FolderAlias + "/spec/spec.json"
}

// GetCacheBaseDir ...
func GetCacheBaseDir(route SteplibRoute) string {
	return CollectionsDirPath + "/" + route.FolderAlias + "/cache"
}

// GetCollectionBaseDirPath ...
func GetCollectionBaseDirPath(route SteplibRoute) string {
	return CollectionsDirPath + "/" + route.FolderAlias + "/collection"
}

// GetAllStepCollectionPath ...
func GetAllStepCollectionPath() []string {
	routes, err := readRouteMap()
	if err != nil {
		log.Error("[STEPMAN] - Failed to read step specs path:", err)
		return []string{}
	}

	sources := []string{}
	for _, route := range routes {
		sources = append(sources, route.SteplibURI)
	}

	return sources
}

// GetStepCacheDirPath ...
// Step's Cache dir path, where it's code lives.
func GetStepCacheDirPath(route SteplibRoute, id, version string) string {
	return GetCacheBaseDir(route) + "/" + id + "/" + version
}

// GetStepCollectionDirPath ...
// Step's Collection dir path, where it's spec (step.yml) lives.
func GetStepCollectionDirPath(route SteplibRoute, id, version string) string {
	return GetCollectionBaseDirPath(route) + "/steps/" + id + "/" + version
}

// Life cycle
func init() {
	StepManDirPath = pathutil.UserHomeDir() + "/" + StepmanDirname
	routingFilePath = StepManDirPath + "/" + RoutingFilename
	CollectionsDirPath = StepManDirPath + "/" + CollectionsDirname
}
