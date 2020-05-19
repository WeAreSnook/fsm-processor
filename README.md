# fsm-processor

Parser script, takes spreadsheets as input and outputs csv results containing which matched pupils should be given free school meals and clothing grant. Used by fsm-app.

The input spreadsheets are:
- Current Awards
- Consent 360
- Dependents SHBE
- Benefit Extract
- School roll
- Universal credit

The output spreadsheets are:
- `report_awards_fsm.csv` - people who have given consent and qualify for FSM and/or CG
- `report_awards_ctr.csv` - people who have not given consent and qualify for CG based on CTR
- `report_education_fsm.csv` - people who couldn't be matched to the school roll when generating report_awards_fsm.csv
- `report_education_ctr.csv` - people who couldn't be matched to the schoo lroll when generating report_awards_ctr.csv

# Usage

Once built with `go build` you can run the processor with the following inputs:
```
  -awardcg
    	if we should award CG (default true)
  -awards string
    	filepath for current awards spreadsheet
  -benefitamount float
    	benefit amount (default 610)
  -benefitextract string
    	filepath for benefit extract spreadsheet
  -consent string
    	filepath for consent spreadsheet
  -ctcfigure float
    	ctc annual income figure (default 16105)
  -ctcwtcfigure float
    	ctc/wtc annual income figure (default 6420)
  -debugclaim int
    	claimnumber to output debug logs for (default -1)
  -dependents string
    	filepath for dependents SHBE spreadsheet
  -dev
    	development mode, use spreadsheets from ./private-data folder without having to specify each one
  -filter string
    	filepath for filter spreadsheet
  -log
    	log output to stdout (for debugging, breaks json output parsing)
  -output string
    	path of the folder outputs should be stored in (default "./")
  -rollover
    	rollover mode
  -schoolroll string
    	filepath for school roll spreadsheet
  -universalcredit string
    	filepath for universal credit spreadsheet
 ```

# Implementation

The app is split into 2 main packages, see the output of `go doc` for [details](./DOCS.md).

### spreadsheet

Input spreadsheets come in various forms, including: csv, tsv, xlsx, xls. This package abstracts the details of working with each individual format to allow them to be treated the same. Also provides convenience functions to e.g. fetch columns by name, override column names, validate certain column names exist in the input, and convert cell values to various types.

### main

Runs the checks to determine who gets FSM and/or CG based on the input spreadsheets.
