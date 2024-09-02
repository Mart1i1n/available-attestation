package actionset

var (
	_ Action = storeSignedAttestAction{}
	_ Action = rePackAttestationAction{}
)

type storeSignedAttestAction struct {
	baseCmdAction
}

func (s storeSignedAttestAction) Desc() string { return "# Store signed attestation" }

func (s storeSignedAttestAction) Name() string { return "storeSignedAttest" }

func (s storeSignedAttestAction) ActionType() AType { return AttestAction }

type rePackAttestationAction struct {
	baseCmdAction
}

func (r rePackAttestationAction) Desc() string { return "# Re-pack attestation" }

func (r rePackAttestationAction) Name() string { return "rePackAttestation" }

func (r rePackAttestationAction) ActionType() AType { return AttestAction }
