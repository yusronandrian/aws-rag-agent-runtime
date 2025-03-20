package model

import (
	"aws-rag-agent-runtime/constant"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"

	"log"
)

type BedrockAgent struct {
	Client bedrockagentruntime.Client
}

func NewBedrock() *BedrockAgent {
	ctx := context.Background()
	var BedrockAgentRuntimeClient *bedrockagentruntime.Client

	// Load AWS Credentials from environment variables
	accessKeyID := "XXX"
	secretAccessKey := "XXX"
	region := "us-east-1"

	if accessKeyID == "" || secretAccessKey == "" || region == "" {
		log.Fatal("AWS credentials not found in environment variables")
		//Or you can set default region
		//region = "us-east-1"
	}

	// Load AWS Credentials
	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"), config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")))
	if err != nil {
		log.Fatal("Failed to load AWS credentials", err)
	}

	// create BedrockAgentClient
	BedrockAgentRuntimeClient = bedrockagentruntime.NewFromConfig(awsConfig)

	return &BedrockAgent{
		Client: *BedrockAgentRuntimeClient,
	}
}

func (bedrockAgent *BedrockAgent) RetrieveResponseFromKnowledgeBase(question string) string {
	// invoke bedrock agent runtime to retrieve and generate
	output, err := bedrockAgent.Client.RetrieveAndGenerate(
		context.TODO(),
		&bedrockagentruntime.RetrieveAndGenerateInput{
			Input: &types.RetrieveAndGenerateInput{
				Text: aws.String(question),
			},
			RetrieveAndGenerateConfiguration: &types.RetrieveAndGenerateConfiguration{
				Type: types.RetrieveAndGenerateTypeKnowledgeBase,
				KnowledgeBaseConfiguration: &types.KnowledgeBaseRetrieveAndGenerateConfiguration{
					KnowledgeBaseId: aws.String(constant.KnowledgeBaseId),
					ModelArn:        aws.String(constant.ModelArn),
					RetrievalConfiguration: &types.KnowledgeBaseRetrievalConfiguration{
						VectorSearchConfiguration: &types.KnowledgeBaseVectorSearchConfiguration{
							NumberOfResults: aws.Int32(6),
						},
					},
					// Add guardrail configuration
					// GenerationConfiguration: &types.GenerationConfiguration{
					// 	GuardrailConfiguration: &types.GuardrailConfiguration{
					// 		GuardrailId:      aws.String(constant.GuardrailId),
					// 		GuardrailVersion: aws.String(constant.GuardrailVersion),
					// 	},
					// },
				},
			},
		},
	)

	if err != nil {
		log.Fatal("RetrieveResponseFromKnowledgeBase::", err)
	}

	// Check if the response was blocked by guardrails
	// if output.Output.GuardrailAction != nil && output.Output.GuardrailAction.Action == types.GuardrailActionTypeBlocked {
	// 	return "I'm sorry, but I cannot provide a response to that query due to content guardrails."
	// }

	result := output.Output.Text
	return *result
}
