#!/bin/bash

if [ "${CIRCLE_BRANCH}" == "main" ]; then
    export TF_VAR_ENVIRONMENT=production
else
    export TF_VAR_ENVIRONMENT=test
fi