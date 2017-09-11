#!/bin/sh
curl -H "Content-Type: application/json" -X POST -d '{"name":"Anamoly Detection", "description": "Feature for Anamoly Detection"}' http://localhost:8800/features

