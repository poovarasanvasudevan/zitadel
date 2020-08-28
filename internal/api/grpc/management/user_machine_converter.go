package management

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func machineCreateToModel(machine *management.CreateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		Name:        machine.Name,
		Description: machine.Description,
	}
}

func updateMachineToModel(machine *management.UpdateMachineRequest) *usr_model.Machine {
	return &usr_model.Machine{
		ObjectRoot:  models.ObjectRoot{AggregateID: machine.Id},
		Description: machine.Description,
	}
}

func machineFromModel(account *usr_model.Machine) *management.MachineResponse {
	return &management.MachineResponse{
		Name:        account.Name,
		Description: account.Description,
	}
}

func machineViewFromModel(machine *usr_model.MachineView) *management.MachineView {
	lastKeyAdded, err := ptypes.TimestampProto(machine.LastKeyAdded)
	logging.Log("MANAG-wGcAQ").OnError(err).Debug("unable to parse date")
	return &management.MachineView{
		Description:  machine.Description,
		Name:         machine.Name,
		LastKeyAdded: lastKeyAdded,
		Keys:         machineKeyViewsFromModel(machine.Keys...),
	}
}

func machineKeyViewsFromModel(keys ...*usr_model.MachineKeyView) []*management.MachineKeyView {
	keyViews := make([]*management.MachineKeyView, len(keys))
	for i, key := range keys {
		keyViews[i] = machineKeyViewFromModel(key)
	}
	return keyViews
}

func machineKeyViewFromModel(key *usr_model.MachineKeyView) *management.MachineKeyView {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	expirationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-gluk7").OnError(err).Debug("unable to parse timestamp")

	return &management.MachineKeyView{
		Id:             key.ID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		Type:           machineKeyTypeFromModel(key.Type),
	}
}

func addMachineKeyToModel(key *management.AddMachineKeyRequest) *usr_model.MachineKey {
	expirationDate, err := ptypes.Timestamp(key.ExpirationDate)
	logging.Log("MANAG-iNshR").OnError(err).Debug("unable to parse expiration date")

	return &usr_model.MachineKey{
		ExpirationDate: expirationDate,
		Type:           machineKeyTypeToModel(key.Type),
	}
}

func addMachineKeyFromModel(key *usr_model.MachineKey) *management.AddMachineKeyResponse {
	creationDate, err := ptypes.TimestampProto(key.CreationDate)
	logging.Log("MANAG-dlb8m").OnError(err).Debug("unable to parse cretaion date")

	expirationDate, err := ptypes.TimestampProto(key.ExpirationDate)
	logging.Log("MANAG-dlb8m").OnError(err).Debug("unable to parse cretaion date")

	details := []byte(`{"type": "serviceaccount, "privateKeyId":}`)

	return &management.AddMachineKeyResponse{
		Id:             key.KeyID,
		CreationDate:   creationDate,
		ExpirationDate: expirationDate,
		Sequence:       key.Sequence,
		KeyDetails:     details,
		Type:           machineKeyTypeFromModel(key.Type),
	}
}

func machineKeyTypeToModel(typ management.MachineKeyType) usr_model.MachineKeyType {
	switch typ {
	case management.MachineKeyType_MACHINEKEY_JSON:
		return usr_model.MachineKeyTypeJSON
	default:
		return usr_model.MachineKeyTypeNONE
	}
}

func machineKeyTypeFromModel(typ usr_model.MachineKeyType) management.MachineKeyType {
	switch typ {
	case usr_model.MachineKeyTypeJSON:
		return management.MachineKeyType_MACHINEKEY_JSON
	default:
		return management.MachineKeyType_MACHINEKEY_UNSPECIFIED
	}
}