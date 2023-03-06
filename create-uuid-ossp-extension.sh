#!/bin/bash

# Execute a shell in the container using kubectl
kubectl exec -it postgres-deployment-0 -n logfire-local sh <<EOF

# Run psql in the container with the appropriate parameters and execute the SQL command
psql -d logfire-dev -U postgresuser -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

EOF

