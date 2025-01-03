# Use the latest Node.js image
FROM node:latest AS base

# Set the working directory inside the container
WORKDIR /app

# Copy package.json and pnpm-lock.yaml (or package-lock.json) for dependency installation
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN npm install -g pnpm && pnpm install

# Copy the rest of the application files
COPY . .

# Build the Vite project
FROM base AS prod
RUN pnpm build

# Expose the port Vite will use
EXPOSE 5173
