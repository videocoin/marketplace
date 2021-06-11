package listener

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// ErrUnknownEvent returned if unknown event is provided for parsing.
	ErrUnknownEvent = errors.New("unknown event")

	// ErrInvalidLog returned has unexpected format.
	ErrInvalidLog = errors.New("invalid log format")
)

// NewParser creates parser instance.
func NewParser() *Parser {
	return &Parser{
		events: map[common.Hash]abi.Event{},
	}
}

// Parser provides utilities to register parsers for different events and parse ethereum logs into map.
type Parser struct {
	// events is a map from an event id to an abi.Event
	// event id is sha3(event sig)
	events map[common.Hash]abi.Event
}

// RegisterABI registers all events from an abi definition.
func (parser *Parser) RegisterABI(data io.Reader) error {
	definition, err := abi.JSON(data)
	if err != nil {
		return err
	}
	for _, ev := range definition.Events {
		if err := parser.RegisterEvent(ev); err != nil {
			return err
		}
	}
	return nil
}

func (parser *Parser) IterateEvents(f func(common.Hash) bool) {
	for etype := range parser.events {
		if !f(etype) {
			return
		}
	}
}

func (parser *Parser) GetTypes() []common.Hash {
	rst := make([]common.Hash, 0, len(parser.events))
	for key := range parser.events {
		rst = append(rst, key)
	}
	return rst
}

// RegisterEvent register single abi event.
func (parser *Parser) RegisterEvent(event abi.Event) error {
	_, exist := parser.events[event.ID]
	if exist {
		return nil
	}
	parser.events[event.ID] = event
	return nil
}

// GetEventName if event is registered returns event name as it was declared in ABI.
func (parser *Parser) GetEventName(eventID common.Hash) string {
	event, exist := parser.events[eventID]
	if !exist {
		return ""
	}
	return event.RawName
}

func (parser *Parser) Unpack(out interface{}, log *types.Log) error {
	if len(log.Topics) == 0 {
		return ErrInvalidLog
	}
	eventParser, exist := parser.events[log.Topics[0]]
	if !exist {
		return fmt.Errorf("%w: event %v", ErrUnknownEvent, log)
	}
	if len(log.Data) > 0 {
		values, err := eventParser.Inputs.Unpack(log.Data)
		if err != nil {
			return err
		}
		if err := eventParser.Inputs.Copy(out, values); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range eventParser.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return parseTopics(out, indexed, log.Topics[1:])
}

// Big batch of reflect types for topic reconstruction.
var (
	reflectHash    = reflect.TypeOf(common.Hash{})
	reflectAddress = reflect.TypeOf(common.Address{})
	reflectBigInt  = reflect.TypeOf(new(big.Int))
)

// capitalise makes a camel-case string which starts with an upper case character.
func capitalise(input string) string {
	return abi.ToCamelCase(input)
}

// parseTopics converts the indexed topic fields into actual log field values.
//
// Note, dynamic types cannot be reconstructed since they get mapped to Keccak256
// hashes as the topic value!
func parseTopics(out interface{}, fields abi.Arguments, topics []common.Hash) error {
	// Sanity check that the fields and topics match up
	if len(fields) != len(topics) {
		return errors.New("topic/field count mismatch")
	}
	// Iterate over all the fields and reconstruct them from topics
	for _, arg := range fields {
		if !arg.Indexed {
			return errors.New("non-indexed field in topic reconstruction")
		}
		field := reflect.ValueOf(out).Elem().FieldByName(capitalise(arg.Name))

		// Try to parse the topic back into the fields based on primitive types
		switch field.Kind() {
		case reflect.Bool:
			if topics[0][common.HashLength-1] == 1 {
				field.Set(reflect.ValueOf(true))
			}
		case reflect.Int8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int8(num.Int64())))

		case reflect.Int16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int16(num.Int64())))

		case reflect.Int32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int32(num.Int64())))

		case reflect.Int64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Int64()))

		case reflect.Uint8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint8(num.Uint64())))

		case reflect.Uint16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint16(num.Uint64())))

		case reflect.Uint32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint32(num.Uint64())))

		case reflect.Uint64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Uint64()))

		default:
			// Ran out of plain primitive types, try custom types

			switch field.Type() {
			case reflectHash: // Also covers all dynamic types
				field.Set(reflect.ValueOf(topics[0]))

			case reflectAddress:
				var addr common.Address
				copy(addr[:], topics[0][common.HashLength-common.AddressLength:])
				field.Set(reflect.ValueOf(addr))

			case reflectBigInt:
				num := new(big.Int).SetBytes(topics[0][:])
				if arg.Type.T == abi.IntTy {
					if num.Cmp(abi.MaxInt256) > 0 {
						num.Add(abi.MaxUint256, big.NewInt(0).Neg(num))
						num.Add(num, big.NewInt(1))
						num.Neg(num)
					}
				}
				field.Set(reflect.ValueOf(num))

			default:
				// Ran out of custom types, try the crazies
				switch {
				// static byte array
				case arg.Type.T == abi.FixedBytesTy:
					reflect.Copy(field, reflect.ValueOf(topics[0][:arg.Type.Size]))
				default:
					return fmt.Errorf("unsupported indexed type: %v", arg.Type)
				}
			}
		}
		topics = topics[1:]
	}
	return nil
}
