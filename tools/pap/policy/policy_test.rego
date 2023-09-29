package policy

test_default_customer_pii {
	input := {
		"Resource": {
			"Classification": "CUSTOMERS_PII",
			"Location": "Belgium",
		},
		"Subject": {
			"Email": "hr@klarrio.com",
			"WorkingLocation": "Belgium",
		},
	}

	not authz with input as input
}

test_badlocation_customer_pii {
	input := {
		"Resource": {
			"Classification": "CUSTOMERS_PII",
			"Location": "Belgium",
		},
		"Subject": {
			"Email": "lander.visterin@klarrio.com",
			"WorkingLocation": "USA",
		},
	}

	testlist := ["lander.visterin@klarrio.com"]

	not authz with input as input
		with data.customerDataAccess as testlist
}

test_allowlist_customer_pii {
	input := {
		"Resource": {
			"Classification": "CUSTOMERS_PII",
			"Location": "Belgium",
		},
		"Subject": {
			"Email": "lander.visterin@klarrio.com",
			"WorkingLocation": "Belgium",
		},
	}

	testlist := ["lander.visterin@klarrio.com"]

	authz with input as input
		with data.customerDataAccess as testlist
}
