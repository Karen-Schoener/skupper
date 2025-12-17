Execution Workflow
Here is how you would execute the new broken-down setup:

Clean up everything:

Bash

./cleanup.sh
Run setup for each cluster manually:

Bash

./setup-cluster-N.sh 1
# Check Skupper status, resolve any hang issues, then proceed:
./setup-cluster-N.sh 2
# Check Skupper status, resolve any hang issues, then proceed:
./setup-cluster-N.sh 3
# Check Skupper status, resolve any hang issues, then proceed:
./setup-cluster-N.sh 4
Perform the hub connection:

Bash

./connect-hub.sh
Link Setup: Open 4 shells and run:

Bash

source dedicate-aliases.sh
# Follow the dedication, token creation, and linking steps manually.
