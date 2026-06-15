package api

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/stellarlinkco/agentsdk-go/pkg/config"
	"github.com/stellarlinkco/agentsdk-go/pkg/skylark"
	"github.com/stellarlinkco/agentsdk-go/pkg/tool"
)

func buildKnowledgeEngine(ctx context.Context, opts Options) (*skylark.Engine, error) {
	if opts.Knowledge == nil || !opts.Knowledge.Enabled {
		return nil, nil
	}
	k := opts.Knowledge
	indexDir := strings.TrimSpace(k.IndexDir)
	if indexDir == "" {
		indexDir = filepath.Join(opts.ProjectRoot, ".agents", "knowledge-index")
	}

	var emb = k.Embedder
	if emb == nil && !k.DisableEmbedding {
		var err error
		emb, err = skylark.NewEmbedderFromEnv()
		if err != nil {
			return nil, fmt.Errorf("knowledge: embedder: %w", err)
		}
	}

	eng, err := skylark.NewEngine(indexDir, emb)
	if err != nil {
		return nil, err
	}

	vaultDir := strings.TrimSpace(k.VaultDir)
	docs, err := skylark.SyncVault(vaultDir)
	if err != nil {
		_ = eng.Close()
		return nil, fmt.Errorf("knowledge: sync vault: %w", err)
	}
	if err := eng.Rebuild(ctx, docs); err != nil {
		_ = eng.Close()
		return nil, fmt.Errorf("knowledge: rebuild index: %w", err)
	}
	return eng, nil
}

func registerKnowledgeTools(registry *tool.Registry, engine *skylark.Engine, opts Options, settings *config.Settings) error {
	if registry == nil || engine == nil {
		return fmt.Errorf("api: registerKnowledgeTools: nil argument")
	}
	dis := effectiveDisallowedToolSet(opts, settings)
	tools := []tool.Tool{
		&memorySearchTool{engine: engine},
		&sessionSearchTool{},
	}
	for _, t := range tools {
		if t == nil {
			continue
		}
		name := canonicalToolName(t.Name())
		if name == "" {
			continue
		}
		if dis != nil {
			if _, blocked := dis[name]; blocked {
				continue
			}
		}
		if err := registry.Register(t); err != nil {
			return fmt.Errorf("api: register tool %s: %w", t.Name(), err)
		}
	}
	return nil
}

// effectiveDisallowedToolSet merges Options.DisallowedTools with settings.json disallowedTools.
func effectiveDisallowedToolSet(opts Options, settings *config.Settings) map[string]struct{} {
	dis := toLowerSet(opts.DisallowedTools)
	if settings != nil && len(settings.DisallowedTools) > 0 {
		if dis == nil {
			dis = map[string]struct{}{}
		}
		for _, name := range settings.DisallowedTools {
			if key := canonicalToolName(name); key != "" {
				dis[key] = struct{}{}
			}
		}
	}
	return dis
}
