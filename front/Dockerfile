# Use an official Node.js runtime as the base image
FROM node:14 as build-stage

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React app
RUN npm run build

# Use a lightweight Nginx image as the production base image
FROM nginx:alpine

# Copy the build files from the previous stage to the Nginx web server directory
COPY --from=build-stage /app/build /usr/share/nginx/html

# Expose port 80 for the Nginx server
EXPOSE 80

# Start Nginx when the container starts
CMD ["nginx", "-g", "daemon off;"]
