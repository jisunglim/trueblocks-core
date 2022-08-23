package pinning

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/cache"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/index/bloom"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type PinResult struct {
	Range   cache.FileRange         `json:"range,omitempty"`
	Local   types.SimpleChunkRecord `json:"local,omitempty"`
	Remote  types.SimpleChunkRecord `json:"remote,omitempty"`
	Matches bool                    `json:"matches,omitempty"`
	err     error
}

func PinChunk(chain, path string, isRemote bool) (PinResult, error) {
	bloomFile := bloom.ToBloomPath(path)
	indexFile := index.ToIndexPath(path)

	rng, _ := cache.RangeFromFilename(path)
	result := PinResult{
		Range:  rng,
		Local:  types.SimpleChunkRecord{Range: rng.String()},
		Remote: types.SimpleChunkRecord{Range: rng.String()},
	}

	localService, _ := NewPinningService(chain, Local)
	remoteService, _ := NewPinningService(chain, Pinata)

	isLocal := LocalDaemonRunning()
	if isLocal {
		if result.Local.BloomHash, result.err = localService.PinFile(bloomFile, true); result.err != nil {
			return PinResult{}, result.err
		}
		if result.Local.IndexHash, result.err = localService.PinFile(indexFile, true); result.err != nil {
			return PinResult{}, result.err
		}
	}

	if isRemote {
		if result.Remote.BloomHash, result.err = remoteService.PinFile(bloomFile, false); result.err != nil {
			return PinResult{}, result.err
		}
		if result.Remote.IndexHash, result.err = remoteService.PinFile(indexFile, false); result.err != nil {
			return PinResult{}, result.err
		}
	}

	// TODO: We used to use this to report an error between local and remote pinning, but it got turned off. Turn it back on
	if isLocal && isRemote {
		result.Matches = result.Local.IndexHash == result.Remote.IndexHash && result.Local.BloomHash == result.Remote.BloomHash
	} else {
		result.Matches = true
	}

	return result, nil
}

func (p *Service) PinFile(filepath string, local bool) (string, error) {
	if local {
		return p.pinFileLocally(filepath)
	}
	return p.pinFileRemotely(filepath)
}
