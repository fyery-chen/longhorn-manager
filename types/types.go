package types

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rancher/longhorn-manager/util"
)

const (
	DefaultAPIPort = 9500

	DefaultEngineBinaryPath          = "/usr/local/bin/longhorn"
	EngineBinaryDirectoryInContainer = "/engine-binaries/"
	EngineBinaryDirectoryOnHost      = "/var/lib/rancher/longhorn/engine-binaries/"
)

type ReplicaMode string

const (
	ReplicaModeRW  = ReplicaMode("RW")
	ReplicaModeWO  = ReplicaMode("WO")
	ReplicaModeERR = ReplicaMode("ERR")

	EnvNodeName     = "NODE_NAME"
	EnvPodNamespace = "POD_NAMESPACE"
	EnvPodIP        = "POD_IP"

	AWSAccessKey = "AWS_ACCESS_KEY_ID"
	AWSSecretKey = "AWS_SECRET_ACCESS_KEY"
	AWSEndPoint  = "AWS_ENDPOINTS"

	OptionFromBackup          = "fromBackup"
	OptionNumberOfReplica     = "numberOfReplicas"
	OptionStaleReplicaTimeout = "staleReplicaTimeout"
	OptionFrontend            = "frontend"

	EngineImageChecksumNameLength = 8
)

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v was not found", e.Name)
}

const (
	engineSuffix    = "-e"
	replicaSuffix   = "-r"
	recurringSuffix = "-c"

	// MaximumJobNameSize is calculated using
	// 1. NameMaximumLength is 40
	// 2. Recurring suffix is 2
	// 3. Maximum kubernetes name length is 63
	// 4. cronjob pod suffix is 11
	// 5. Dash and buffer for 2
	MaximumJobNameSize = 8

	engineImagePrefix = "ei-"
)

func GetEngineNameForVolume(vName string) string {
	return vName + engineSuffix
}

func GenerateReplicaNameForVolume(vName string) string {
	return vName + replicaSuffix + "-" + util.RandomID()
}

func GetCronJobNameForVolumeAndJob(vName, job string) string {
	return vName + "-" + job + recurringSuffix
}

func GetAPIServerAddressFromIP(ip string) string {
	return ip + ":" + strconv.Itoa(DefaultAPIPort)
}

func GetImageCanonicalName(image string) string {
	return strings.Replace(strings.Replace(image, ":", "-", -1), "/", "-", -1)
}

func GetEngineBinaryDirectoryOnHostForImage(image string) string {
	cname := GetImageCanonicalName(image)
	return filepath.Join(EngineBinaryDirectoryOnHost, cname)
}

func GetEngineBinaryDirectoryInContainerForImage(image string) string {
	cname := GetImageCanonicalName(image)
	return filepath.Join(EngineBinaryDirectoryInContainer, cname)
}

func EngineBinaryExistOnHostForImage(image string) bool {
	st, err := os.Stat(filepath.Join(GetEngineBinaryDirectoryOnHostForImage(image), "longhorn"))
	return err == nil && !st.IsDir()
}

var (
	LonghornSystemKey              = "longhorn"
	LonghornSystemValueManager     = "manager"
	LonghornSystemValueEngineImage = "engine-image"
)

func GetEngineImageLabel() map[string]string {
	return map[string]string{
		LonghornSystemKey: LonghornSystemValueEngineImage,
	}
}

func GetEngineImageChecksumName(image string) string {
	return engineImagePrefix + util.GetStringChecksum(strings.TrimSpace(image))[:EngineImageChecksumNameLength]
}
