package stabilityai

import (
	"context"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	pb "github.com/insomniacslk/stabilityai/generation"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

// some client connection defaults.
const (
	DefaultAPIHost = "grpc.stability.ai:443"
	DefaultEngine  = "stable-diffusion-v1"
)

// Default parameters for image generation.
var (
	DefaultSteps    uint64  = 50
	DefaultSamples  uint64  = 1
	DefaultCfgScale float32 = 7.0
)

// Client is a Stability AI client.
type Client struct {
	ctx    context.Context
	log    *log.Logger
	host   string
	apikey string
	engine string
	client pb.GenerationServiceClient
}

// NewClient returns a new Stability AI Client object.
func NewClient(opts ...Option) *Client {
	c := Client{
		ctx:    context.Background(),
		host:   DefaultAPIHost,
		engine: DefaultEngine,
		log:    log.New(os.Stderr, "stabilityai: ", log.LstdFlags),
	}
	for _, o := range opts {
		o(&c)
	}
	return &c
}

// Connect connects to the gRPC endpoint of Stability AI.
func (c *Client) Connect() error {
	token := oauth2.Token{
		AccessToken: c.apikey,
	}
	perRPCToken := oauth.NewOauthAccess(&token)
	creds := credentials.NewTLS(nil)
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(perRPCToken),
		grpc.WithTransportCredentials(creds),
	}
	c.log.Printf("Dialling %s", c.host)
	conn, err := grpc.Dial(c.host, opts...)
	if err != nil {
		return err
	}
	c.client = pb.NewGenerationServiceClient(conn)
	return nil
}

// Generate calls the Stability AI gRPC method Generate.
func (c *Client) Generate(req *pb.Request) ([]*pb.Answer, error) {
	if c.client == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}
	c.log.Printf("Request: %+v", req)
	stream, err := c.client.Generate(c.ctx, req)
	if err != nil {
		return nil, err
	}
	answers := make([]*pb.Answer, 0)
	for {
		ans, err := stream.Recv()
		if err == io.EOF {
			return answers, nil
		}
		if err != nil {
			return nil, err
		}
		if len(ans.Artifacts) > 0 {
			answers = append(answers, ans)
		}
	}
}

// GenerateImage generates an image using the Stability AI API with some default
// values.
func (c *Client) GenerateImage(text string, width, height uint64) ([]*pb.Answer, error) {
	// TODO make random seeder configurable
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	seed := r.Intn(0xffffffff)

	req := pb.Request{
		EngineId:  c.engine,
		RequestId: uuid.New().String(),
		Prompt: []*pb.Prompt{
			&pb.Prompt{Prompt: &pb.Prompt_Text{Text: text}},
		},
		Params: &pb.Request_Image{
			Image: &pb.ImageParameters{
				Width:   &width,
				Height:  &height,
				Seed:    []uint32{uint32(seed)},
				Steps:   &DefaultSteps,
				Samples: &DefaultSamples,
				Transform: &pb.TransformType{
					Type: &pb.TransformType_Diffusion{
						Diffusion: pb.DiffusionSampler_SAMPLER_K_LMS,
					},
				},
				Parameters: []*pb.StepParameter{
					&pb.StepParameter{
						Sampler: &pb.SamplerParameters{
							CfgScale: &DefaultCfgScale,
						},
					},
				},
			},
		},
		RequestedType: pb.ArtifactType_ARTIFACT_IMAGE,
	}
	return c.Generate(&req)
}
