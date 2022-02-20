// Package chainlink  handler is the representation layer for the  chainlink node
// uses the k8 client to CRUD the nodes
package chainlink

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kotalco/api/internal/chainlink"
	restError "github.com/kotalco/api/pkg/errors"
	"github.com/kotalco/api/pkg/logger"
	"github.com/kotalco/api/pkg/shared"
	chainlinkv1alpha1 "github.com/kotalco/kotal/apis/chainlink/v1alpha1"
	sharedAPI "github.com/kotalco/kotal/apis/shared"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sort"
	"strconv"
)

var service = chainlink.ChainlinkService

// Get returns a single chainlink node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(*chainlinkv1alpha1.Node)
	return c.JSON(shared.NewResponse(new(Node).FromChainlinkNode(node)))
}

// Create creates chainlink node from the given spec
func Create(c *fiber.Ctx) error {
	//Todo add Validation
	request := new(Node)
	if err := c.BodyParser(request); err != nil {
		go logger.Error("error parsing create chainlink node", err)
		badReqErr := restError.NewBadRequestError("invalid request body")
		return c.Status(badReqErr.Status).JSON(badReqErr)
	}

	node := &chainlinkv1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: "default",
		},
		Spec: chainlinkv1alpha1.NodeSpec{
			EthereumChainId:            request.EthereumChainId,
			LinkContractAddress:        request.LinkContractAddress,
			EthereumWSEndpoint:         request.EthereumWSEndpoint,
			DatabaseURL:                request.DatabaseURL,
			KeystorePasswordSecretName: request.KeystorePasswordSecretName,
			APICredentials: chainlinkv1alpha1.APICredentials{
				Email:              request.APICredentials.Email,
				PasswordSecretName: request.APICredentials.PasswordSecretName,
			},
		},
	}

	node, err := service.Create(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(Node).FromChainlinkNode(node)))
}

// Update updates a single chainlink node by name from spec
func Update(c *fiber.Ctx) error {
	request := new(Node)

	if err := c.BodyParser(request); err != nil {
		badReq := restError.NewBadRequestError("invalid request body")
		return c.Status(badReq.Status).JSON(err)
	}

	node := c.Locals("node").(*chainlinkv1alpha1.Node)

	if request.EthereumWSEndpoint != "" {
		node.Spec.EthereumWSEndpoint = request.EthereumWSEndpoint
	}

	if request.DatabaseURL != "" {
		node.Spec.DatabaseURL = request.DatabaseURL
	}

	if len(request.EthereumHTTPEndpoints) != 0 {
		node.Spec.EthereumHTTPEndpoints = request.EthereumHTTPEndpoints
	}

	if request.KeystorePasswordSecretName != "" {
		node.Spec.KeystorePasswordSecretName = request.KeystorePasswordSecretName
	}

	if request.APICredentials != nil {
		node.Spec.APICredentials.Email = request.APICredentials.Email
		node.Spec.APICredentials.PasswordSecretName = request.APICredentials.PasswordSecretName
	}

	if len(request.CORSDomains) != 0 {
		node.Spec.CORSDomains = request.CORSDomains
	}

	if request.CertSecretName != "" {
		node.Spec.CertSecretName = request.CertSecretName
	}

	if request.TLSPort != 0 {
		node.Spec.TLSPort = request.TLSPort
	}

	if request.P2PPort != 0 {
		node.Spec.P2PPort = request.P2PPort
	}

	if request.APIPort != 0 {
		node.Spec.APIPort = request.APIPort
	}

	if request.SecureCookies != nil {
		node.Spec.SecureCookies = *request.SecureCookies
	}

	if request.Logging != "" {
		node.Spec.Logging = sharedAPI.VerbosityLevel(request.Logging)
	}

	if request.CPU != "" {
		node.Spec.CPU = request.CPU
	}
	if request.CPULimit != "" {
		node.Spec.CPULimit = request.CPULimit
	}
	if request.Memory != "" {
		node.Spec.Memory = request.Memory
	}
	if request.MemoryLimit != "" {
		node.Spec.MemoryLimit = request.MemoryLimit
	}
	if request.Storage != "" {
		node.Spec.Storage = request.Storage
	}

	node, err := service.Update(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}
	return c.Status(http.StatusOK).JSON(shared.NewResponse(node))

}

// List returns all chainlink nodes
func List(c *fiber.Ctx) error {
	nodes, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	// default page to 0
	//Todo move to the service =>create pagination package
	page, _ := strconv.Atoi(c.Query("page"))
	start, end := shared.Page(uint(len(nodes.Items)), uint(page))
	sort.Slice(nodes.Items[:], func(i, j int) bool {
		return nodes.Items[j].CreationTimestamp.Before(&nodes.Items[i].CreationTimestamp)
	})

	nodeModels := new(Nodes).FromChainlinkNode(nodes.Items[start:end])
	return c.Status(http.StatusOK).JSON(shared.NewResponse(nodeModels))
}

// Delete a single chainlink node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(*chainlinkv1alpha1.Node)

	err := service.Delete(node)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

// Count returns total number of nodes
func Count(c *fiber.Ctx) error {
	nodes, err := service.List()
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodes.Items)))

	return c.SendStatus(http.StatusOK)
}

// validateNodeExist validate node by name exist
func ValidateNodeExist(c *fiber.Ctx) error {
	name := c.Params("name")

	node, err := service.Get(name)
	if err != nil {
		return c.Status(err.Status).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}
