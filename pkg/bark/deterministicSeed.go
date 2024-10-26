package bark

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IceCodeNew/telesend/internal/app/aead"
	"github.com/IceCodeNew/telesend/internal/app/config"
	"github.com/IceCodeNew/telesend/pkg/crypto"
	"github.com/IceCodeNew/telesend/pkg/random"
	"lukechampine.com/frand"
)

func (sender *BarkSender) deterministicKeyAndNonce() (key, nonce []byte, err error) {
	_seed, _token, found := strings.Cut(config.TSConfig.BotToken, ":")
	if !found {
		return nil, nil, fmt.Errorf("ERROR: [Internal] Invalid bot token format")
	}

	passphrase := append(
		make([]byte, 0, len(_token)+len(sender.ID)),
		_token...,
	)
	passphrase = append(passphrase, sender.ID...)

	seed1, err := strconv.ParseUint(_seed, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"ERROR: [Internal] Failed to parse seed into uint64: %s", _seed,
		)
	}
	rng := frand.NewCustom(
		deterministicSeed(seed1, uint64(sender.Creator)),
		-1, -1,
	)
	salt := make([]byte, crypto.KeySizeAES128)
	_, _ = rng.Read(salt)

	_kn := aead.DeriveKey(passphrase, salt, crypto.KeySizeAES128<<1)
	key, nonce, _kn = _kn[:crypto.KeySizeAES128], _kn[crypto.KeySizeAES128:], nil
	return key, nonce, nil
}

// deterministicSeed returns a deterministic seed,
// which length is always 32 bytes.
func deterministicSeed(seed1, seed2 uint64) []byte {
	seed1, seed2 = enlargeSeed(seed1), enlargeSeed(seed2)
	// It does not matter if the product overflowed
	//
	// The final result of the product would length at least 16 digits, at most 20 digits.
	// math.MaxUint64: 18446744073709551615
	product := seed1 * seed2
	for i := 0; product < uint64(1000000000000000); i++ {
		product *= seed2
		if product == 0 {
			product = seed1 + seed2 - minimumSeed<<1
			product = enlargeSeed(uint64(i) * product)
		}
	}

	_strProduct := strconv.FormatUint(product, 10)
	seed := append(make([]byte, 0, seedLen), _strProduct...)

	// make sure the _randomWordIndex is in the range of [0, 1023]
	_randomWordIndex := int(product & random.Mask_1023)
	for _padLen := seedLen - len(_strProduct); _padLen > 0; _randomWordIndex++ {
		_randomWord := random.WordList[_randomWordIndex&random.Mask_1023]
		_randomWordLen := min(len(_randomWord), _padLen)

		seed = append(seed, _randomWord[:_randomWordLen]...)
		_padLen -= _randomWordLen
	}

	return seed
}

func enlargeSeed(seed uint64) uint64 {
	if seed >= minimumSeed {
		return seed
	} else if seed < 2 {
		return minimumSeed
	}
	seed *= seed
	return enlargeSeed(seed)
}

const (
	minimumSeed uint64 = 1<<20 - 1
	seedLen            = 32
)
