package main

// Letter represents all the types of letters that can be sent
type Letter int

const (
	// NoLetter when no next step is available
	NoLetter Letter = 0

	// Non-rollover

	// AwardFSMAndCG letter
	AwardFSMAndCG Letter = 1

	// AwardCG letter
	AwardCG Letter = 2

	// AwardFSM letter
	AwardFSM Letter = 3

	// AwardCGAndRequestConsent letter
	AwardCGAndRequestConsent Letter = 4

	// RequestConsent letter
	RequestConsent Letter = 5

	// Rollover

	// RolloverFSMAndCG letter
	RolloverFSMAndCG Letter = 6

	// RolloverCG letter
	RolloverCG Letter = 7

	// RolloverFSM letter
	RolloverFSM Letter = 8

	// RolloverCGAndRequestConsent letter
	RolloverCGAndRequestConsent Letter = 9

	// RolloverRequestConsent letter
	RolloverRequestConsent Letter = 10
)

func (l Letter) String() string {
	return []string{
		"",
		"3. Award FSM and CG",
		"2. Award CG",
		"4. Award FSM",
		"1. Award CG + request consent",
		"5. Request consent",
		"6. Rollover FSM and CG",
		"8. Rollover CG",
		"7. Rollover FSM",
		"9. Rollover CG + request consent",
		"10. Rollover Request consent",
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

			if !d.ExistingCG && d.NewCG {
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
