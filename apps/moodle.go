package apps

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/containerinstance/mgmt/containerinstance"
	"github.com/writeameer/aci/azure"
	"github.com/writeameer/aci/helpers"
)

// RunMoodle Runs wordpress on ACI
func RunMoodle(resourceGroupName string, containerGroupName string) {

	//	wp_config_extra := "\n\tdefine('WP_HOME', 'http://localhost:8000/');\n\tdefine('WP_SITEURL', 'http://localhost:8000/');\n"

	// Define Images
	wordpressImage := "wordpress:4.9-apache"
	mysqlImage := "mysql:5.6"

	// Define containers to run
	containerSpecs := []azure.ContainerSpec{
		azure.ContainerSpec{
			ContainerName:  "wordpress",
			ContainerImage: wordpressImage,
			Ports:          []int32{80},
			CPU:            0.5,
			MemoryInGB:     0.5,
			EnvironmentVariables: map[string]string{
				"WORDPRESS_DB_HOST":     "127.0.0.1:3306",
				"WORDPRESS_DB_PASSWORD": "0rsmP@ssw0rd",
				//"WORDPRESS_CONFIG_EXTRA": wp_config_extra,
			},
		},
		azure.ContainerSpec{
			ContainerName:  "bitnami/mariadb:latest",
			ContainerImage: mysqlImage,
			Ports:          []int32{3306},
			CPU:            0.5,
			MemoryInGB:     0.5,
			EnvironmentVariables: map[string]string{
				"MYSQL_ROOT_PASSWORD": "0rsmP@ssw0rd",
			},
		},
	}

	// Define the container group's specifications
	containerGroupSpecs := azure.ContainerGroupSpec{
		ResourceGroupName: resourceGroupName,
		Name:              containerGroupName,
		Ports:             []int32{80},
		DNSNameLabel:      "hiberapp",
		OsType:            containerinstance.Linux,
		IPAddressType:     containerinstance.Public,
	}

	// Get ARM group to inspect location
	armGroup, err := azure.GetGroup(resourceGroupName)
	helpers.PrintError(err)

	deployedGroup, err := azure.DeployContainer(*armGroup.Location, resourceGroupName, containerGroupName, containerSpecs, containerGroupSpecs)
	log.Printf(*deployedGroup.IPAddress.Fqdn)
	helpers.PrintError(err)
}
