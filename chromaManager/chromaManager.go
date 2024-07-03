package chromaManager

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	chromem "github.com/philippgille/chromem-go"
)

func init() {
	godotenv.Load(".env")
}

type ChromaManager struct {
	client      *chromem.DB
	collections map[string]*chromem.Collection
	defaultEF   chromem.EmbeddingFunc
}

func NewChromaManager() (*ChromaManager, error) {
	path := os.Getenv("EMBED_DIR")
	if path == "" {
		path = ".docs_db"
	}

	client, err := chromem.NewPersistentDB(path, false)
	if err != nil {
		return nil, err
	}

	ollamaEmbURL := os.Getenv("OLLAMA_EMBED_URL")

	ollamaEmbMode := os.Getenv("OLLAMA_EMBED_MODE")

	var defaultEF chromem.EmbeddingFunc
	if ollamaEmbURL != "" && ollamaEmbMode != "" {
		log.Printf("using ollama embedded %s %s", ollamaEmbURL, ollamaEmbMode)
		defaultEF = chromem.NewEmbeddingFuncOllama(ollamaEmbMode, ollamaEmbURL)
	} else {
		defaultEF = chromem.NewEmbeddingFuncDefault()
	}

	return &ChromaManager{
		client:      client,
		collections: make(map[string]*chromem.Collection),
		defaultEF:   defaultEF,
	}, nil
}

func (cm *ChromaManager) GetOrCreateCollection(collectionName string) (*chromem.Collection, error) {
	if collection, ok := cm.collections[collectionName]; ok {
		return collection, nil
	}

	collection, err := cm.client.GetOrCreateCollection(collectionName, nil, cm.defaultEF)
	if err != nil {
		return nil, err
	}

	cm.collections[collectionName] = collection
	return collection, nil
}

func (cm *ChromaManager) UpsertDocuments(collectionName string, documents []string, metadatas []map[string]string, ids []string) error {
	collection, err := cm.GetOrCreateCollection(collectionName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	for i := range documents {
		// Generate embedding for the document
		embeddings, err := cm.defaultEF(ctx, documents[i])
		if err != nil {
			return err
		}

		doc, err := chromem.NewDocument(ctx, ids[i], metadatas[i], embeddings, documents[i], cm.defaultEF)
		if err != nil {
			return err
		}
		if err := collection.AddDocument(ctx, doc); err != nil {
			return err
		}
	}
	return nil
}

func (cm *ChromaManager) QueryCollection(collectionName, queryText string, nResults int) ([]chromem.Result, error) {
	collection, err := cm.GetOrCreateCollection(collectionName)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	// Adding empty maps for where and whereDocument parameters
	return collection.Query(ctx, queryText, nResults, map[string]string{}, map[string]string{})
}

func (cm *ChromaManager) UpsertDoc(collectionName, text, id string, metadata map[string]string) error {
	collection, err := cm.GetOrCreateCollection(collectionName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	// Generate embedding for the document
	embeddings, err := cm.defaultEF(ctx, text)
	if err != nil {
		return err
	}

	doc, err := chromem.NewDocument(ctx, id, metadata, embeddings, text, cm.defaultEF)
	if err != nil {
		return err
	}
	return collection.AddDocument(ctx, doc)
}

func (cm *ChromaManager) DeleteByID(collectionName string, id string) error {
	collection, err := cm.GetOrCreateCollection(collectionName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return collection.Delete(ctx, nil, nil, id)
}

func (cm *ChromaManager) DeleteByIDs(collectionName string, ids []string) error {
	collection, err := cm.GetOrCreateCollection(collectionName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return collection.Delete(ctx, nil, nil, ids...)
}

func (cm *ChromaManager) Search(collectionName, query string, topk int) ([]chromem.Result, error) {
	collection, err := cm.GetOrCreateCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if topk > collection.Count() {
		topk = collection.Count()
	}

	ctx := context.Background()
	// Adding empty maps for where and whereDocument parameters
	results, err := collection.Query(ctx, query, topk, map[string]string{}, map[string]string{})
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return results, nil
	}
	return nil, nil
}
