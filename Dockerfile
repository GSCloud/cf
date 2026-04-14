FROM node:alpine
RUN npm install -g wrangler && npm cache clean --force
WORKDIR /app
ENTRYPOINT ["wrangler"]