package main

// Letter represents all the types of letters that can be sent
type Letter int

const (
	// NoLetter when no next step is available
	NoLetter Letter = iota

	// Non-rollover

	// AwardFSMAndCG letter
	AwardFSMAndCG

	// AwardCG letter
	AwardCG

	// AwardFSM letter
	AwardFSM

	// AwardCGAndRequestConsent letter
	AwardCGAndRequestConsent

	// RequestConsent letter
	RequestConsent

	// Rollover

	// RolloverFSMAndCG letter
	RolloverFSMAndCG

	// RolloverCG letter
	RolloverCG

	// RolloverFSM letter
	RolloverFSM

	// RolloverCGAndRequestConsent letter
	RolloverCGAndRequestConsent

	// RolloverRequestConsent letter
	RolloverRequestConsent
)

func (l Letter) String() string {
	return []string{
		"",
		"AwardFSMAndCG",
		"AwardCG",
		"AwardFSM",
		"AwardCGAndRequestConsent",
		"RequestConsent",
		"RolloverFSMAndCG",
		"RolloverCG",
		"RolloverFSM",
		"RolloverCGAndRequestConsent",
		"RolloverRequestConsent",
	}[l]
}

// LetterForDependent returns the next-step letter for the given dependent
func LetterForDependent(d Dependent, rollover bool) Letter {
	consent := d.Person.HasConsent()

	if rollover {
		if consent {
			if d.NewFSM && d.NewCG {
				return RolloverFSMAndCG
			}

			if !d.NewFSM && d.NewCG {
				return RolloverCG
			}

			if d.NewFSM && !d.NewCG {
				return RolloverFSM
			}
		} else {
			if d.NewCG {
				return RolloverCGAndRequestConsent
			}

			if d.ExistingFSM {
				return RolloverRequestConsent
			}
		}
	} else {
		if consent {
			if !d.ExistingFSM && d.NewFSM && !d.ExistingCG && d.NewCG {
				return AwardFSMAndCG
			}

			if (d.ExistingFSM || d.NewFSM) && d.NewCG {
				return AwardCG
			}

			if d.NewFSM {
				return AwardFSM
			}
		} else {
			if d.NewCG {
				return AwardCGAndRequestConsent
			}

			if !d.ExistingFSM && d.ExistingCG {
				return RequestConsent
			}
		}

	}

	return NoLetter
}
