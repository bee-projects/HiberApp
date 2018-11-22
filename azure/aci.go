package azure

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
	"github.com/writeameer/aci/helpers"
)

// ContainerSpec defines the details of the container to launch
type ContainerSpec struct {
	ContainerName        string
	Ports                []int32
	ContainerImage       string
	CPU                  float64
	MemoryInGB           float64
	EnvironmentVariables map[string]string
	//VolumeMount          AzureFileMount
}

// ContainerGroupSpec defines the details of the container to launch
type ContainerGroupSpec struct {
	ResourceGroupName string
	Name              string
	Ports             []int32
	DNSNameLabel      string
	OsType            containerinstance.OperatingSystemTypes
	RestartPolicy     containerinstance.ContainerGroupRestartPolicy
	IPAddressType     containerinstance.ContainerGroupIPAddressType
}

// AzureFileMount describes the Azure File Mount for a container
// type AzureFileMount struct {
// 	ShareName          string
// 	StorageAccountKey  string
// 	StorageAccountName string
// }

// GetContainerFromSpec returns a container struct with provided config
func GetContainerFromSpec(containerSpec ContainerSpec) (container containerinstance.Container) {

	// Define Env Variables
	var envVars []containerinstance.EnvironmentVariable
	for key, value := range containerSpec.EnvironmentVariables {
		k := key
		v := value
		envVars = append(envVars, containerinstance.EnvironmentVariable{
			Name:  &k,
			Value: &v,
		})
	}

	// Define container's properties
	containerProperties := containerinstance.ContainerProperties{
		Image: &containerSpec.ContainerImage,
		Ports: setTCPPort(containerSpec.Ports),
		Resources: &containerinstance.ResourceRequirements{
			Requests: setResourceRequests(containerSpec.CPU, containerSpec.MemoryInGB),
		},
		EnvironmentVariables: &envVars,
	}

	// Define a container with given properties
	container = containerinstance.Container{
		ContainerProperties: &containerProperties,
		Name:                &containerSpec.ContainerName,
	}

	// return containers
	return
}

// GetContainersFromSpec returns an array of Container structs from the specs provided
func GetContainersFromSpec(containerSpecs []ContainerSpec) (containers *[]containerinstance.Container) {
	var myContainerSpecs []containerinstance.Container
	for _, containerSpec := range containerSpecs {
		myContainerSpecs = append(myContainerSpecs, GetContainerFromSpec(containerSpec))
	}

	return &myContainerSpecs
}

// GetContainerGroupFromSpec converts a ContainerGroupSpec and ContainerSpec struct to a containerinstance.ContainerGroupProperties struct
func GetContainerGroupFromSpec(containerGroupSpec ContainerGroupSpec, containerSpecs []ContainerSpec) (containerGroup *containerinstance.ContainerGroupProperties) {
	log.Println("Starting GetContainerGroupFromSpec...")

	cgroup := containerinstance.ContainerGroupProperties{
		Containers:    GetContainersFromSpec(containerSpecs),
		OsType:        containerGroupSpec.OsType,
		RestartPolicy: containerGroupSpec.RestartPolicy,
		IPAddress: &containerinstance.IPAddress{
			Type:         containerinstance.Public,
			DNSNameLabel: &containerGroupSpec.DNSNameLabel,
			Ports:        setContainerGroupTCPPort(containerGroupSpec.Ports),
		},
	}

	return &cgroup
}

func setTCPPort(ports []int32) (containerPorts *[]containerinstance.ContainerPort) {
	log.Println("Starting setTCPPort...")
	var portList []containerinstance.ContainerPort

	for _, port := range ports {
		//log.Printf("%d, %d", i, port)
		portList = append(portList, containerinstance.ContainerPort{
			Port:     &port,
			Protocol: "tcp",
		})
	}

	return &portList
}

func setContainerGroupTCPPort(ports []int32) (containerPorts *[]containerinstance.Port) {

	var portList []containerinstance.Port

	for _, port := range ports {
		//log.Printf("%d, %d", i, port)
		portList = append(portList, containerinstance.Port{
			Port:     &port,
			Protocol: "tcp",
		})
	}

	return &portList
}

func setResourceRequests(cpu float64, memoryInGB float64) (resourceRequirements *containerinstance.ResourceRequests) {
	requirements := containerinstance.ResourceRequests{
		CPU:        &cpu,
		MemoryInGB: &memoryInGB,
	}

	return &requirements
}

// DeployContainer Deploys a container to ACI
func DeployContainer(containerLocation string, resourceGroupName string, containerGroupName string, containersSpec []ContainerSpec, containerGroupSpec ContainerGroupSpec) (deployedGroup containerinstance.ContainerGroup, err error) {

	// Define container group
	containerGroupProperties := GetContainerGroupFromSpec(containerGroupSpec, containersSpec)
	log.Println("Created containerGroupProperties")

	cgroup := containerinstance.ContainerGroup{
		ContainerGroupProperties: containerGroupProperties,
		Location:                 &containerLocation,
		Name:                     &containerGroupName,
	}

	log.Println("Created containnerGroup")

	// Authenticate with Azure
	authorizer, sid := Auth()

	// Get container service client and create container group
	client := containerinstance.NewContainerGroupsClient(sid)
	client.Authorizer = authorizer

	deploymentFuture, err := client.CreateOrUpdate(ctx, resourceGroupName, containerGroupName, cgroup)
	helpers.PrintError(err)

	err = deploymentFuture.Future.WaitForCompletion(ctx, client.BaseClient.Client)
	helpers.PrintError(err)

	log.Printf("Deployment completed...")

	deployedGroup, err = deploymentFuture.Result(client)
	helpers.PrintError(err)

	return
}
