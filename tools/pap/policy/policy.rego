package policy

import future.keywords

default authz = false

# allow staging data access for everyone
# OR
# allow HR_PII for HR in same location
# OR
# allow CUSTOMERS_PII in same location if whitelisted in data
