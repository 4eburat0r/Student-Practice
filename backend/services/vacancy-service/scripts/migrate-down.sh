#!/bin/bash
set -e

migrate -path ./migrations -database "$DATABASE_URL" down
