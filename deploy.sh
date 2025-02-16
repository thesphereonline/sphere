#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Deploying Sphere Blockchain...${NC}"

# Deploy backend to Railway
echo -e "${GREEN}Deploying backend to Railway...${NC}"
railway up

# Get the backend URL
BACKEND_URL=$(railway show | grep "URL:" | cut -d' ' -f2)
echo -e "${GREEN}Backend deployed to: $BACKEND_URL${NC}"

# Update frontend configuration
echo -e "${GREEN}Updating frontend configuration...${NC}"
echo "REACT_APP_API_URL=$BACKEND_URL/api" > frontend/.env.production

# Deploy frontend to Netlify
echo -e "${GREEN}Deploying frontend to Netlify...${NC}"
cd frontend
netlify deploy --prod

echo -e "${GREEN}Deployment complete!${NC}" 