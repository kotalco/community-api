package bitcoin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kotalco/community-api/internal/bitcoin"
	"github.com/kotalco/community-api/internal/core/secret"
	restErrors "github.com/kotalco/community-api/pkg/errors"
	"github.com/kotalco/community-api/pkg/k8s"
	"github.com/kotalco/community-api/pkg/shared"
	bitcoinv1alpha1 "github.com/kotalco/kotal/apis/bitcoin/v1alpha1"
	"github.com/ybbus/jsonrpc/v2"
	apiError "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type request struct {
	endpoint string
	method   string
	name     string
}
type result struct {
	err  error
	name string
	data interface{}
}

const (
	nameKeyword = "name"
)

var (
	service       = bitcoin.NewBitcoinService()
	secretService = secret.NewSecretService()
	k8sClient     = k8s.NewClientService()
)

// Get returns a single bitcoin node by name
func Get(c *fiber.Ctx) error {
	node := c.Locals("node").(bitcoinv1alpha1.Node)
	return c.JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
}

// List returns all bitcoin nodes
func List(c *fiber.Ctx) error {
	// default page to 0
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	nodeList, err := service.List(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	start, end := shared.Page(uint(len(nodeList.Items)), uint(page), uint(limit))
	sort.Slice(nodeList.Items[:], func(i, j int) bool {
		return nodeList.Items[j].CreationTimestamp.Before(&nodeList.Items[i].CreationTimestamp)
	})

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", len(nodeList.Items)))

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(bitcoin.BitcoinListDto).FromBitcoinNode(nodeList.Items[start:end])))
}

// Create created bitcoin node from given specs
func Create(c *fiber.Ctx) error {
	dto := new(bitcoin.BitcoinDto)
	if err := c.BodyParser(dto); err != nil {
		badReqErr := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReqErr.StatusCode()).JSON(badReqErr)
	}

	dto.Namespace = c.Locals("namespace").(string)
	if err := dto.MetaDataDto.Validate(); err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	//check for bitcoin json rpc default user secret
	_, err := secretService.Get(types.NamespacedName{
		Name:      bitcoin.BitcoinJsonRpcDefaultUserPasswordName,
		Namespace: dto.Namespace,
	})
	if err != nil {
		if err.StatusCode() != http.StatusNotFound {
			return c.Status(err.StatusCode()).JSON(err)
		}
		//create bitcoin user default secret
		_, err = secretService.Create(secret.SecretDto{
			MetaDataDto: k8s.MetaDataDto{Name: bitcoin.BitcoinJsonRpcDefaultUserPasswordName, Namespace: dto.Namespace},
			Type:        "password",
			Data:        map[string]string{"password": bitcoin.BitcoinJsonRpcDefaultUserPasswordSecret},
		})
		if err != nil {
			return c.Status(err.StatusCode()).JSON(err)
		}
	}

	node, err := service.Create(*dto)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}
	return c.Status(http.StatusCreated).JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
}

// Update updates a single bitcoin node by name from spec
func Update(c *fiber.Ctx) error {
	dto := new(bitcoin.BitcoinDto)
	if err := c.BodyParser(dto); err != nil {
		badReq := restErrors.NewBadRequestError("invalid request body")
		return c.Status(badReq.StatusCode()).JSON(badReq)
	}

	node := c.Locals("node").(bitcoinv1alpha1.Node)

	err := service.Update(*dto, &node)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(shared.NewResponse(new(bitcoin.BitcoinDto).FromBitcoinNode(node)))
}

// Count returns total number of nodes
func Count(c *fiber.Ctx) error {
	length, err := service.Count(c.Locals("namespace").(string))
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprintf("%d", length))

	return c.SendStatus(http.StatusOK)
}

// Delete a single bitcoin node by name
func Delete(c *fiber.Ctx) error {
	node := c.Locals("node").(bitcoinv1alpha1.Node)

	err := service.Delete(&node)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	return c.SendStatus(http.StatusNoContent)
}

func ValidateNodeExist(c *fiber.Ctx) error {
	nameSpacedName := types.NamespacedName{
		Name:      c.Params(nameKeyword),
		Namespace: c.Locals("namespace").(string),
	}

	node, err := service.Get(nameSpacedName)
	if err != nil {
		return c.Status(err.StatusCode()).JSON(err)
	}

	c.Locals("node", node)
	return c.Next()
}

// Stats returns a websocket that emits bitcoin block and node count stats
func Stats(c *websocket.Conn) {
	defer c.Close()
	name := c.Params("name")
	node := &bitcoinv1alpha1.Node{}
	nameSpacedName := types.NamespacedName{
		Namespace: c.Locals("namespace").(string),
		Name:      name,
	}
	err := k8sClient.Get(context.Background(), nameSpacedName, node)
	if err != nil {
		if apiError.IsNotFound(err) {
			c.WriteJSON(fiber.Map{
				"error": fmt.Sprintf("node by name %s doesn't exist", name),
			})
			return
		}
		c.WriteJSON(fiber.Map{
			"error": err.Error(),
		})
		return
	}

	if !node.Spec.RPC {
		c.WriteJSON(fiber.Map{
			"error": "JSON-RPC server is not enabled",
		})
		return
	}

	for {
		jobs := make(chan request, 2)
		results := make(chan result, 2)

		for i := 0; i < 2; i++ {
			go worker(jobs, results)
		}

		endpoint := fmt.Sprintf("http://%s:%s@%s.%s:%d/", bitcoin.BitcoinJsonRpcDefaultUserName, bitcoin.BitcoinJsonRpcDefaultUserPasswordSecret, nameSpacedName.Name, nameSpacedName.Namespace, node.Spec.RPCPort)
		jobs <- request{name: "blockCount", endpoint: endpoint, method: "getblockcount"}
		jobs <- request{name: "peerCount", endpoint: endpoint, method: "getconnectioncount"}

		close(jobs)

		var bitcoinStatResponseDto struct {
			BlockCount int64 `json:"blockCount"`
			PeerCount  int64 `json:"peerCount"`
		}

		newBitcoinResponseDto := bitcoinStatResponseDto

		for i := 0; i < 2; i++ {
			resp := <-results
			if resp.err != nil {
				c.WriteJSON(fiber.Map{
					"error": resp.err,
				})
				return
			}

			switch resp.name {
			case "blockCount":
				newBitcoinResponseDto.BlockCount, err = resp.data.(json.Number).Int64()
				if err != nil {
					c.WriteJSON(fiber.Map{
						"error": err.Error(),
					})
					return
				}
				break
			case "peerCount":
				newBitcoinResponseDto.PeerCount, err = resp.data.(json.Number).Int64()
				if err != nil {
					c.WriteJSON(fiber.Map{
						"error": err.Error(),
					})
					return
				}
				break
			}
		}
		close(results)

		err := c.WriteJSON(newBitcoinResponseDto)
		if err != nil {
			return
		}

		time.Sleep(time.Second * 3)
	}
}

// worker is a  collection of threads for the bitcoin stats
func worker(jobs <-chan request, results chan<- result) {
	chanRes := result{}
	for job := range jobs {
		chanRes.name = job.name

		client := jsonrpc.NewClient(job.endpoint)

		res, err := client.Call(job.method)
		if err != nil {
			chanRes.err = err
		} else {
			if res.Error != nil {
				chanRes.err = errors.New(res.Error.Message)
			} else {
				chanRes.data = res.Result
			}
		}

		results <- chanRes
	}
}
